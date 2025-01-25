package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const package_loc string = "/workspace/github.com/mwtopping/gator/"
const configfile string = ".gatorconfig.json"

type Config struct {
	DbURL    string `json:"db_url"`
	Username string `json:"current_user_name"`
}

func (c *Config) SetUser(username string) {
	c.Username = username

	// convert to bytes
	bytes, _ := json.Marshal(c)

	home_dir, _ := os.UserHomeDir()

	err := os.WriteFile(home_dir+package_loc+configfile, bytes, 0666)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Username has been set to: %v", username)
}

func Read() (Config, error) {
	config := Config{}
	home_dir, dir_err := os.UserHomeDir()
	if dir_err != nil {
		return Config{}, dir_err
	}

	// read in json bytes
	bytes, read_err := os.ReadFile(home_dir + package_loc + configfile)
	if read_err != nil {
		return Config{}, read_err
	}

	jsonerr := json.Unmarshal(bytes, &config)
	if jsonerr != nil {
		return Config{}, jsonerr
	}

	return config, nil
}
