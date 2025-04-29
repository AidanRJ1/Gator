package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/AidanRJ1/gator/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		ctx := context.Background()
		currentUser := s.cfg.CurrentUserName

		user, err := s.db.GetUser(ctx, sql.NullString{String: currentUser, Valid: true})
		if err != nil {
			return fmt.Errorf("error getting user '%s' from database: %w", currentUser, err)
		}

		return handler(s, cmd, user)
	}
}
