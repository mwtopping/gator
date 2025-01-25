package main

import (
	"fmt"
	"gator/internal/config"
	"os"
)

func main() {
	fmt.Println("Initializing Gator...")
	config, err := config.Read()
	if err != nil {
		fmt.Println(err)
	}
	prg_state := state{}
	prg_state.config = &config

	all_commands := commands{coms: make(map[string]func(*state, command) error, 0)}
	all_commands.register("login", handlerLogin)

	args_input := os.Args
	if len(args_input) < 2 {
		fmt.Println("Not enough arguments to do anything")
		os.Exit(1)
	}

	t_command := command{name: args_input[1], args: args_input[2:]}
	cerr := all_commands.run(&prg_state, t_command)
	if cerr != nil {
		fmt.Println(cerr)
	}
}
