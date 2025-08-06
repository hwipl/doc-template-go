// Package doctempl contains document templates.
package doctempl

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
)

// parseJSON converts the json in b to a map.
func parseJSON(b []byte) (map[string]any, error) {
	m := map[string]any{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return m, nil
}

// parseJSONFile converts the json in file to a map.
func parseJSONFile(file string) (map[string]any, error) {
	b, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return parseJSON(b)
}

// parseJSONArg converts the json in arg to a map.
func parseJSONArg(arg string) (map[string]any, error) {
	return parseJSON([]byte(arg))
}

// parseArgs converts command line arguments in args to a map.
func parseArgs(args []string) (map[string]any, error) {
	m := map[string]any{}

	for _, arg := range args {
		s := strings.SplitN(arg, "=", 2)
		if len(s) != 2 {
			continue
		}
		key := s[0]
		val := s[1]

		// json
		var j any
		if err := json.Unmarshal([]byte(val), &j); err == nil {
			m[key] = j
			continue
		}

		// lists
		if strings.HasPrefix(val, "[") && strings.HasSuffix(val, "]") {
			l := val[1 : len(val)-1]
			m[key] = strings.Split(l, ",")
			continue
		}

		// maps
		if strings.HasPrefix(val, "{") && strings.HasSuffix(val, "}") {
			kv := map[string]string{}
			s := val[1 : len(val)-1]
			for pair := range strings.SplitSeq(s, ",") {
				p := strings.SplitN(pair, ":", 2)
				if len(p) == 2 {
					k := p[0]
					v := p[1]
					kv[k] = v
				}
			}
			m[key] = kv
			continue
		}

		m[key] = val
	}
	return m, nil
}

// runTemplate runs the template in file with data and writes to the file in
// output or to Stdout if output is empty.
func runTemplate(file string, data any, output string) (err error) {
	t, err := NewDocTemplate(file)
	if err != nil {
		return err
	}

	if output == "" {
		return t.Execute(os.Stdout, data)
	}

	f, err := os.OpenFile(output, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		return err
	}
	defer func() {
		err = f.Close()
	}()
	return t.Execute(f, data)
}

// runTemplates runs the templates in config.
func runTemplates(config *Config) error {
	for _, tmpl := range config.Templates {
		if err := runTemplate(tmpl.File, tmpl.Data, tmpl.Output); err != nil {
			return fmt.Errorf("error executing template %s: %w",
				tmpl.File, err)
		}
	}
	return nil
}

// flagIsSet returns whether command line flag is set in flag set.
func flagIsSet(fs *flag.FlagSet, name string) bool {
	isSet := false
	fs.Visit(func(f *flag.Flag) {
		if name == f.Name {
			isSet = true
		}
	})
	return isSet
}

// getConfig returns a config with fields loaded from command line arguments in
// args and from a config file.
func getConfig(args []string) (*Config, error) {
	// command line arguments
	fs := flag.NewFlagSet(args[0], flag.ExitOnError)
	flagConfig := fs.String("config", ".doc-template-go.json",
		"read configuration from `file`")
	flagFile := fs.String("file", "", "read template from `file`")
	flagOutput := fs.String("output", "", "write output to `file`")
	flagDataFile := fs.String("data-file", "", "load data from json `file`")
	flagData := fs.String("data", "", "read input data from `json`")
	_ = fs.Parse(args[1:])

	// read config file
	config := NewConfig()
	config.ConfigFile = *flagConfig
	err := config.Load()
	if err != nil {
		if flagIsSet(fs, "config") {
			return nil, fmt.Errorf("error loading config file: %w", err)
		}
	}

	// template file argument
	if flagIsSet(fs, "file") {
		config.File = *flagFile

		config.Templates = []*ConfigTemplate{
			{File: config.File},
		}
	}
	if len(config.Templates) == 0 {
		fs.Usage()
		return nil, fmt.Errorf("no templates found")
	}

	// data file
	if flagIsSet(fs, "data-file") {
		config.DataFile = *flagDataFile

		// data file from command line arguments
		data, err := parseJSONFile(config.DataFile)
		if err != nil {
			return nil, fmt.Errorf("error parsing data file: %w", err)
		}
		for _, tmpl := range config.Templates {
			tmpl.Data = data
		}
	} else {
		// data file from config
		for _, tmpl := range config.Templates {
			if tmpl.DataFile == "" {
				continue
			}
			data, err := parseJSONFile(tmpl.DataFile)
			if err != nil {
				return nil, fmt.Errorf("error parsing data file: %w", err)
			}
			tmpl.Data = data
		}
	}

	// data arguments
	data, err := parseArgs(fs.Args())
	if err != nil {
		return nil, fmt.Errorf("error parsing arguments: %w", err)
	}

	if flagIsSet(fs, "data") {
		config.DataString = *flagData

		data, err = parseJSONArg(config.DataString)
		if err != nil {
			return nil, fmt.Errorf("error parsing json: %w", err)
		}
	}

	// set output in config
	if flagIsSet(fs, "output") {
		config.Output = *flagOutput

		for _, tmpl := range config.Templates {
			tmpl.Output = config.Output
		}
	}

	// set data in config
	if len(data) != 0 {
		for _, tmpl := range config.Templates {
			tmpl.Data = data
		}
	}

	return config, nil
}

// Run is the main entry point.
func Run() {
	// get config
	config, err := getConfig(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting configuration: %v\n", err)
		os.Exit(1)
	}

	// run templates in config
	if err := runTemplates(config); err != nil {
		fmt.Fprintf(os.Stderr, "Error running templates: %v\n", err)
		os.Exit(1)
	}
}
