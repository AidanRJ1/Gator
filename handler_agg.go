package main

import (
	"context"
	"fmt"
	"html"
)

func handlerAgg(s *state, cmd command) error {
	/*
		if len(cmd.Args) != 1 {
			return fmt.Errorf("usage: %s <url>", cmd.Name)
		}
	*/

	//url := cmd.Args[0]
	url := "https://www.wagslane.dev/index.xml"
	ctx := context.Background()
	feed, err := fetchFeed(ctx, url)
	if err != nil {
		return err
	}

	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	for i, item := range feed.Channel.Item {
		feed.Channel.Item[i].Title = html.UnescapeString(item.Title)
		feed.Channel.Item[i].Description = html.UnescapeString(item.Description)
	}

	fmt.Printf("Title: %s \n", feed.Channel.Title)
	fmt.Printf("Link: %s \n", feed.Channel.Link)
	fmt.Printf("Description: %s \n", feed.Channel.Description)
	fmt.Printf("Items: \n")
	for _, item := range feed.Channel.Item {
		fmt.Printf("  - Title: %s \n", item.Title)
		fmt.Printf("  - Link: %s \n", item.Link)
		fmt.Printf("  - Description: %s \n", item.Description)
		fmt.Printf("  - Publication Date: %s \n", item.PubDate)
		fmt.Println("---------- ----------- -----------")
	}
	return nil
}
