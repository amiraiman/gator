package main

import (
	"fmt"
)

type command struct {
	Name string
	Args []string
}

type commands struct {
	registeredCommands map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	name := cmd.Name
	fn, ok := c.registeredCommands[name]
	if !ok {
		return fmt.Errorf("Command %v does not exists.", name)
	}

	return fn(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.registeredCommands[name] = f
}
