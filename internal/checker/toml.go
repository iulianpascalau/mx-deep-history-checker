package checker

import (
	"fmt"
	"os"

	"github.com/pelletier/go-toml/v2"
)

// DbConfig represents the expected structure of config.toml in a sharded DB folder.
type DbConfig struct {
	NumShards uint   `toml:"NumShards"`
	Type      string `toml:"Type"`
}

// readDBConfig detects and parses a config.toml file if it exists.
func readDBConfig(tomlPath string) (*DbConfig, error) {
	data, err := os.ReadFile(tomlPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil // Not an error if it doesn't exist
		}
		return nil, fmt.Errorf("failed to read %s: %w", tomlPath, err)
	}

	var cfg DbConfig
	err = toml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to parse %s: %w", tomlPath, err)
	}

	return &cfg, nil
}
