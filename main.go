package main

import (
	"github.com/raphael-fua/blog-aggregator/internal/config"
	"github.com/raphael-fua/blog-aggregator/internal/database"
	"fmt"
	"os"
	_ "github.com/lib/pq"
	"database/sql"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	db, err := sql.Open("postgres", cfg.DbUrl)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	s := state{
		db: database.New(db),
		cfg: &cfg,
	}
	cmds := commands{
 		m: map[string]func(*state, command) error{
			"login": handlerLogin,
			"register": handlerRegister,
			"reset": handlerReset,
			"users": handlerUsers,
		},
	}
	cmdLine := os.Args
	if len(cmdLine) < 2 {
		fmt.Println("not enough args")
		os.Exit(1)
	}
	cmd := command{
		name: cmdLine[1],
		args: cmdLine[2:],
	}

	err = cmds.run(&s, cmd)
	if err != nil{
		fmt.Println(err)
		os.Exit(1)
	}
}






