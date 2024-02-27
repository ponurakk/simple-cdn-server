package api

import (
	"encoding/json"
	"os"
)

type Config struct {
	Token string `json:"token"`
	Port  string `json:"port"`
}

func ReadConfig() (Config, error) {
	file, err := os.Open("config.json")
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	var config Config
	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}
