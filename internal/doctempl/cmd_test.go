package doctempl

import (
	"path/filepath"
	"reflect"
	"testing"
)

// TestParseArgs tests parseArgs.
func TestParseArgs(t *testing.T) {
	for i, args := range []struct {
		v    []string
		want map[string]any
	}{
		// test invalid
		{
			[]string{"invalid"},
			map[string]any{},
		},
		// test empty
		{
			[]string{},
			map[string]any{},
		},
		// test string
		{
			[]string{"string=string"},
			map[string]any{"string": "string"},
		},
		// test json
		{
			[]string{"list=[]"},
			map[string]any{"list": []any{}},
		},
		{
			[]string{"list=[\"s1\",\"s2\"]"},
			map[string]any{"list": []any{"s1", "s2"}},
		},
		{
			[]string{"object={}"},
			map[string]any{"object": map[string]any{}},
		},
		{
			[]string{"object={\"o1\": \"s1\", \"o2\": \"s2\"}"},
			map[string]any{"object": map[string]any{"o1": "s1", "o2": "s2"}},
		},
		// test list
		{
			[]string{"list=[s1,s2]"},
			map[string]any{"list": []string{"s1", "s2"}},
		},
		// test map
		{
			[]string{"map={o1:s1,o2:s2}"},
			map[string]any{"map": map[string]string{"o1": "s1", "o2": "s2"}},
		},
	} {
		got, err := parseArgs(args.v)
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(got, args.want) {
			t.Errorf("%d: got %v, want %v", i, got, args.want)
		}
	}
}

// TestGetConfig tests getConfig.
func TestGetConfig(t *testing.T) {
	// files
	dir := t.TempDir()
	config := filepath.Join(dir, ".doc-template-go.json")
	tmpl := filepath.Join(dir, "template.tmpl")
	data := filepath.Join(dir, "data")

	// no config, no template
	if _, err := getConfig([]string{"test", "-config", config}); err == nil {
		t.Error("no templates should return error")
	}

	// no config, not existing template
	if _, err := getConfig([]string{"test", "-config", config, "-file", tmpl}); err != nil {
		t.Error(err)
	}

	// no config, not existing template, not existing data file
	if _, err := getConfig([]string{"test", "-config", config, "-file", tmpl, "-data-file", data}); err == nil {
		t.Error("not existing data file should return error")
	}
}
