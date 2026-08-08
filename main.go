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
	cfg.SetUser("Fua-Marcou")
	cfg, err = config.Read()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(cfg.DbUrl)
	fmt.Println(cfg.UserName)
}

