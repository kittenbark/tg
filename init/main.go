package main

import (
	"embed"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

//go:embed all:templates
var templatesFS embed.FS

func main() {
	var module string
	var overwrite bool

	flag.StringVar(&module, "module", "", "go module name (defaults to current directory name)")
	flag.BoolVar(&overwrite, "overwrite", false, "overwrite existing files")
	flag.Parse()

	if module == "" {
		abs, err := filepath.Abs(".")
		if err != nil {
			fatal(err)
		}
		module = filepath.Base(abs)
	}

	files := []struct {
		tmpl string
		dest string
	}{
		{"templates/main.go.tmpl", "cmd/main.go"},
		{"templates/Dockerfile.tmpl", "Dockerfile"},
		{"templates/compose.yml.tmpl", "compose.yml"},
		{"templates/go.mod.tmpl", "go.mod"},
		{"templates/.env.tmpl", ".env"},
		{"templates/.gitignore.tmpl", ".gitignore"},
	}

	name := filepath.Base(module)

	goModCreated := false
	for _, f := range files {
		data, err := templatesFS.ReadFile(f.tmpl)
		if err != nil {
			fatal(fmt.Errorf("read template %s: %w", f.tmpl, err))
		}
		content := strings.NewReplacer("{{module}}", module, "{{name}}", name).Replace(string(data))
		created, err := writeFile(f.dest, content, overwrite)
		if err != nil {
			fatal(err)
		}
		if f.dest == "go.mod" && created {
			goModCreated = true
		}
	}

	if goModCreated {
		fmt.Println()
		cmd := exec.Command("go", "mod", "tidy")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "warning: go mod tidy failed: %v\n", err)
			fmt.Fprintln(os.Stderr, "run 'go mod tidy' manually to generate go.sum")
		}
	}

	fmt.Printf("\ninitialized %s\n", module)
	fmt.Println("set KITTENBARK_TG_TOKEN in .env, then: docker compose up --build")
}

func writeFile(path, content string, overwrite bool) (created bool, err error) {
	if !overwrite {
		if _, err := os.Stat(path); err == nil {
			fmt.Printf("skip   %s\n", path)
			return false, nil
		}
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return false, err
	}
	fmt.Printf("create %s\n", path)
	return true, os.WriteFile(path, []byte(content), 0644)
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
