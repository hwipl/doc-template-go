package doctempl

import (
	"encoding/json"
	"os"
)

// ConfigTemplate is a template configuration in Config.
type ConfigTemplate struct {
	File     string
	Output   string
	DataFile string
	Data     map[string]any
}

// Config is a document template configuration.
type Config struct {
	Templates []*ConfigTemplate
}

// NewConfig returns a new Config.
func NewConfig() *Config {
	return &Config{}
}

// LoadConfig loads the configuration from file.
func LoadConfig(file string) (*Config, error) {
	f, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	c := NewConfig()
	if err := json.Unmarshal(f, c); err != nil {
		return nil, err
	}
	return c, nil
}
