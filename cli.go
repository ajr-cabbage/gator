package main

import (
	"errors"
	"fmt"

	"github.com/ajr-cabbage/gator/internal/config"
)

type state struct {
	conf *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	cmdFuncs map[string]func(*state, command) error
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

	err := s.conf.SetUser(cmd.args[0])
	if err != nil {
		return err
	}

	fmt.Printf("User has been set to '%s'\n", cmd.args[0])

	return nil
}
