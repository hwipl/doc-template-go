package doctempl

import (
	"encoding/json"
	"os"
)

// ConfigTemplate is a template configuration in Config.
type ConfigTemplate struct {
	File     string
	Output   string
	Data     map[string]any
	DataFile string
}

// Config is a document template configuration.
type Config struct {
	Templates []*ConfigTemplate
}

// Load loads the configuration from file.
func (c *Config) Load(file string) error {
	f, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	return json.Unmarshal(f, c)
}

// NewConfig returns a new Config.
func NewConfig() *Config {
	return &Config{}
}
