package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/AidanRJ1/gator/internal/database"
	"github.com/google/uuid"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedUrl string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedUrl, nil)
	if err != nil {
		return &RSSFeed{}, err
	}

	req.Header.Set("User-Agent", "gator")
	client := http.DefaultClient

	res, err := client.Do(req)
	if err != nil {
		return &RSSFeed{}, err
	}

	dat, err := io.ReadAll(res.Body)
	if err != nil {
		return &RSSFeed{}, err
	}

	var feed RSSFeed
	err = xml.Unmarshal(dat, &feed)
	if err != nil {
		return &RSSFeed{}, err
	}

	return &feed, nil
}


func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.Args) !=  2 {
		return fmt.Errorf("usage: %s <name> <url>", cmd.Name)
	}
	ctx := context.Background()

	name := cmd.Args[0]
	url := cmd.Args[1]
	userName := s.cfg.CurrentUserName

	user, err := s.db.GetUser(ctx, sql.NullString{String: userName, Valid: true})
	if err != nil {
		return fmt.Errorf("error getting user '%s' from database: %w", userName, err)
	}

	feed, err := s.db.CreateFeed(ctx, database.CreateFeedParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: name,
		Url: url,
		UserID: user.ID,
	})
	if err != nil {
		return fmt.Errorf("error adding feed to database: %w", err)
	}

	fmt.Println("user successfully added to database!")
	fmt.Println(feed)

	return nil
}

func handlerFeeds(s *state, cmd command) error {
	return nil
}