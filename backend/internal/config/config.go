package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/joho/godotenv"
)

const (
	defaultHost = "0.0.0.0"
	defaultPort = 8080
)

// Config holds server configuration resolved from defaults, .env, and OS env.
type Config struct {
	Host string
	Port int
}

// Load reads configuration: compiled defaults, then .env at repo root, then OS env.
func Load() (Config, error) {
	_ = loadDotEnv()

	host := envOrDefault("TTRPG_SERVER_HOST", defaultHost)
	portStr := envOrDefault("TTRPG_SERVER_PORT", strconv.Itoa(defaultPort))
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return Config{}, fmt.Errorf("invalid TTRPG_SERVER_PORT %q: %w", portStr, err)
	}
	if port < 1 || port > 65535 {
		return Config{}, fmt.Errorf("TTRPG_SERVER_PORT out of range: %d", port)
	}

	return Config{Host: host, Port: port}, nil
}

func loadDotEnv() error {
	root, err := findRepoRoot()
	if err != nil {
		return err
	}
	envPath := filepath.Join(root, ".env")
	if _, err := os.Stat(envPath); os.IsNotExist(err) {
		return nil
	}
	return godotenv.Load(envPath)
}

func findRepoRoot() (string, error) {
	if wd, err := os.Getwd(); err == nil {
		if fileExists(filepath.Join(wd, ".env.example")) || fileExists(filepath.Join(wd, "frontend")) {
			return wd, nil
		}
		if fileExists(filepath.Join(wd, "..", ".env.example")) {
			return filepath.Clean(filepath.Join(wd, "..")), nil
		}
	}

	exe, err := os.Executable()
	if err != nil {
		return ".", nil
	}
	dir := filepath.Dir(exe)
	for i := 0; i < 6; i++ {
		if fileExists(filepath.Join(dir, ".env.example")) || fileExists(filepath.Join(dir, "frontend")) {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return ".", nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// Addr returns the listen address host:port.
func (c Config) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
