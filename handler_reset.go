package main

import (
	"context"
	"fmt"
)

func handlerReset(s *state, cmd command) error {
	ctx := context.Background()
	
	err := s.db.DeleteAllUsers(ctx,)
	if err != nil {
		return fmt.Errorf("error reseting database: %w", err)
	}

	fmt.Println("database reset successfully!")
	return nil
}
