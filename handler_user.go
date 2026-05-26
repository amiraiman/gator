package main

import (
	"context"
	"errors"
	"fmt"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}

	name := cmd.Args[0]
	ctx := context.Background()
	_, err := s.db.GetUser(ctx, name)
	if err != nil {
		return errors.New("User does not exists, please register first")
	}

	err = s.cfg.SetUser(name)
	if err != nil {
		return fmt.Errorf("cannot login as: %w", err)
	}

	fmt.Printf("You have logged in as %v!\n", cmd.Args[0])
	return nil
}
