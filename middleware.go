package main

import (
	"context"
	"fmt"
	"gator/internal/database"
)

// middlewareLoggedIn wraps a command handler requiring authentication.
// It fetches the current user from the database and passes it directly to the inner handler.
func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		currentUser := s.cfg.CurrentUserName
		user, err := s.db.GetUser(context.Background(), currentUser)
		if err != nil {
			return fmt.Errorf("authentication failed for user '%s': %w", currentUser, err)
		}

		return handler(s, cmd, user)
	}
}
