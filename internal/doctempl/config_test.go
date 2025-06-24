package doctempl

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

// TestConfigLoad tests Load of Config.
func TestConfigLoad(t *testing.T) {
	d := t.TempDir()

	// not existing file
	c := NewConfig()
	c.ConfigFile = filepath.Join("does not exist")
	if err := c.Load(); err == nil {
		t.Error("not existing file should return error")
	}

	// empty file
	c = NewConfig()
	f := filepath.Join(d, "config.json")
	c.ConfigFile = f
	if err := os.WriteFile(f, []byte(""), 0600); err != nil {
		t.Fatal(err)
	}
	if err := c.Load(); err == nil {
		t.Error("empty file should return error")
	}

	// valid file
	c = NewConfig()
	c.ConfigFile = f
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
		ConfigFile: f,
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
	err := c.Load()
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(c, want) {
		t.Errorf("got %v, want %v", *c, *want)
	}
}
