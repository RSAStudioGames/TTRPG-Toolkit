package server

import (
	"io/fs"
	"net/http"
	"strings"

	"github.com/gabriel/ttrpg-toolkit/backend/internal/api"
	"github.com/gabriel/ttrpg-toolkit/backend/ui"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// New returns a configured Fiber application.
func New() *fiber.App {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: false,
	})

	app.Use(logger.New())

	staticFS, err := fs.Sub(ui.Static, "static")
	if err != nil {
		panic("embedded static files: " + err.Error())
	}

	api.Register(app)

	// Embedded SvelteKit build: assets + SPA fallback to index.html.
	app.Use("/", filesystem.New(filesystem.Config{
		Root:         http.FS(staticFS),
		Browse:       false,
		Index:        "index.html",
		NotFoundFile: "index.html",
		Next: func(c *fiber.Ctx) bool {
			return strings.HasPrefix(c.Path(), "/api")
		},
	}))

	return app
}
