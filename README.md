# Gator: CLI RSS Feed Aggregator

**Gator** is a feature-rich, command-line RSS feed aggregator written in Go and backed by PostgreSQL. It allows users to manage feed subscriptions, follow/unfollow feeds, and run an automated background scraper that periodically fetches and aggregates RSS items across feeds.

---

## Features

* **User Management:** Multi-user support with account registration, authentication, and state management.
* **Feed Management:** Add, track, and list public RSS/Atom feeds.
* **Subscription System:** Follow and unfollow feeds with unique user-feed tracking.
* **Automated Scraper Loop:** Background worker using ticker polling and persistent state tracking (`last_fetched_at`) to fetch oldest/unfetched feeds first.
* **Database & Migrations:** Type-safe database queries generated via `sqlc` and versioned schema migrations managed by `goose`.

---

## Tech Stack & Dependencies

* **Language:** Go (1.20+)
* **Database:** PostgreSQL
* **Code Generation:** [sqlc](https://sqlc.dev/)
* **Database Migrations:** [goose](https://github.com/pressly/goose)
* **Driver:** `github.com/lib/pq`
* **UUID Management:** `github.com/google/uuid`

---

## Installation & Setup

### 1. Prerequisites

Ensure you have the following installed on your system:
* [Go](https://go.dev/doc/install)
* [PostgreSQL](https://www.postgresql.org/download/)
* [Goose](https://github.com/pressly/goose#installation) (for database migrations)
* [sqlc](https://docs.sqlc.dev/en/stable/overview/install.html) (optional, if modifying SQL queries)

### 2. Database Configuration

Create a local PostgreSQL database named `gator`:

```bash
createdb gator

Run schema migrations to set up required tables (users, feeds, feed_follows):

cd sql/schema
goose postgres "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable" up
cd ../..

3. Application Configuration
Create a .gatorconfig.json file in your home directory (~/.gatorconfig.json) to specify your database connection string:
{
  "db_url": "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable",
  "current_user_name": ""
}

4. Build the Project
Compile the binary from the repository root:
go build -o gator .

gator/
├── main.go                  # Entry point, CLI command routing, and state setup
├── state.go                 # State and config context definitions
├── commands.go              # CLI command registration registry
├── middleware.go            # Higher-order middleware for authenticated routes
├── handler_user.go          # User auth and management handlers
├── handler_feed.go          # Feed creation and retrieval handlers
├── handler_feed_follow.go   # Subscription and follow/unfollow handlers
├── handler_agg.go           # Background scraper ticker loop & fetcher
├── rss.go                   # HTTP client and XML RSS parser
├── sql/
│   ├── schema/              # Goose SQL schema migrations
│   └── queries/             # Raw SQL queries compiled by sqlc
└── internal/
    ├── config/              # Configuration file read/write logic
    ├── database/            # Generated type-safe Go database code from sqlc
    └── rss/                 # Core RSS fetching models
