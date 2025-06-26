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

// Run is the main entry point.
func Run() {
	// create config
	config := NewConfig()

	// command line arguments
	flag.StringVar(&config.ConfigFile, "config", ".doc-template-go.json",
		"read configuration from `file`")
	flag.StringVar(&config.File, "file", "", "read template from `file`")
	flag.StringVar(&config.Output, "output", "", "write output to `file`")
	flag.StringVar(&config.DataFile, "data-file", "", "load data from `file`")
	flag.StringVar(&config.DataString, "data", "", "read input data from `json`")
	flag.Parse()

	// read config file
	err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not load config file: %v\n", err)
		config = &Config{}
	}

	// template file argument
	if config.File != "" {
		config.Templates = []*ConfigTemplate{
			{File: config.File},
		}
	}
	if len(config.Templates) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	// data file
	if config.DataFile != "" {
		// data file from command line arguments
		data, err := parseJSONFile(config.DataFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing data file: %v\n", err)
			os.Exit(1)
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
				fmt.Fprintf(os.Stderr, "Error parsing data file: %v\n", err)
				os.Exit(1)
			}
			tmpl.Data = data
		}
	}

	// data arguments
	data, err := parseArgs(flag.Args())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing arguments: %v\n", err)
		os.Exit(1)
	}

	if config.DataString != "" {
		data, err = parseJSONArg(config.DataString)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing json: %v\n", err)
			os.Exit(1)
		}
	}

	// set output in config
	if config.Output != "" {
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

	if err := runTemplates(config); err != nil {
		fmt.Fprintf(os.Stderr, "Error running templates: %v\n", err)
		os.Exit(1)
	}
}
