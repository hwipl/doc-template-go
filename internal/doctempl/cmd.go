// Package doctempl contains document templates.
package doctempl

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"
)

// DocTemplate is a document template.
type DocTemplate struct {
	File     string
	Template *template.Template
}

// Execute runs the template with data and writes output to wr.
func (d *DocTemplate) Execute(wr io.Writer, data any) error {
	return d.Template.Execute(wr, data)
}

// NewDocTemplate returns a new DocTemplate.
func NewDocTemplate(file string) (*DocTemplate, error) {
	tmpl, err := template.ParseFiles(file)
	if err != nil {
		return nil, err
	}

	return &DocTemplate{file, tmpl}, nil
}

// parseJSON converts the json in s to a map.
func parseJSON(s string) (map[string]any, error) {
	m := map[string]any{}
	if err := json.Unmarshal([]byte(s), &m); err != nil {
		return nil, err
	}
	return m, nil
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

// runTemplateStdout runs the template in file with data and writes to Stdout.
func runTemplateStdout(file string, data any) error {
	t, err := NewDocTemplate(file)
	if err != nil {
		return err
	}

	return t.Execute(os.Stdout, data)
}

// Run is the main entry point.
func Run() {
	// command line arguments
	file := flag.String("file", "", "read template from `file`")
	json := flag.String("json", "", "read input data from `json`")

	flag.Parse()

	if *file == "" {
		flag.Usage()
		os.Exit(1)
	}

	data, err := parseArgs(flag.Args())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing arguments: %v\n", err)
		os.Exit(1)
	}

	if *json != "" {
		data, err = parseJSON(*json)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing json: %v\n", err)
			os.Exit(1)
		}
	}

	if err := runTemplateStdout(*file, data); err != nil {
		fmt.Fprintf(os.Stderr, "Error executing template: %v\n", err)
		os.Exit(1)
	}
}
