package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/AidanRJ1/gator/internal/config"
	"github.com/AidanRJ1/gator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	db *database.Queries
	cfg *config.Config
}

func main() {
	// Read config from local file '.gator.config'
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config file: %v", err)
	}

	// Open connection to database
	db, err := sql.Open("postgres", cfg.DbUrl)
	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}
	dbQueries := database.New(db)

	// Assign config and database to app state
	var programState = state{
		cfg: &cfg,
		db: dbQueries,
	}

	// Create list of cli commands and register each
	var cmds = commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("feeds", handlerFeeds)
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("following", middlewareLoggedIn(handlerFollowing))
	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	cmds.register("browse", middlewareLoggedIn(handlerBrowse))

	// Check at least two arguments are provided: command + args
	if len(os.Args) < 2 {
		log.Fatal("Usage: cli <command> [args...]")
		return
	}

	commandName := os.Args[1]
	commandArgs := os.Args[2:]

	err = cmds.run(&programState, command{Name: commandName, Args: commandArgs})
	if err != nil {
		log.Fatal(err)
	}
}
