package scaffold

import (
	"embed"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
)

type BuildScaffoldOptions struct {
	Dir   string
	Force bool
}

//go:embed templates/*
var templatesFS embed.FS

func BuildScaffold(options BuildScaffoldOptions) error {
	if _, err := os.Stat(options.Dir); os.IsNotExist(err) {
		if err := os.MkdirAll(options.Dir, 0755); err != nil {
			return fmt.Errorf("failed to create target dir: %w", err)
		}
	}

	// Define files to generate
	files := []struct {
		Name     string
		Template string
	}{
		{"zyra.config", "templates/zyra.config.tmpl"},
		{filepath.Join("requests", "example.zyra"), "templates/example.zyra.tmpl"},
		{".gitignore", "templates/gitignore.tmpl"},
	}

	for _, f := range files {
		targetPath := filepath.Join(options.Dir, f.Name)

		if _, err := os.Stat(targetPath); err == nil && !options.Force {
			fmt.Printf("⚠ %s already exists, skipping\n", f.Name)
			continue
		}

		tmplContent, err := templatesFS.ReadFile(f.Template)
		if err != nil {
			return fmt.Errorf("failed to read template %s: %w", f.Template, err)
		}

		tmpl, err := template.New(f.Name).Parse(string(tmplContent))
		if err != nil {
			return fmt.Errorf("failed to parse template %s: %w", f.Template, err)
		}

		if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
			return fmt.Errorf("failed to create dir for %s: %w", targetPath, err)
		}

		file, err := os.Create(targetPath)
		if err != nil {
			return fmt.Errorf("failed to create file %s: %w", targetPath, err)
		}

		if err := tmpl.Execute(file, nil); err != nil {
			file.Close()
			return fmt.Errorf("failed to execute template %s: %w", f.Name, err)
		}

		file.Close()
		fmt.Printf("✔ %s created\n", f.Name)
	}

	return nil
}
