package main

import (
	"gator/internal/config"
	"gator/internal/database"
)

// state holds application state shared across command handlers.
type state struct {
	db  *database.Queries
	cfg *config.Config
}
