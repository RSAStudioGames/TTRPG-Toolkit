package server

import (
	"io/fs"
	"net/http"
	"os"
	"strings"

	"github.com/gabriel/ttrpg-toolkit/backend/internal/api"
	"github.com/gabriel/ttrpg-toolkit/backend/internal/config"
	"github.com/gabriel/ttrpg-toolkit/backend/internal/repository"
	"github.com/gabriel/ttrpg-toolkit/backend/internal/services"
	"github.com/gabriel/ttrpg-toolkit/backend/ui"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// New returns a configured Fiber application and a cleanup function.
func New(cfg config.Config) (*fiber.App, func(), error) {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: false,
	})

	app.Use(logger.New())

	var cleanup func() = func() {}
	deps := api.Deps{UploadDir: cfg.UploadDir}

	if cfg.DatabaseURL != "" {
		db, err := repository.Connect(cfg.DatabaseURL)
		if err != nil {
			return nil, nil, err
		}
		cleanup = func() { _ = db.Close() }
		if err := os.MkdirAll(cfg.UploadDir, 0o755); err != nil {
			_ = db.Close()
			return nil, nil, err
		}
		repo := repository.NewSystemRepository(db)
		deps.Systems = services.NewSystemService(repo, cfg.UploadDir)
		mechanicsRepo := repository.NewMechanicsRepository(db)
		deps.Formula = services.NewFormulaService()
		deps.Mechanics = services.NewMechanicsService(mechanicsRepo, repo, deps.Formula)
	}

	staticFS, err := fs.Sub(ui.Static, "static")
	if err != nil {
		cleanup()
		return nil, nil, err
	}

	api.Register(app, deps)

	app.Use("/", filesystem.New(filesystem.Config{
		Root:         http.FS(staticFS),
		Browse:       false,
		Index:        "index.html",
		NotFoundFile: "index.html",
		Next: func(c *fiber.Ctx) bool {
			return strings.HasPrefix(c.Path(), "/api")
		},
	}))

	return app, cleanup, nil
}
