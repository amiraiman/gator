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

	currentState := state{cfg: &cfg}
	availableCommands := commands{
		cmd: make(map[string]func(*state, command) error),
	}

	availableCommands.register("login", handlerLogin)

	if len(os.Args) < 2 {
		log.Fatal("Usage: cli <command> [args...]")
	}

	currentCommand := command{
		name: os.Args[1],
		args: os.Args[2:],
	}

	err = availableCommands.run(&currentState, currentCommand)
	if err != nil {
		log.Fatal(err)
	}
}
