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

func handlerFollowFeed(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <url>", cmd.Name)
	}

	feed, err := s.db.GetFeedByUrl(context.Background(), cmd.Args[0])
	if err != nil {
		return fmt.Errorf("feed is not created yet: %v", err)
	}

	follow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("couldn't follow feed: %v", err)
	}

	fmt.Printf("New feed followed:\n")
	fmt.Printf("Feed: %v\n", follow.UserName)
	fmt.Printf("User: %v\n", follow.FeedName)
	return nil
}

func handlerFollowingFeed(s *state, cmd command, user database.User) error {
	following, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("couldn't get following feed: %v", err)
	}

	if len(following) == 0 {
		fmt.Println("Not following any feed")
		return nil
	}

	fmt.Printf("Feeds Followed By You (%v):\n", user.Name)
	for i, feed := range following {
		fmt.Printf("%v) %v\n", (i + 1), feed.FeedName)
	}

	return nil
}

func handlerDeleteFollowingFeed(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <feed-url>", cmd.Name)
	}

	feed, err := s.db.GetFeedByUrl(context.Background(), cmd.Args[0])
	if err != nil {
		return fmt.Errorf("feed url doesnt exist: %v", err)
	}

	err = s.db.DeleteFollowFeedsByUserAndFeedId(context.Background(), database.DeleteFollowFeedsByUserAndFeedIdParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("couldn't unfollow: %v", err)
	}

	fmt.Printf("%f unfollowed successfully\n", feed.Name)
	return nil
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("usage: %s <feed-name> <feed-url>", cmd.Name)
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

	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("couldn't follow feed: %v", err)
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
