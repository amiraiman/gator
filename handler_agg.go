package main

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/amiraiman/gator/internal/database"
	"github.com/google/uuid"
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

	for _, post := range fetchedFeed.Channel.Items {
		publishedDate, err := time.Parse(time.RFC1123, post.PublishedDate)
		if err != nil {
			fmt.Printf("couldn't parse published date for \"%v\": %v\n", post.Title, post.PublishedDate)
			continue
		}

		savedPost, err := s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Title:       post.Title,
			Url:         post.Link,
			Description: post.Description,
			PublishedAt: publishedDate,
			FeedID:      feed.ID,
		})

		if err != nil {
			if !strings.Contains(err.Error(), "pq: duplicate key value violates unique constraint") {
				fmt.Printf("couldn't save post %v to database: %v\n", post.Title, err)
			} else {
				fmt.Printf("post %v already exists in the database\n", post.Title)
			}
		} else {
			fmt.Printf("post %v saved to database\n", savedPost.Title)
		}
	}

	return nil
}
