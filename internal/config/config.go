package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

const configFileName string = ".gatorconfig.json"

// write a Config{} to the default config file location
func write(cfg Config) error {
	configData, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	targetFilepath := filepath.Join(homeDir, configFileName)

	err = os.WriteFile(targetFilepath, configData, 0644)
	if err != nil {
		return err
	}
	return nil
}

// read the config file from default location and unmarshal to Config{}
func Read() (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return &Config{}, err
	}

	targetFilepath := filepath.Join(homeDir, configFileName)
	confFile, err := os.ReadFile(targetFilepath)
	if err != nil {
		return &Config{}, err
	}

	newConfig := Config{}

	err = json.Unmarshal(confFile, &newConfig)
	return &newConfig, nil
}

// set the current user to the struct and write()
func (c *Config) SetUser(userName string) error {
	if c == nil {
		return errors.New("Error: nil config")
	}

	c.CurrentUserName = userName

	err := write(*c)
	if err != nil {
		return err
	}

	return nil
}
