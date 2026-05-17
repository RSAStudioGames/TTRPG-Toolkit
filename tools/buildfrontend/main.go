// buildfrontend compiles the SvelteKit app and copies assets into backend/ui/static
// for go:embed. Invoked via go generate; not part of the runtime binary.
package main

import (
	"fmt"
	"io"
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

	frontend := filepath.Join(root, "frontend")
	staticOut := filepath.Join(root, "build", "static")
	embedDir := filepath.Join(root, "backend", "ui", "static")

	if _, err := os.Stat(filepath.Join(frontend, "package.json")); err != nil {
		fatal(fmt.Errorf("frontend/package.json not found: %w", err))
	}

	npm := "npm"
	if runtime.GOOS == "windows" {
		npm = "npm.cmd"
	}

	fmt.Println("buildfrontend: npm run build")
	cmd := exec.Command(npm, "run", "build")
	cmd.Dir = frontend
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fatal(fmt.Errorf("npm run build: %w", err))
	}

	if _, err := os.Stat(filepath.Join(staticOut, "index.html")); err != nil {
		fatal(fmt.Errorf("missing %s — adapter-static output not found", staticOut))
	}

	fmt.Println("buildfrontend: copy to backend/ui/static")
	if err := os.RemoveAll(embedDir); err != nil {
		fatal(err)
	}
	if err := os.MkdirAll(embedDir, 0o755); err != nil {
		fatal(err)
	}
	if err := copyDir(staticOut, embedDir); err != nil {
		fatal(err)
	}

	fmt.Println("buildfrontend: done")
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

func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)
		if info.IsDir() {
			return os.MkdirAll(target, info.Mode())
		}
		return copyFile(path, target)
	})
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "buildfrontend: %v\n", err)
	os.Exit(1)
}
