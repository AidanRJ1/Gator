package config

import (
	"encoding/json"
	"os"
)

const configFileName string = ".gatorconfig.json"

type Config struct {
	DbUrl string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func (config Config) SetUser(currentUser string) error {
	return nil
}

func Read() (Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return Config{}, err
	}
	filePath := homeDir + "/" + configFileName

	dat, err := os.ReadFile(filePath)
	if err != nil {
		return Config{}, err
	}

	var config Config
	err = json.Unmarshal(dat, &config)
	if err != nil {
		return Config{}, err
	}	
	
	return config, nil
}

func config() {

}