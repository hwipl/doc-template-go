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
	ConfigFile string `json:"-"`
	DataString string `json:"-"`
	Templates  []*ConfigTemplate
}

// Load loads the configuration from file.
func (c *Config) Load() error {
	f, err := os.ReadFile(c.ConfigFile)
	if err != nil {
		return err
	}
	return json.Unmarshal(f, c)
}

// NewConfig returns a new Config.
func NewConfig() *Config {
	return &Config{}
}
