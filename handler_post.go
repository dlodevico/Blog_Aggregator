package main

import (
	"context"
	"fmt"
	"strconv"

	"gator/internal/database"
)

// handlerBrowse lists posts from feeds followed by the logged-in user.
func handlerBrowse(s *state, cmd command, user database.User) error {
	limit := int32(2) // Default limit

	if len(cmd.Args) > 0 {
		parsedLimit, err := strconv.Atoi(cmd.Args[0])
		if err != nil || parsedLimit <= 0 {
			return fmt.Errorf("invalid limit parameter '%s': must be a positive integer", cmd.Args[0])
		}
		limit = int32(parsedLimit)
	}

	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  limit,
	})
	if err != nil {
		return fmt.Errorf("could not fetch posts: %w", err)
	}

	if len(posts) == 0 {
		fmt.Printf("No posts found for user '%s'.\n", user.Name)
		return nil
	}

	fmt.Printf("--- Showing latest %d posts for %s ---\n\n", len(posts), user.Name)
	for _, post := range posts {
		pubDate := "Unknown"
		if post.PublishedAt.Valid {
			pubDate = post.PublishedAt.Time.Format("2006-01-02 15:04:05 UTC")
		}

		fmt.Printf("Title:       %s\n", post.Title)
		fmt.Printf("Published:   %s\n", pubDate)
		fmt.Printf("Link:        %s\n", post.Url)
		if post.Description.Valid && post.Description.String != "" {
			fmt.Printf("Description: %s\n", post.Description.String)
		}
		fmt.Println("--------------------------------------------------")
	}

	return nil
}
