package config

import (
	"encoding/json"
	"fmt"
	"os"
)

//This package should have the following functionality exported so the main package can use it:
//
//    Export a Config struct that represents the JSON file structure, including struct tags.
//    Export a Read function that reads the JSON file found at ~/.gatorconfig.json and returns a Config struct. It should read the file from the HOME directory, then decode the JSON string into a new Config struct. I used os.UserHomeDir to get the location of HOME.
//    Export a SetUser method on the Config struct that writes the config struct to the JSON file after setting the current_user_name field.
//
//I also wrote a few non-exported helper functions and added a constant to hold the filename.
//
//    getConfigFilePath() (string, error)
//    write(cfg Config) error
//    const configFileName = ".gatorconfig.json"
//
//But you can implement the internals of the package however you like.

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
