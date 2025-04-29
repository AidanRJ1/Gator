package main

import (
	"context"
	"fmt"
	"time"

	"github.com/AidanRJ1/gator/internal/database"
	"github.com/google/uuid"
)

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <url>", cmd.Name)
	}
	ctx := context.Background()
	url := cmd.Args[0]

	feed, err := s.db.GetFeedByUrl(ctx, url)
	if err != nil {
		return fmt.Errorf("error getting feed '%s' from database: %w", url, err)
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

func handlerFollowing(s *state, cmd command, user database.User) error {
	ctx := context.Background()

	follows, err := s.db.GetFeedFollowsForUser(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("error getting follows for user '%s' from database: %w", user.Name.String, err)
	}
	if len(follows) == 0 {
		return fmt.Errorf("user '%s' not following any feeds", user.Name.String)
	}

	for _, follow := range follows {
		fmt.Printf("Feed: %s \n", follow.FeedName)
		fmt.Printf("User: %s \n", follow.UserName.String)
		fmt.Println("---------- ---------- ----------")
	}

	return nil
}
