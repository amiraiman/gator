package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/amiraiman/gator/internal/database"
)

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}

	name := cmd.Args[0]
	ctx := context.Background()
	_, err := s.db.GetUser(ctx, name)
	if err != nil {
		// User not exists
		s.db.CreateUser(ctx, database.CreateUserParams{Name: name})
		fmt.Printf("Successfully registered as %v\n", name)
		return nil
	}

	return errors.New("The username has been taken")
}
