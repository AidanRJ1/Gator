package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/AidanRJ1/gator/internal/database"
	"github.com/google/uuid"
)

func handlerFollow(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <url>", cmd.Name)
	}
	ctx := context.Background()
	url := cmd.Args[0]
	currentUser := s.cfg.CurrentUserName

	feed, err := s.db.GetFeedByUrl(ctx, url)
	if err != nil {
		return fmt.Errorf("error getting feed '%s' from database: %w", url, err)
	}

	user, err := s.db.GetUser(ctx, sql.NullString{String: currentUser, Valid: true})
	if err != nil {
		return fmt.Errorf("error getting user '%s' from database: %w", currentUser, err)
	}

	follow, err := s.db.CreateFeedFollows(ctx, database.CreateFeedFollowsParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("error adding follow to database: %w", err)
	}

	fmt.Println("follow successfully added to database!")
	fmt.Printf("Feed: %s \n", follow.FeedName)
	fmt.Printf("User: %s \n", follow.UserName.String)

	return nil
}

func handlerFollowing(s *state, cmd command) error {
	ctx := context.Background()
	currentUser := s.cfg.CurrentUserName

	user, err := s.db.GetUser(ctx, sql.NullString{String: currentUser, Valid: true})
	if err != nil {
		return fmt.Errorf("error getting user '%s' from database: %w", currentUser, err)
	}

	follows, err := s.db.GetFeedFollowsForUser(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("error getting follows for user '%s' from database: %w", currentUser, err)
	}
	if len(follows) == 0 {
		return fmt.Errorf("user '%s' not following any feeds", currentUser)
	}

	for _, follow := range follows {
		fmt.Printf("Feed: %s \n", follow.FeedName)
		fmt.Printf("User: %s \n", follow.UserName.String)
		fmt.Println("---------- ---------- ----------")
	}

	return nil
}
