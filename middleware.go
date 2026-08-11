package main


import (
	"context"
	// "database/sql"
	// "errors"
	// "github.com/raphael-fua/blog-aggregator/internal/config"
	"github.com/raphael-fua/blog-aggregator/internal/database"
)


func middlewareLoggedIn(
	handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.cfg.UserName)
		if err != nil {
			return err
		}
		return handler(s, cmd, user)
	}
}







