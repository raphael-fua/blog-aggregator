package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/raphael-fua/blog-aggregator/internal/config"
	"github.com/raphael-fua/blog-aggregator/internal/database"
)


type state struct {
	db *database.Queries
	cfg *config.Config
}


type command struct {
	name string
	args []string
}


type commands struct {
	m map[string]func(*state, command) error
}


func (c *commands) run(s *state, cmd command) error {
	handlerFunc, ok := c.m[cmd.name]
	if !ok {
		return fmt.Errorf("command %s not registered", cmd.name)
	}
	return handlerFunc(s, cmd)
}


func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("the `login` handler expects a single argument, the username")
	}
	name := cmd.args[0]
	ctx := context.Background()
	_, err := s.db.GetUser(ctx, name)
    if err != nil {
		return errors.New("cannot login to an account that does not exist")
	}
	err = s.cfg.SetUser(name)
	if err != nil {
		return err
	}
	fmt.Println("User has been set")
	return nil
}



func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("the `register` handler expects a single argument, the username")
	}

	name := cmd.args[0]

	ctx := context.Background()
	_, err := s.db.GetUser(ctx, name)
	if err == nil {
		return fmt.Errorf("%s already registered", name)
	}
	if !errors.Is(err, sql.ErrNoRows) {
        return err
	}

	t := time.Now() 
	user, err := s.db.CreateUser(ctx, database.CreateUserParams{
		ID: uuid.New(),
		CreatedAt: t,
		UpdatedAt: t,
		Name: name,

	})
	if err != nil {
		return err
	}
    err = s.cfg.SetUser(name)
	if err != nil {
		return err
	}
	fmt.Printf("user %s was created\n", name)
	fmt.Printf("  ID: %v\n", user.ID)
	fmt.Printf("  CreatedAt: %v\n", user.CreatedAt)
	fmt.Printf("  UpdatedAt: %v\n", user.UpdatedAt)
	fmt.Printf("  Name: %v\n", user.Name)
	return nil
}


func handlerReset(s *state, cmd command) error {
	ctx := context.Background()
	err := s.db.ResetDatabase(ctx)
    if err != nil {
		return errors.New("failed to reset database")
	}
	fmt.Println("database has been reset")
	return nil
}


func handlerUsers(s *state, cmd command) error {
	ctx := context.Background()
	users, err := s.db.GetUsers(ctx)
    if err != nil {
		return errors.New("failed to get users")
	}
	for _, user := range users {
		if s.cfg.UserName == user.Name {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}
	return nil
}




