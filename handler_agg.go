package main

import (
	"context"
	"fmt"
)

func handlerRss(s *state, cmd command) error {
	var feedUrl string
	if len(cmd.Args) != 1 {
		feedUrl = "https://www.wagslane.dev/index.xml"
	} else {
		feedUrl = cmd.Args[0]
	}

	feed, err := fetchFeed(context.Background(), feedUrl)
	if err != nil {
		return fmt.Errorf("couldn't fetch feed: %v", err)
	}

	feed.printFeed()
	return nil
}
