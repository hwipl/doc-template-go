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

// LoadConfig loads a Config from file.
func LoadConfig(file string) (*Config, error) {
	f, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	c := &Config{}
	if err := json.Unmarshal(f, c); err != nil {
		return nil, err
	}
	return c, nil
}
