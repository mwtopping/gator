package main

import (
	"context"
	"fmt"
	"gator/internal/database"
	"gator/internal/rss"
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

func handlerReset(s *state, cmd command) error {

	err := s.db.Reset(context.Background())
	if err != nil {
		fmt.Println("Unable to Reset Table")
		return err
	}

	return nil
}

func handlerAgg(s *state, cmd command) error {

	feed, err := rss.FetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		fmt.Println(err)
		return nil
	}

	fmt.Println(feed)

	return nil
}

func handlerListUsers(s *state, cmd command) error {

	users, err := s.db.ListUsers(context.Background())
	if err != nil {
		fmt.Println("Unable to Reset Table")
		return err
	}

	if len(users) == 0 {
		fmt.Println("No users to list")
	} else {
		for u := range users {
			if users[u] == s.cfg.Username {
				fmt.Printf("* %v (current)\n", users[u])
			} else {
				fmt.Printf("* %v\n", users[u])
			}
		}
	}

	return nil
}

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.ListFeeds(context.Background())
	if err != nil {
		fmt.Println(err)
		return err
	}

	for _, f := range feeds {
		username, err := s.db.GetFeedUser(context.Background(), f.UserID)
		if err != nil {
			fmt.Println("Could not get username for ID", f.UserID)
			fmt.Println(err)
			return err
		}
		fmt.Println("---")
		fmt.Printf("Name: %v", f.Name)
		fmt.Printf("URL: %v", f.Url)
		fmt.Printf("Creator: %v", username)
	}

	return nil
}

func handlerFollow(s *state, cmd command) error {

	if len(cmd.args) < 1 {
		fmt.Println("Not enough arguments")
		return fmt.Errorf("Not enough arguments")
	}

	url := cmd.args[0]

	userID, err := s.db.GetUser(context.Background(), s.cfg.Username)
	if err != nil {
		fmt.Println(err)
		return err
	}
	feed, err := s.db.GetFeedFromURL(context.Background(), url)
	if err != nil {
		fmt.Println(err)
		return err
	}

	newFollow := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    userID.ID,
		FeedID:    feed.ID}

	_, err = s.db.CreateFeedFollow(context.Background(), newFollow)
	if err != nil {
		fmt.Println("Error creating new feed")
		return err
	}

	fmt.Printf("Created Feed with name: %v for User: %v", feed.Name, s.cfg.Username)

	return nil
}

func handlerAddFeed(s *state, cmd command) error {

	if len(cmd.args) < 2 {
		fmt.Println("Not enough arguments")
		return fmt.Errorf("Not enough arguments")
	}

	name := cmd.args[0]
	url := cmd.args[1]

	user, err := s.db.GetUser(context.Background(), s.cfg.Username)
	if err != nil {
		fmt.Println("Unable to get User")
		fmt.Println(err)
		return err
	}

	newFeed := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		Url:       url,
		UserID:    user.ID}

	feed, err := s.db.CreateFeed(context.Background(), newFeed)
	if err != nil {
		fmt.Println("Error creating new feed")
		return err
	}

	fmt.Println(feed)

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
