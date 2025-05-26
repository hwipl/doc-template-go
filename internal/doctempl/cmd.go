// Package doctempl contains document templates.
package doctempl

import (
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

// parseArgs converts command line arguments in args to a map.
func parseArgs(args []string) (map[string]any, error) {
	m := map[string]any{}
	for _, arg := range args {
		s := strings.SplitN(arg, "=", 2)
		if len(s) != 2 {
			continue
		}
		m[s[0]] = s[1]
	}
	return m, nil
}

// Run is the main entry point.
func Run() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Command line arguments missing")
		os.Exit(1)
	}
	args := os.Args[2:]

	_, err := parseArgs(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing arguments: %v\n", err)
		os.Exit(1)
	}
}
