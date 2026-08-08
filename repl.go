package main

import (
	// "bufio"
	"fmt"
	// "os"
	// "strings"
	"errors"
	"github.com/raphael-fua/blog-aggregator/internal/config"
)


type state struct {
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


func (c *commands) register(name string, f func(*state, command) error) {
	// if c.m == nil {
	// 	fmt.Println("`commands` field `m map[string]func(*state, command) error` is nil")
	// 	return
	// }
	c.m[name] = f
}


func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("the `login` handler expects a single argument, the username")
	}
	err := s.cfg.SetUser(cmd.args[0])
	if err != nil {
		return err
	}
	fmt.Println("User has been set")
	return nil
}

// func getCommands() map[string]command {
// 	return map[string]command{
// 		"login": {
// 			name: "login",
// 			args: []string,
// 		},
// 		"map": {
// 			name:        "map",
// 			description: "Gets the next page of locations",
// 			callback:    commandMapf,
// 		},
// 		"mapb": {
// 			name:        "mapb",
// 			description: "Gets the previous page of locations",
// 			callback:    commandMapb,
// 		},
// 		"exit": {
// 			name:        "exit",
// 			description: "Exits the Pokedex",
// 			callback:    commandExit,
// 		},
// 		"explore": {
// 			name:        "explore",
// 			description: "Explores the location passed to the command",
// 			callback:    commandExplore,
// 		},
// 		"catch": {
// 			name:        "catch <pokemon_name>",
// 			description: "Attempt to catch a pokemon",
// 			callback:    commandCatch,
// 		},
// 		"inspect": {
// 			name:        "inspect <pokemon_name>",
// 			description: "Get info of a captured pokemon",
// 			callback:    commandInspect,
// 		},
// 		"pokedex": {
// 			name:        "pokedex",
// 			description: "prints list of all names of caught pokemon",
// 			callback:    commandPokedex,
// 		},
// 	}
// }






