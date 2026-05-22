package repository

import (
	"embed"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
)

//go:embed migrations/*.up.sql
var migrationFS embed.FS

// RunMigrations applies pending numbered migrations.
func RunMigrations(db *sqlx.DB) error {
	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INT PRIMARY KEY,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`); err != nil {
		return fmt.Errorf("schema_migrations table: %w", err)
	}

	entries, err := migrationFS.ReadDir("migrations")
	if err != nil {
		return fmt.Errorf("read migrations dir: %w", err)
	}

	var versions []int
	files := make(map[int]string)
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".up.sql") {
			continue
		}
		parts := strings.SplitN(e.Name(), "_", 2)
		v, err := strconv.Atoi(parts[0])
		if err != nil {
			continue
		}
		versions = append(versions, v)
		b, err := migrationFS.ReadFile("migrations/" + e.Name())
		if err != nil {
			return err
		}
		files[v] = string(b)
	}
	sort.Ints(versions)

	for _, v := range versions {
		var count int
		_ = db.Get(&count, `SELECT COUNT(*) FROM schema_migrations WHERE version = $1`, v)
		if count > 0 {
			continue
		}

		tx, err := db.Beginx()
		if err != nil {
			return err
		}
		if _, err := tx.Exec(files[v]); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("migration %d: %w", v, err)
		}
		if _, err := tx.Exec(`INSERT INTO schema_migrations (version) VALUES ($1)`, v); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("record migration %d: %w", v, err)
		}
		if err := tx.Commit(); err != nil {
			return err
		}
	}
	return nil
}
