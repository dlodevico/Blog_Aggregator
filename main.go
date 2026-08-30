package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq" // Drivers are imported anonymously for database/sql registration

	"gator/internal/config"
	"gator/internal/database"
)

func main() {
	// 1. Read config
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	// 2. Open connection to database using db_url from config
	db, err := sql.Open("postgres", cfg.DbURL)
	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}
	defer db.Close()

	// 3. Initialize queries instance
	dbQueries := database.New(db)

	// 4. Store state with config and db queries
	appState := &state{
		db:  dbQueries,
		cfg: &cfg,
	}

	// 5. Setup commands registry
	cmds := commands{
		registeredCmds: make(map[string]func(*state, command) error),
	}

	// Unauthenticated / public commands
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("feeds", handlerFeeds)

	// Authenticated commands wrapped with middleware
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("following", middlewareLoggedIn(handlerFollowing))
	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	cmds.register("browse", middlewareLoggedIn(handlerBrowse))

	// 6. Handle CLI arguments
	args := os.Args
	if len(args) < 2 {
		log.Fatalf("error: not enough arguments provided (usage: gator <command> [args...])")
	}

	cmdName := args[1]
	cmdArgs := args[2:]

	cmd := command{
		Name: cmdName,
		Args: cmdArgs,
	}

	// 7. Run command
	err = cmds.run(appState, cmd)
	if err != nil {
		log.Fatalf("error running command: %v", err)
	}
}
