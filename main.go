package main

import (
	"fmt"
	"log"
	"os"

	"github.com/amiraiman/gator/internal/config"
)

type state struct {
	cfg *config.Config
}

type commands struct {
	cmd map[string]func(*state, command) error
}

type command struct {
	name string
	args []string
}

func main() {
	cfg := config.Read()
	currentState := state{cfg: &cfg}
	availableCommands := commands{
		cmd: make(map[string]func(*state, command) error),
	}

	availableCommands.register("login", handlerLogin)

	if len(os.Args) < 2 {
		log.Fatalf("No command found: %v", os.Args)
	}

	currentCommand := command{
		name: os.Args[1],
		args: os.Args[2:],
	}

	err := availableCommands.run(&currentState, currentCommand)
	if err != nil {
		log.Fatal(err)
	}
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("Username is required")
	}

	err := s.cfg.SetUser(cmd.args[0])
	if err != nil {
		return err
	}

	fmt.Printf("You have logged in as %v!\n", cmd.args[0])
	return nil
}

func (c *commands) run(s *state, cmd command) error {
	fn, ok := c.cmd[cmd.name]
	if !ok {
		return fmt.Errorf("Command %v does not exists.", cmd.name)
	}

	err := fn(s, cmd)
	if err != nil {
		return fmt.Errorf("Failed to execute command: %v", err)
	}

	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.cmd[name] = f
}
