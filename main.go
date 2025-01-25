package main

import (
	"database/sql"
	"fmt"
	"gator/internal/config"
	"gator/internal/database"
	_ "github.com/lib/pq"
	"os"
)

func main() {
	fmt.Println("Initializing Gator...")
	config, err := config.Read()
	if err != nil {
		fmt.Println(err)
	}

	// open connection
	db, err := sql.Open("postgres", config.DbURL)
	if err != nil {
		fmt.Println("Error opening database", db)
		os.Exit(1)
	}
	defer db.Close()

	// put db in state
	dbQueries := database.New(db)

	prg_state := state{}
	prg_state.cfg = &config
	prg_state.db = dbQueries

	// add commands
	all_commands := commands{coms: make(map[string]func(*state, command) error, 0)}
	all_commands.register("login", handlerLogin)
	all_commands.register("register", handlerRegister)

	args_input := os.Args
	if len(args_input) < 2 {
		fmt.Println("Not enough arguments to do anything")
		os.Exit(1)
	}

	t_command := command{name: args_input[1], args: args_input[2:]}
	cerr := all_commands.run(&prg_state, t_command)
	if cerr != nil {
		os.Exit(1)
		fmt.Println(cerr)
	}
}
