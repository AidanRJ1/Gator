package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/AidanRJ1/gator/internal/database"
	"github.com/google/uuid"
)

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.Args) != 2 {
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
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		Url:       url,
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("error adding feed to database: %w", err)
	}

	fmt.Println("user successfully added to database!")
	fmt.Println(feed)

	return nil
}

func handlerFeeds(s *state, cmd command) error {
	ctx := context.Background()

	feeds, err := s.db.GetFeeds(ctx)
	if err != nil {
		return fmt.Errorf("error getting feeds from database: %w", err)
	}

	for i, feed := range feeds {
		fmt.Printf("Name: %s \n", feed.Name)
		fmt.Printf("URL: %s \n", feed.Url)
		user, err := s.db.GetUserById(ctx, feeds[i].UserID)
		if err != nil {
			return fmt.Errorf("error getting user '%v' from database: %w", feeds[i].UserID, err)
		}
		fmt.Printf("User: %s \n", user.Name.String)
		fmt.Println("---------- ---------- ----------")
	}

	return nil
}
