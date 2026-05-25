package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

// Export a Read function that reads the JSON file found at ~/.gatorconfig.json
// and returns a Config struct. It should read the file from the HOME directory,
// then decode the JSON string into a new Config struct. I used os.UserHomeDir
// to get the location of HOME.
func Read() Config {
	fileName, err := getConfigFilePath()
	if err != nil {
		log.Fatal(err)
	}

	fileContent, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatalf("Cannot read config: %v", err)
	}

	cfg := Config{}
	err = json.Unmarshal(fileContent, &cfg)
	if err != nil {
		log.Fatalf("Cannot decode config: %v", err)
	}

	return cfg
}

// Export a SetUser method on the Config struct that writes the config struct
// to the JSON file after setting the current_user_name field.
func (c Config) SetUser(userName string) {
	c.CurrentUserName = userName
	err := write(c)
	if err != nil {
		log.Fatalf("Error when writing config: %v", err)
	}
}

func write(c Config) error {
	f, err := json.Marshal(c)
	if err != nil {
		return fmt.Errorf("Cannot encode config: %v", err)
	}

	fileName, err := getConfigFilePath()
	if err != nil {
		return fmt.Errorf("Cannot get config file name: %v", err)
	}

	err = os.WriteFile(fileName, f, 0644)
	if err != nil {
		return fmt.Errorf("Cannot write config: %v", err)
	}

	return nil
}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("Cannot get home dir: %v", err)
	}

	return homeDir + string(os.PathSeparator) + configFileName, nil
}
