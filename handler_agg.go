package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/AidanRJ1/gator/internal/database"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <duration>", cmd.Name)
	}

	duration, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("duration was invalid: %w", err)
	}
	fmt.Printf("Collecting feeds every %s", cmd.Args[0])

	ticker := time.NewTicker(duration)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func scrapeFeeds(s *state) {
	ctx := context.Background()

	feed, err := s.db.GetNextFeedToFetch(ctx)
	if err != nil {
		log.Println("error fetching next feed", err)
		return
	}
	log.Println("Found a feed to fetch!")
	scrapeFeed(s.db, feed)
}

func scrapeFeed(db *database.Queries, feed database.Feed) {
	ctx := context.Background()
	err := db.MarkFeedFetched(ctx, database.MarkFeedFetchedParams{
		ID: feed.ID,
		UpdatedAt: time.Now().UTC(),
		LastFetchedAt: sql.NullTime{Time: time.Now().UTC(), Valid: true},
	})
	if err != nil {
		log.Printf("Couldn't mark feed '%s' fetched: %v \n", feed.Name, err)
	}

	rssFeed, err := fetchFeed(ctx, feed.Url)
	if err != nil {
		log.Printf("Couldn't collect feed '%s': %v \n", feed.Name, err)
	}

	/*
	rssFeed.Channel.Title = html.UnescapeString(rssFeed.Channel.Title)
	rssFeed.Channel.Description = html.UnescapeString(rssFeed.Channel.Description)
	for i, item := range rssFeed.Channel.Item {
		rssFeed.Channel.Item[i].Title = html.UnescapeString(item.Title)
		rssFeed.Channel.Item[i].Description = html.UnescapeString(item.Description)
	}

	fmt.Printf("Title: %s \n", rssFeed.Channel.Title)
	fmt.Printf("Link: %s \n", rssFeed.Channel.Link)
	fmt.Printf("Description: %s \n", rssFeed.Channel.Description)
	fmt.Printf("Items: \n")
	for _, item := range rssFeed.Channel.Item {
		fmt.Printf("  - Title: %s \n", item.Title)
		fmt.Printf("  - Link: %s \n", item.Link)
		fmt.Printf("  - Description: %s \n", item.Description)
		fmt.Printf("  - Publication Date: %s \n", item.PubDate)
		fmt.Println("---------- ----------- -----------")
	}
	*/

	for _, item := range rssFeed.Channel.Item {
		fmt.Printf("Found post: %s \n", item.Title)
	}
	log.Printf("Feed %s collected, %v posts found", feed.Name, len(rssFeed.Channel.Item))
}