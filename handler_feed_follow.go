package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"gator/internal/database"
)

// handlerFollow creates a feed follow record for the logged-in user.
func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) < 1 {
		return errors.New("usage: follow <url>")
	}

	url := cmd.Args[0]

	feed, err := s.db.GetFeedByURL(context.Background(), url)
	if err != nil {
		return fmt.Errorf("could not find feed with URL '%s': %w", url, err)
	}

	now := time.Now().UTC()
	ffRow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("could not follow feed: %w", err)
	}

	fmt.Printf("User '%s' is now following feed '%s'\n", ffRow.UserName, ffRow.FeedName)
	return nil
}

// handlerFollowing lists all feeds followed by the logged-in user.
func handlerFollowing(s *state, cmd command, user database.User) error {
	follows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("could not retrieve feed follows: %w", err)
	}

	if len(follows) == 0 {
		fmt.Printf("User '%s' is not following any feeds.\n", user.Name)
		return nil
	}

	fmt.Printf("Feeds followed by %s:\n", user.Name)
	for _, follow := range follows {
		fmt.Printf("* %s\n", follow.FeedName)
	}

	return nil
}

// handlerUnfollow removes a feed follow record for the logged-in user by feed URL.
func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) < 1 {
		return errors.New("usage: unfollow <url>")
	}

	url := cmd.Args[0]

	err := s.db.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{
		UserID: user.ID,
		Url:    url,
	})
	if err != nil {
		return fmt.Errorf("could not unfollow feed: %w", err)
	}

	fmt.Printf("User '%s' unfollowed feed at URL '%s'\n", user.Name, url)
	return nil
}
