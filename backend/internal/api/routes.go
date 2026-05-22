package api

import (
	"github.com/gabriel/ttrpg-toolkit/backend/internal/services"
	"github.com/gofiber/fiber/v2"
)

// Deps holds shared API dependencies.
type Deps struct {
	Systems   *services.SystemService
	Mechanics *services.MechanicsService
	Formula   *services.FormulaService
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

		if deps.Mechanics != nil && deps.Formula != nil {
			mh := NewMechanicsHandler(deps.Mechanics, deps.Formula)
			systems.Post("/:id/validate-formula", mh.ValidateFormula)
			systems.Get("/:id/mechanics", mh.GetMechanics)
			systems.Put("/:id/mechanics", mh.UpsertMechanics)
			systems.Get("/:id/attributes", mh.ListAttributes)
			systems.Post("/:id/attributes", mh.CreateAttribute)
			systems.Get("/:id/attributes/:attrId", mh.GetAttribute)
			systems.Put("/:id/attributes/:attrId", mh.UpdateAttribute)
			systems.Delete("/:id/attributes/:attrId", mh.DeleteAttribute)
			systems.Get("/:id/skills", mh.ListSkills)
			systems.Post("/:id/skills", mh.CreateSkill)
			systems.Get("/:id/skills/:skillId", mh.GetSkill)
			systems.Put("/:id/skills/:skillId", mh.UpdateSkill)
			systems.Delete("/:id/skills/:skillId", mh.DeleteSkill)
			systems.Get("/:id/resources", mh.ListResources)
			systems.Post("/:id/resources", mh.CreateResource)
			systems.Get("/:id/resources/:resourceId", mh.GetResource)
			systems.Put("/:id/resources/:resourceId", mh.UpdateResource)
			systems.Delete("/:id/resources/:resourceId", mh.DeleteResource)
		}

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
