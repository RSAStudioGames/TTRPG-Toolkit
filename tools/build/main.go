// build runs full prep: npm ci, frontend embed, and ttrpg-toolkit at repo root.
package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func main() {
	root, err := repoRoot()
	if err != nil {
		fatal(err)
	}

	toolsDir := filepath.Join(root, "tools")
	frontend := filepath.Join(root, "frontend")

	npm := "npm"
	if runtime.GOOS == "windows" {
		npm = "npm.cmd"
	}

	fmt.Println("build: npm ci")
	if err := run(npm, frontend, "ci"); err != nil {
		fatal(fmt.Errorf("npm ci: %w", err))
	}

	fmt.Println("build: frontend embed")
	if err := run("go", toolsDir, "run", "./buildfrontend"); err != nil {
		fatal(fmt.Errorf("buildfrontend: %w", err))
	}

	out := filepath.Join(root, "ttrpg-toolkit")
	if runtime.GOOS == "windows" {
		out += ".exe"
	}

	backend := filepath.Join(root, "backend")
	fmt.Printf("build: go build -> %s\n", out)
	if err := run("go", backend, "build", "-o", out, "./cmd"); err != nil {
		fatal(fmt.Errorf("go build: %w", err))
	}

	fmt.Println("build: done")
}

func repoRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if fileExists(filepath.Join(dir, "backend", "go.mod")) && fileExists(filepath.Join(dir, "frontend", "package.json")) {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("repository root not found")
		}
		dir = parent
	}
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func run(name, dir string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	return cmd.Run()
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "build: %v\n", err)
	os.Exit(1)
}
