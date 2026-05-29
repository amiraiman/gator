package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/amiraiman/gator/internal/database"
)

func handlerRss(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %v <time-between-req>", cmd.Name)
	}

	timeString := cmd.Args[0]
	timeBetweenRequests, err := time.ParseDuration(timeString)
	if err != nil {
		return fmt.Errorf("invalid time argument, expecting (1s / 1m / 1h) but received: %v", timeString)
	}

	fmt.Printf("Collecting feeds every %v\n", timeBetweenRequests.String())
	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		err = scrapeFeeds(s, cmd)
		if err != nil {
			return err
		}
	}
}

func scrapeFeeds(s *state, cmd command) error {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't fetch feed to scrape: %v", err)
	}

	err = s.db.MarkFeedFetchedById(context.Background(), database.MarkFeedFetchedByIdParams{
		ID:            feed.ID,
		UpdatedAt:     time.Now().UTC(),
		LastFetchedAt: sql.NullTime{Time: time.Now().UTC(), Valid: true},
	})
	if err != nil {
		return fmt.Errorf("couldn't mark feed %s as fetched: %v", feed.Name, err)
	}

	fmt.Printf("Fetching %v\n", feed.Name)
	fetchedFeed, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		return fmt.Errorf("couldn't scrape feed: %v", err)
	}

	fetchedFeed.printFeed()
	return nil
}
