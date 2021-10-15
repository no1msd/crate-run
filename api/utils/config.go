package utils

import (
	"fmt"
	"io"
	"io/ioutil"

	"github.com/pelletier/go-toml"
)

// Config holds the configuration information of the application.
type Config struct {
	Listener struct {
		Address string
		Port    uint
	}

	Database struct {
		Path string
	}
}

// LoadOrCreateConfig loads a configuration from the given ReadWriter and returns the parsed
// Config object. Missing values are substituted with default values, and written back to the
// file. If no configuration file existed, it creates one with all default values.
func LoadOrCreateConfig(rw io.ReadWriter) (*Config, error) {
	config := Config{}
	config.Listener.Address = "0.0.0.0"
	config.Listener.Port = 9080
	config.Database.Path = "crateapi.db3"

	data, err := ioutil.ReadAll(rw)
	if err != nil {
		return nil, fmt.Errorf("could not read config: %w", err)
	}

	err = toml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal config: %w", err)
	}

	newData, err := toml.Marshal(config)
	if err != nil {
		return &config, fmt.Errorf("could not marshal config: %w", err)
	}

	if string(newData) != string(data) {
		_, err := rw.Write(newData)
		if err != nil {
			return &config, fmt.Errorf("could not write config: %w", err)
		}
	}

	return &config, nil
}
