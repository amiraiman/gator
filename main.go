package main

import (
	"log"
	"os"

	"github.com/amiraiman/gator/internal/config"
)

type state struct {
	cfg *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Error when reading: %v", err)
	}

	programState := state{cfg: &cfg}
	availableCommands := commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}

	availableCommands.register("login", handlerLogin)

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
