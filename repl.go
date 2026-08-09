package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"errors"
	"fmt"
	"html"
	"io"
	"net/http"
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


func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.args) < 2 {
		return errors.New("the `addfeed` handler expects a two arguments: the `name` of the feed and the `url` of the feed")
	}
	ctx := context.Background()
	name := cmd.args[0]
	url := cmd.args[1]
	t := time.Now()

	user, err := s.db.GetUser(ctx, name)
	if err != nil {
		return err
	}

	rssfeed, err := s.db.CreateFeed(ctx, database.CreateFeedParams{
		ID: uuid.New(),
		CreatedAt: t,
		UpdatedAt: t,
		Name: name,
		Url: url,
		UserID: user.ID})
	if err != nil {
		return err
	}

	fmt.Println("Channel")
	fmt.Printf("  ID: %s\n", rssfeed.ID)
	fmt.Printf("  CreatedAt: %s\n", rssfeed.CreatedAt)
	fmt.Printf("  UpdatedAt: %s\n", rssfeed.UpdatedAt)
	fmt.Printf("  Name: %s\n", rssfeed.Name)
	fmt.Printf("  Url: %s\n", rssfeed.Url)
	fmt.Printf("  UserID: %s\n", rssfeed.UserID)

	return nil
}



func handlerAgg(s *state, cmd command) error {
	ctx := context.Background()
	rssfeed, err := fetchFeed(ctx, "https://www.wagslane.dev/index.xml")
    if err != nil {
		return err
	}
	// fmt.Printf("FEED: %+v\n", rssfeed)

	fmt.Println("Channel")
	fmt.Printf("  Title: %s\n", rssfeed.Channel.Title)
	fmt.Printf("  Link: %s\n", rssfeed.Channel.Link)
	fmt.Printf("  Description: %s\n", rssfeed.Channel.Description)

	for i, item := range rssfeed.Channel.Item {
		fmt.Printf("    Item %d\n", i)
		fmt.Printf("      Title: %s\n", item.Title)
		fmt.Printf("      Link: %s\n", item.Link)
		fmt.Printf("      Description: %s\n", item.Description)
		fmt.Printf("      PubDate: %s\n", item.PubDate)
	}
	return nil
}


type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`

}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, err
	}
	client := http.Client{
    	Timeout: 10 * time.Second,
	}
	req.Header.Set("User-Agent", "gator")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	rssResp := RSSFeed{}
	err = xml.Unmarshal(dat, &rssResp)
	if err != nil {
		return nil, err
	}

	rssResp.Channel.Title = html.UnescapeString(rssResp.Channel.Title)
	rssResp.Channel.Description = html.UnescapeString(rssResp.Channel.Description)
	for i := range rssResp.Channel.Item {
		rssResp.Channel.Item[i].Title = html.UnescapeString(rssResp.Channel.Item[i].Title)
		rssResp.Channel.Item[i].Description = html.UnescapeString(rssResp.Channel.Item[i].Description)
	}

	return &rssResp, nil
}









