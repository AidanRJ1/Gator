package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/AidanRJ1/gator/internal/database"
	"github.com/google/uuid"
)

func checkUserExists(s *state, ctx context.Context, name string) bool {
	_, err := s.db.GetUser(ctx, sql.NullString{String: name, Valid: true}) 
	return err == nil 
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}

	ctx := context.Background()
	name := cmd.Args[0]

	exists := checkUserExists(s, ctx, name)
	if !exists {
		return fmt.Errorf("%s doesn't exist, can't login", name)
	}

	err := s.cfg.SetUser(name)
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}

	fmt.Println("User set successfully!")
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}

	ctx := context.Background()
	name := cmd.Args[0]

	exists := checkUserExists(s, ctx, name)
	if exists {
		return fmt.Errorf("%s is already in use, pick a different name", name)
	}

	user, err := s.db.CreateUser(ctx, database.CreateUserParams{
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name: sql.NullString{String: name, Valid: true},
	})
	if err != nil {
		return fmt.Errorf("error adding user to database: %w", err)
	}
	fmt.Println("user successfully added to database!")
	fmt.Println(user) 

	err = s.cfg.SetUser(name)
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}

	return nil
}

func handlerUsers(s *state, cmd command) error {
	ctx := context.Background()
	users, err := s.db.GetUsers(ctx)
	if err != nil {
		return fmt.Errorf("error getting users from database: %w", err)
	}

	if len(users) == 0 {
		return fmt.Errorf("no users in database")
	}

	for _ ,user := range users {
		name := user.Name.String
		if name == s.cfg.CurrentUserName {
			fmt.Printf("* %s (current) \n", name)
			continue
		}
		fmt.Printf("* %s \n", name)
	}

	return nil
}