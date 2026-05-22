package api

import (
	"github.com/gabriel/ttrpg-toolkit/backend/internal/services"
	"github.com/gofiber/fiber/v2"
)

// Deps holds shared API dependencies.
type Deps struct {
	Systems *services.SystemService
	UploadDir string
}

// Register mounts API routes on the Fiber app.
func Register(app *fiber.App, deps Deps) {
	app.Get("/api/config", ConfigHandler)

	if deps.Systems != nil {
		h := NewSystemHandler(deps.Systems)
		systems := app.Group("/api/systems")
		systems.Post("/", h.Create)
		systems.Get("/", h.List)
		systems.Post("/import", h.Import)
		systems.Get("/:id", h.Get)
		systems.Put("/:id", h.Update)
		systems.Delete("/:id", h.Delete)
		systems.Get("/:id/delete-preview", h.DeletePreview)
		systems.Post("/:id/icon", h.UploadIcon)
		systems.Post("/:id/cover", h.UploadCover)
		systems.Post("/:id/publish", h.Publish)
		systems.Post("/:id/lock", h.Lock)
		systems.Post("/:id/unlock", h.Unlock)
		systems.Post("/:id/archive", h.Archive)
		systems.Post("/:id/restore", h.Restore)
		systems.Post("/:id/clone", h.Clone)
		systems.Post("/:id/fork", h.Fork)
		systems.Post("/:id/save-template", h.SaveTemplate)
		systems.Get("/:id/export", h.Export)
	}

	if deps.UploadDir != "" {
		app.Get("/api/uploads/*", ServeUploads(deps.UploadDir))
	}
}
