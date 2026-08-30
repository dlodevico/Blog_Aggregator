package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"os"
	"time"

	"github.com/google/uuid"

	"gator/internal/database"
	"gator/internal/rss"
)

// handlerRegister creates a new user in the database and sets them as the current active user.
func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) == 0 {
		return errors.New("usage: register <name>")
	}

	name := cmd.Args[0]

	now := time.Now().UTC()
	userParams := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		Name:      name,
	}

	user, err := s.db.CreateUser(context.Background(), userParams)
	if err != nil {
		fmt.Printf("Error creating user: %v\n", err)
		os.Exit(1)
	}

	err = s.cfg.SetUser(name)
	if err != nil {
		return fmt.Errorf("could not set active user: %w", err)
	}

	fmt.Printf("User '%s' created successfully!\n", name)
	fmt.Printf("User Details: %+v\n", user)

	return nil
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) == 0 {
		return errors.New("usage: login <username>")
	}

	username := cmd.Args[0]

	// Verify the user exists in the database
	_, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		fmt.Printf("User '%s' does not exist\n", username)
		os.Exit(1)
	}

	// Set the current user in the config file
	err = s.cfg.SetUser(username)
	if err != nil {
		return fmt.Errorf("could not set active user: %w", err)
	}

	fmt.Printf("User has been set to: %s\n", username)
	return nil
}

// handlerReset deletes all users from the database.
func handlerReset(s *state, cmd command) error {
	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		fmt.Printf("Failed to reset database: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Database successfully reset!")
	return nil
}

// handlerUsers fetches and prints all registered users from the database.
func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		fmt.Printf("Failed to retrieve users: %v\n", err)
		os.Exit(1)
	}

	currentUser := s.cfg.CurrentUserName

	for _, user := range users {
		if user.Name == currentUser {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}

	return nil
}

func scrapeFeeds(s *state) error {
	ctx := context.Background()

	feed, err := s.db.GetNextFeedToFetch(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return fmt.Errorf("could not get next feed: %w", err)
	}

	now := time.Now().UTC()
	_, err = s.db.MarkFeedFetched(ctx, database.MarkFeedFetchedParams{
		LastFetchedAt: sql.NullTime{Time: now, Valid: true},
		ID:            feed.ID,
	})
	if err != nil {
		return fmt.Errorf("could not mark feed as fetched: %w", err)
	}

	rssFeed, err := rss.FetchFeed(ctx, feed.Url)
	if err != nil {
		return fmt.Errorf("could not fetch feed '%s': %w", feed.Name, err)
	}

	for _, item := range rssFeed.Channel.Item {
		description := sql.NullString{
			String: item.Description,
			Valid:  item.Description != "",
		}

		pubTime := rss.ParsePublishedAt(item.PubDate)

		_, err := s.db.CreatePost(ctx, database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   now,
			UpdatedAt:   now,
			Title:       item.Title,
			Url:         item.Link,
			Description: description,
			PublishedAt: pubTime,
			FeedID:      feed.ID,
		})
		if err != nil {
			// Ignore duplicate URL constraint violations
			if strings.Contains(err.Error(), "unique constraint") || strings.Contains(err.Error(), "posts_url_key") {
				continue
			}
			fmt.Printf("Error saving post '%s': %v\n", item.Title, err)
		}
	}

	fmt.Printf("Fetched feed '%s' (%d items processed)\n", feed.Name, len(rssFeed.Channel.Item))
	return nil
}

// handlerAgg runs continuous aggregation on a given ticker duration string (e.g., 1s, 1m).
func handlerAgg(s *state, cmd command) error {
	if len(cmd.Args) < 1 {
		return errors.New("usage: agg <time_between_reqs>")
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("invalid duration format '%s': %w", cmd.Args[0], err)
	}

	fmt.Printf("Collecting feeds every %s\n", timeBetweenRequests)

	ticker := time.NewTicker(timeBetweenRequests)
	defer ticker.Stop()

	// Execute immediately on start, then wait for subsequent ticks
	for ; ; <-ticker.C {
		if err := scrapeFeeds(s); err != nil {
			fmt.Printf("Scraping error: %v\n", err)
		}
	}
}
