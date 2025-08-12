package doctempl

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

// TestLoadConfig tests LoadConfig.
func TestLoadConfig(t *testing.T) {
	d := t.TempDir()

	// not existing file
	f := filepath.Join("does not exist")
	if _, err := LoadConfig(f); err == nil {
		t.Error("not existing file should return error")
	}

	// empty file
	f = filepath.Join(d, "config.json")
	if err := os.WriteFile(f, []byte(""), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadConfig(f); err == nil {
		t.Error("empty file should return error")
	}

	// valid file
	data := []byte(`{
	"Templates": [
		{
			"File": "template1.tmpl",
			"Output": "output1.txt",
			"DataFile": "datafile1.json"
		},
		{
			"File": "template2.tmpl"
		}
	]
	}`)
	want := &Config{
		Templates: []*ConfigTemplate{
			{
				File:     "template1.tmpl",
				Output:   "output1.txt",
				DataFile: "datafile1.json",
			},
			{
				File: "template2.tmpl",
			},
		},
	}
	if err := os.WriteFile(f, data, 0600); err != nil {
		t.Fatal(err)
	}
	c, err := LoadConfig(f)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(c, want) {
		t.Errorf("got %v, want %v", *c, *want)
	}
}
