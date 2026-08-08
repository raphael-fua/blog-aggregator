package main

import (
	"github.com/raphael-fua/blog-aggregator/internal/config"
	"fmt"
	"os"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	s := state{
		cfg: &cfg,
	}
	cmds := commands{
 		m: map[string]func(*state, command) error{
			"login": handlerLogin,
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





