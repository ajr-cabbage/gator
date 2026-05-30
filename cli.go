package main

import (
	"context"
	"errors"
	"fmt"
	"html"
	"log"
	"os"
	"time"

	"github.com/ajr-cabbage/gator/internal/config"
	"github.com/ajr-cabbage/gator/internal/database"
	"github.com/ajr-cabbage/gator/internal/rss"
	"github.com/google/uuid"
)

type state struct {
	conf *config.Config
	db   *database.Queries
}

type command struct {
	name string
	args []string
}

type commands struct {
	cmdFuncs map[string]func(*state, command) error
}

func getCommands() *commands {
	// initialize commands{} and underlying map
	cmds := new(commands)
	cmdMap := make(map[string]func(*state, command) error)
	cmds.cmdFuncs = cmdMap
	// register commands
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", handlerAddFeed)
	cmds.register("feeds", handlerFeeds)

	return cmds
}

func buildCommand() command {
	// get CLI arguments
	cliArgs := os.Args
	if len(cliArgs) < 2 {
		log.Fatal("No command given")
	}
	// set cmd name and args
	cmd := command{
		name: cliArgs[1],
		args: cliArgs[2:],
	}
	return cmd
}

func (c *commands) run(s *state, cmd command) error {
	f, ok := c.cmdFuncs[cmd.name]
	if !ok {
		return errors.New("Error: unknown command.")
	}

	err := f(s, cmd)
	if err != nil {
		return err
	}

	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.cmdFuncs[name] = f
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("Error: A username is required")
	}

	_, err := s.db.GetUser(context.Background(), cmd.args[0])
	if err != nil {
		return errors.New("Error: Username not found")
	}

	err = s.conf.SetUser(cmd.args[0])
	if err != nil {
		return err
	}

	fmt.Printf("User has been set to '%s'\n", cmd.args[0])

	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("Error: A name is required")
	}
	newUserParams := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
	}
	_, err := s.db.CreateUser(context.Background(), newUserParams)
	if err != nil {
		return err
	}
	s.conf.SetUser(cmd.args[0])
	fmt.Printf("New user '%s' created\n", s.conf.CurrentUserName)
	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.db.ResetUsers(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}

	for _, user := range users {
		if user == s.conf.CurrentUserName {
			fmt.Printf("* %s (current)\n", user)
		} else {
			fmt.Printf("* %s\n", user)
		}
	}

	return nil
}

func handlerAgg(s *state, cmd command) error {
	url := "https://www.wagslane.dev/index.xml"
	feed, err := rss.FetchFeed(context.Background(), url)
	if err != nil {
		return err
	}
	// unescape Title and Description fields
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	for i := range feed.Channel.Item {
		feed.Channel.Item[i].Title = html.UnescapeString(feed.Channel.Item[i].Title)
		feed.Channel.Item[i].Description = html.UnescapeString(feed.Channel.Item[i].Description)
	}

	fmt.Println(*feed)

	return nil
}

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.args) < 2 {
		return errors.New("Error: name and url required")
	}

	usr, err := s.db.GetUser(context.Background(), s.conf.CurrentUserName)

	feedParams := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
		Url:       cmd.args[1],
		UserID:    usr.ID,
	}

	feed, err := s.db.CreateFeed(context.Background(), feedParams)
	if err != nil {
		return err
	}

	fmt.Printf("ID: %s\n", feed.ID)
	fmt.Printf("CreatedAt: %s\n", feed.CreatedAt)
	fmt.Printf("UpdatedAt: %s\n", feed.UpdatedAt)
	fmt.Printf("Name: %s\n", feed.Name)
	fmt.Printf("URL: %s\n", feed.Url)
	fmt.Printf("UserID: %s\n", feed.UserID)

	return nil
}

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}

	for _, feed := range feeds {
		user, err := s.db.GetUserName(context.Background(), feed.UserID)
		if err != nil {
			user = "Invalid UserID"
		}
		fmt.Printf("Name: %s\n", feed.Name)
		fmt.Printf("URL: %s\n", feed.Url)
		fmt.Printf("User: %s\n", user)
	}

	return nil
}
