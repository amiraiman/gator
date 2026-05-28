package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/amiraiman/gator/internal/config"
	"github.com/amiraiman/gator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	cfg *config.Config
	db  *database.Queries
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Error when reading: %v", err)
	}

	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		log.Fatalf("Cant connect to the database: %v", err)
	}
	defer db.Close()

	queries := database.New(db)
	programState := state{
		cfg: &cfg,
		db:  queries,
	}

	availableCommands := commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}
	availableCommands.register("login", handlerLogin)
	availableCommands.register("register", handlerRegister)
	availableCommands.register("reset", handlerReset)
	availableCommands.register("users", handlerUsers)
	availableCommands.register("agg", handlerRss)
	availableCommands.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	availableCommands.register("feeds", handlerListFeeds)
	availableCommands.register("follow", middlewareLoggedIn(handlerFollowFeed))
	availableCommands.register("following", middlewareLoggedIn(handlerFollowingFeed))
	availableCommands.register("unfollow", middlewareLoggedIn(handlerDeleteFollowingFeed))

	if len(os.Args) < 2 {
		log.Fatal("Usage: cli <command> [args...]")
	}

	err = availableCommands.run(&programState, command{
		Name: os.Args[1],
		Args: os.Args[2:],
	})

	if err != nil {
		log.Fatal(err)
	}
}
