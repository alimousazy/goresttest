package chatserver

import (
	"encoding/json"
	"os"
)

type Config struct {
	Host           string
	ChatPort       uint16
	WebServicePort uint16
	MaxClient      int64
	LogPath        string
}

func ReadConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration := Config{
		Host:           "localhost",
		ChatPort:       6000,
		WebServicePort: 8080,
		LogPath:        "./log_me",
	}
	err = decoder.Decode(&configuration)
	if err != nil {
		return nil, err
	}
	return &configuration, nil
}
