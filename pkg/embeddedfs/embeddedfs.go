// Package embeddedfs provides utility functions for reading files and rendering templates from an embedded file system.
package embeddedfs

import (
	"bytes"
	"embed"
	"fmt"
	"io"
	"os"
	"text/tabwriter"
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

// EmbeddedFS wraps an embedded file system.
type EmbeddedFS struct {
	Dirname string
	FS      embed.FS
}

// NewEmbeddedFS creates a new EmbeddedFS instance.
func NewEmbeddedFS(dirname string, fs embed.FS) *EmbeddedFS {
	return &EmbeddedFS{
		Dirname: dirname,
		FS:      fs,
	}
}

// ReadFile reads a file from the embedded file system.
func (fs *EmbeddedFS) ReadFile(dir, filename string) ([]byte, error) {
	return fs.FS.ReadFile(fs.toEmbeddedFilePath(dir, filename))
}

// RenderTemplate renders a template from the embedded file system and writes it to disk.
func (fs *EmbeddedFS) RenderTemplate(args interface{}, dir, filename, outputPath string) error {
	data, err := fs.RenderTemplateBytes(args, dir, filename)
	if err != nil {
		return err
	}
	return writeFile(outputPath, data)
}

// RenderTemplateBytes renders a template from the embedded file system and returns the resulting bytes.
func (fs *EmbeddedFS) RenderTemplateBytes(args interface{}, dir, filename string) ([]byte, error) {
	var writer bytes.Buffer
	if err := fs.render(args, &writer, dir, filename); err != nil {
		return nil, err
	}
	return writer.Bytes(), nil
}

// render renders a template from the embedded file system.
func (fs *EmbeddedFS) render(args interface{}, writer *bytes.Buffer, dir, filename string) error {
	tfm := sprig.TxtFuncMap() // include sprig library functions
	tmpl, err := template.New(filename).Funcs(tfm).ParseFS(fs.FS, fs.toEmbeddedFilePath(dir, filename))
	if err != nil {
		return err
	}
	if err := tmpl.Option("missingkey=zero").Execute(writer, args); err != nil {
		return fmt.Errorf("failed to parse template: args: %+v. Error: %w", args, err)
	}
	return nil
}

// PrintTableTemplate renders a tabular data template and flushes it to stdout.
func (fs *EmbeddedFS) PrintTableTemplate(out io.Writer, args interface{}, dir, filename string) error {
	tabw := tabwriter.NewWriter(out, 8, 8, 8, ' ', 0)
	tmpl, err := template.ParseFS(fs.FS, fs.toEmbeddedFilePath(dir, filename))
	if err != nil {
		return err
	}
	if err := tmpl.Execute(tabw, args); err != nil {
		return fmt.Errorf("failed to parse template: args: %+v. Error: %w", args, err)
	}
	return tabw.Flush()
}

// toEmbeddedFilePath retrieves the full path of a file within the embedded file system.
// Note that filepath.Join is NOT used here, as embed requires the '/' separator.
func (fs *EmbeddedFS) toEmbeddedFilePath(dir, filename string) string {
	if dir == "" {
		return fmt.Sprintf("%s/%s", fs.Dirname, filename)
	}
	return fmt.Sprintf("%s/%s/%s", fs.Dirname, dir, filename)
}

// writeFile writes bytes to the specified file.
func writeFile(outFilename string, data []byte) error {
	if err := os.WriteFile(outFilename, data, 0600); err != nil {
		return fmt.Errorf("failed to write rendered template to disk: %w", err)
	}
	return nil
}
