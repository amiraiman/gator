package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DBURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

// Export a Read function that reads the JSON file found at ~/.gatorconfig.json
// and returns a Config struct. It should read the file from the HOME directory,
// then decode the JSON string into a new Config struct. I used os.UserHomeDir
// to get the location of HOME.
func Read() (Config, error) {
	fileName, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	fileContent, err := os.ReadFile(fileName)
	if err != nil {
		return Config{}, err
	}

	cfg := Config{}
	err = json.Unmarshal(fileContent, &cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}

// Export a SetUser method on the Config struct that writes the config struct
// to the JSON file after setting the current_user_name field.
func (c *Config) SetUser(userName string) error {
	c.CurrentUserName = userName
	return write(*c)
}

func write(c Config) error {
	fileName, err := getConfigFilePath()
	if err != nil {
		return err
	}

	f, err := json.Marshal(c)
	if err != nil {
		return err
	}

	err = os.WriteFile(fileName, f, 0644)
	if err != nil {
		return err
	}

	return nil
}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(homeDir, configFileName), nil
}
