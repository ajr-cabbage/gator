package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ajr-cabbage/gator/internal/config"
	"github.com/ajr-cabbage/gator/internal/database"
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
