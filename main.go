package main

import (
	"fmt"
	"gator/internal/config"
)

func main() {
	fmt.Println("Initializing Gator...")
	config, err := config.Read()
	if err != nil {
		fmt.Println(err)
	}
	config.SetUser("michael")
	fmt.Println(config)

}
