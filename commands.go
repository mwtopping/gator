package main

import (
	"fmt"
)

type command struct {
	name string
	args []string
}

type commands struct {
	coms map[string]func(*state, command) error
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("Not enough arguments, login command expects username")
	}

	s.config.SetUser(cmd.args[0])

	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.coms[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	if f, ok := c.coms[cmd.name]; ok {
		err := f(s, cmd)
		if err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("Command %v not found in command list", cmd.name)
}
