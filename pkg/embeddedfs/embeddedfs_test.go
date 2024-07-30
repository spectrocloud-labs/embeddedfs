package embeddedfs

import (
	"embed"
	"os"
	"strings"
	"testing"
)

//go:embed testdata/*
var testFS embed.FS

func TestReadFile(t *testing.T) {
	fs := NewEmbeddedFS("testdata", testFS)

	t.Run("Success", func(t *testing.T) {
		data, err := fs.ReadFile("", "test.txt")
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if string(data) != "Hello, World!" {
			t.Errorf("expected data to be 'Hello, World!', got %s", string(data))
		}
	})

	t.Run("Error", func(t *testing.T) {
		_, err := fs.ReadFile("", "nonexistent.txt")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestRenderTemplate(t *testing.T) {
	fs := NewEmbeddedFS("testdata", testFS)

	t.Run("Success", func(t *testing.T) {
		args := struct {
			Name string
		}{
			Name: "Bob",
		}
		if err := fs.RenderTemplate(args, "", "test.tmpl", "test.out"); err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		defer os.Remove("test.out")
		data, err := os.ReadFile("test.out")
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if string(data) != "Name: Bob" {
			t.Errorf("expected data to be 'Name: Bob', got %s", string(data))
		}
	})

	t.Run("Error", func(t *testing.T) {
		err := fs.RenderTemplate(nil, "", "nonexistent.tmpl", "test.txt")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestRenderTemplateBytes(t *testing.T) {
	fs := NewEmbeddedFS("testdata", testFS)

	t.Run("Success", func(t *testing.T) {
		args := struct {
			Name string
		}{
			Name: "Alice",
		}
		data, err := fs.RenderTemplateBytes(args, "", "test.tmpl")
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if string(data) != "Name: Alice" {
			t.Errorf("expected data to be 'Name: Alice', got %s", string(data))
		}
	})

	t.Run("Error", func(t *testing.T) {
		_, err := fs.RenderTemplateBytes(nil, "", "nonexistent.tmpl")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestPrintTableTemplate(t *testing.T) {
	fs := NewEmbeddedFS("testdata", testFS)

	t.Run("Success", func(t *testing.T) {
		args := []struct {
			Name string
			Age  string
		}{
			{
				Name: "Alice",
				Age:  "30",
			},
			{
				Name: "Bob",
				Age:  "40",
			},
		}
		sb := &strings.Builder{}
		if err := fs.PrintTableTemplate(sb, args, "", "table.tmpl"); err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		data := sb.String()
		expected := `Name         Age
Alice        30
Bob          40`
		if string(data) != expected {
			t.Errorf("expected data to be %q, got %s", expected, data)
		}
	})

	t.Run("Error", func(t *testing.T) {
		err := fs.RenderTemplate(nil, "", "nonexistent.tmpl", "test.txt")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}
