package main

import (
	"context"
	"fmt"
	"gator/internal/database"
	"github.com/google/uuid"
	"time"
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

	_, err := s.db.GetUser(context.Background(), cmd.args[0])
	if err != nil {
		fmt.Println("User Doesn't exist")
		return err
	}

	s.cfg.SetUser(cmd.args[0])

	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("Not enough arguments, register command expects username")
	}

	// check if name exists already
	_, err := s.db.GetUser(context.Background(), cmd.args[0])
	if err == nil {
		fmt.Println("User already exists")
		return fmt.Errorf("User already exists")
	}

	// create user in db
	newUser := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0]}

	s.db.CreateUser(context.Background(), newUser)
	s.cfg.SetUser(cmd.args[0])
	fmt.Printf("Successfully added User: %v\n", cmd.args[0])
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
