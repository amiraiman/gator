package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/amiraiman/gator/internal/database"
	"github.com/google/uuid"
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
		u, err := s.db.CreateUser(ctx, database.CreateUserParams{Name: name, ID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now()})
		if err != nil {
			return err
		}

		s.cfg.SetUser(name)
		fmt.Printf("Successfully registered as %v\n", u)
		return nil
	}

	return errors.New("The username has been taken")
}
