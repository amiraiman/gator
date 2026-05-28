package main

import (
	"context"
	"fmt"
	"time"

	"github.com/amiraiman/gator/internal/database"
	"github.com/google/uuid"
)

func handlerListFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't get feeds: %v", err)
	}

	if len(feeds) == 0 {
		fmt.Println("No feeds found.")
		return nil
	}

	fmt.Printf("Found %d feeds:\n", len(feeds))
	for i, f := range feeds {
		creator, err := s.db.GetUserById(context.Background(), f.UserID)
		if err != nil {
			continue
		}

		fmt.Printf("%v) %v\n", (i + 1), f.Name)
		fmt.Printf("* URL: %v\n", f.Url)
		fmt.Printf("* User: %v\n", creator.Name)
		if i != (len(feeds) - 1) {
			fmt.Println()
		}
	}

	return nil
}

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("usage: %s <feed-name> <feed-url>", cmd.Name)
	}

	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("no user found, please login first. %v", s.cfg.CurrentUserName)
	}

	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		Name:      cmd.Args[0],
		Url:       cmd.Args[1],
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("couldn't create feed: %v", err)
	}

	fmt.Println("Feed created successfully:")
	printFeed(feed)
	fmt.Println()
	fmt.Println("==========================")
	return nil
}

func printFeed(f database.Feed) {
	fmt.Printf("* ID:\t\t\t%v\n", f.ID)
	fmt.Printf("* Created:\t\t\t%v\n", f.CreatedAt)
	fmt.Printf("* Updated:\t\t\t%v\n", f.UpdatedAt)
	fmt.Printf("* Name:\t\t\t%v\n", f.Name)
	fmt.Printf("* URL:\t\t\t%v\n", f.Url)
	fmt.Printf("* UserID:\t\t\t%v\n", f.UserID)
}
