package rss

import (
	"database/sql"
	"time"
)

// ParsePublishedAt attempts to parse common RSS/Atom date formats into a sql.NullTime.
func ParsePublishedAt(pubDate string) sql.NullTime {
	if pubDate == "" {
		return sql.NullTime{Valid: false}
	}

	formats := []string{
		time.RFC1123Z,
		time.RFC1123,
		time.RFC822Z,
		time.RFC822,
		time.RFC3339,
		time.RFC3339Nano,
		"Mon, 02 Jan 2006 15:04:05 -0700",
		"Mon, 02 Jan 2006 15:04:05 MST",
		"2006-01-02T15:04:05Z",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, pubDate); err == nil {
			return sql.NullTime{Time: t.UTC(), Valid: true}
		}
	}

	return sql.NullTime{Valid: false}
}
