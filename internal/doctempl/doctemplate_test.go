package doctempl

import (
	"os"
	"path/filepath"
	"testing"
)

// TestDocTemplateExecute tests Execute of DocTemplate.
func TestDocTemplateExecute(t *testing.T) {
	// create doc template
	f := filepath.Join(t.TempDir(), "template.tmpl")
	if err := os.WriteFile(f, []byte(""), 0600); err != nil {
		t.Fatal(err)
	}
	tmpl, err := NewDocTemplate(f)
	if err != nil {
		t.Fatal(err)
	}

	// execute
	if err := tmpl.Execute(os.Stdout, nil); err != nil {
		t.Error(err)
	}
}

// TestNewDocTemplate tests NewDocTemplate.
func TestNewDocTemplate(t *testing.T) {
	d := t.TempDir()

	// not existing file
	if _, err := NewDocTemplate(filepath.Join(d, "does not exist")); err == nil {
		t.Error("not existing file should return error")
	}

	// existing file
	f := filepath.Join(d, "template.tmpl")
	if err := os.WriteFile(f, []byte(""), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := NewDocTemplate(f); err != nil {
		t.Error(err)
	}
}
