package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"gator/internal/database"
)

// handlerAddFeed creates a new feed and automatically follows it for the logged-in user.
func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.Args) < 2 {
		return errors.New("usage: addfeed <name> <url>")
	}

	name := cmd.Args[0]
	url := cmd.Args[1]

	now := time.Now().UTC()
	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		Name:      name,
		Url:       url,
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("could not create feed: %w", err)
	}

	ffRow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("created feed but failed to follow: %w", err)
	}

	fmt.Println("Feed created successfully:")
	fmt.Printf("  * ID:         %s\n", feed.ID)
	fmt.Printf("  * Created At: %s\n", feed.CreatedAt.Format(time.RFC3339))
	fmt.Printf("  * Updated At: %s\n", feed.UpdatedAt.Format(time.RFC3339))
	fmt.Printf("  * Name:       %s\n", feed.Name)
	fmt.Printf("  * URL:        %s\n", feed.Url)
	fmt.Printf("  * User ID:    %s\n", feed.UserID)
	fmt.Printf("Followed as '%s' by '%s'\n", ffRow.FeedName, ffRow.UserName)

	return nil
}

// handlerFeeds lists all feeds in the database alongside the user who created each feed.
func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("could not retrieve feeds: %w", err)
	}

	if len(feeds) == 0 {
		fmt.Println("No feeds found in the database.")
		return nil
	}

	for _, feed := range feeds {
		fmt.Printf("* Name: %s\n", feed.FeedName)
		fmt.Printf("  URL:  %s\n", feed.FeedUrl)
		fmt.Printf("  User: %s\n", feed.UserName)
		fmt.Println()
	}

	return nil
}
