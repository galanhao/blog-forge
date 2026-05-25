package theme

import (
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// Engine renders pages using Go templates from the theme directory.
type Engine struct {
	theme *Theme
	tpl   *template.Template
}

// NewEngine creates a template engine for the given theme.
func NewEngine(t *Theme) (*Engine, error) {
	e := &Engine{theme: t}

	tpl := template.New("").Funcs(defaultFuncs())

	// parse all layout templates (base.html provides page-start/page-end defines)
	if err := parseTemplateDir(tpl, filepath.Join(t.Dir, "layouts")); err != nil {
		return nil, fmt.Errorf("theme engine: parse layouts: %w", err)
	}

	e.tpl = tpl
	return e, nil
}

// Execute renders a layout template with the given data.
func (e *Engine) Execute(layout string, data any) (string, error) {
	var sb strings.Builder
	if err := e.tpl.ExecuteTemplate(&sb, layout, data); err != nil {
		return "", fmt.Errorf("theme engine: execute %s: %w", layout, err)
	}
	return sb.String(), nil
}

// parseTemplateDir walks a directory and adds all .html files as named templates.
func parseTemplateDir(tpl *template.Template, dir string) error {
	return filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		if filepath.Ext(path) != ".html" {
			return nil
		}

		data, readErr := os.ReadFile(path)
		if readErr != nil {
			return readErr
		}

		name := filepath.Base(path)
		if _, execErr := tpl.New(name).Parse(string(data)); execErr != nil {
			return fmt.Errorf("parse %s: %w", path, execErr)
		}

		return nil
	})
}
