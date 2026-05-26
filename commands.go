package main

import (
	"fmt"
)

type command struct {
	name string
	args []string
}

type commands struct {
	cmd map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	fn, ok := c.cmd[cmd.name]
	if !ok {
		return fmt.Errorf("Command %v does not exists.", cmd.name)
	}

	return fn(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.cmd[name] = f
}
