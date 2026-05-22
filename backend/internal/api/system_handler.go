package api

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gabriel/ttrpg-toolkit/backend/internal/models"
	"github.com/gabriel/ttrpg-toolkit/backend/internal/repository"
	"github.com/gabriel/ttrpg-toolkit/backend/internal/services"
	"github.com/gofiber/fiber/v2"
)

type SystemHandler struct {
	svc *services.SystemService
}

func NewSystemHandler(svc *services.SystemService) *SystemHandler {
	return &SystemHandler{svc: svc}
}

func (h *SystemHandler) serviceError(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		return WriteError(c, fiber.StatusNotFound, "System not found", nil)
	case errors.Is(err, services.ErrLocked):
		return WriteError(c, fiber.StatusConflict, "This system is locked and cannot be edited", nil)
	case errors.Is(err, services.ErrProtected):
		return WriteError(c, fiber.StatusConflict, "Disable Deletion Protection before deleting this system", nil)
	case errors.Is(err, services.ErrInvalidTransition):
		return WriteError(c, fiber.StatusConflict, "This status change is not allowed", nil)
	case errors.Is(err, services.ErrSlugInUse):
		return WriteError(c, fiber.StatusConflict, "URL Identifier already in use", []string{"URL Identifier already in use"})
	case errors.Is(err, services.ErrCampaignsActive):
		return WriteError(c, fiber.StatusConflict, "Cannot unlock: active Campaigns reference this system", nil)
	default:
		if msg := err.Error(); msg != "" {
			return WriteError(c, fiber.StatusBadRequest, msg, []string{msg})
		}
		return WriteError(c, fiber.StatusInternalServerError, "Something went wrong", nil)
	}
}

func (h *SystemHandler) Create(c *fiber.Ctx) error {
	var req models.CreateSystemRequest
	if err := c.BodyParser(&req); err != nil {
		return WriteError(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}
	if errs := ValidateStruct(req); len(errs) > 0 {
		return WriteError(c, fiber.StatusBadRequest, "Validation failed", errs)
	}
	s, err := h.svc.Create(req)
	if err != nil {
		return h.serviceError(c, err)
	}
	return WriteSuccess(c, s)
}

func (h *SystemHandler) List(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "20"))
	f := repository.ListFilter{Page: page, PerPage: perPage, Status: c.Query("status")}
	if v := c.Query("is_active"); v != "" {
		active := v == "true" || v == "1"
		f.IsActive = &active
	}
	resp, err := h.svc.List(f)
	if err != nil {
		return h.serviceError(c, err)
	}
	return WriteSuccess(c, resp)
}

func (h *SystemHandler) Get(c *fiber.Ctx) error {
	s, err := h.svc.Get(c.Params("id"))
	if err != nil {
		return h.serviceError(c, err)
	}
	return WriteSuccess(c, s)
}

func (h *SystemHandler) Update(c *fiber.Ctx) error {
	var req models.UpdateSystemRequest
	if err := c.BodyParser(&req); err != nil {
		return WriteError(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}
	if errs := ValidateStruct(req); len(errs) > 0 {
		return WriteError(c, fiber.StatusBadRequest, "Validation failed", errs)
	}
	s, err := h.svc.Update(c.Params("id"), req)
	if err != nil {
		return h.serviceError(c, err)
	}
	return WriteSuccess(c, s)
}

func (h *SystemHandler) Delete(c *fiber.Ctx) error {
	if err := h.svc.Delete(c.Params("id")); err != nil {
		return h.serviceError(c, err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *SystemHandler) DeletePreview(c *fiber.Ctx) error {
	p, err := h.svc.DeletePreview(c.Params("id"))
	if err != nil {
		return h.serviceError(c, err)
	}
	return WriteSuccess(c, p)
}

func (h *SystemHandler) uploadImage(c *fiber.Ctx, kind string, maxBytes int64) error {
	file, err := c.FormFile("file")
	if err != nil {
		return WriteError(c, fiber.StatusBadRequest, "No file uploaded", nil)
	}
	f, err := file.Open()
	if err != nil {
		return WriteError(c, fiber.StatusBadRequest, "Could not read file", nil)
	}
	defer f.Close()
	data, err := services.ReadUploadFile(f, maxBytes)
	if err != nil {
		return WriteError(c, fiber.StatusBadRequest, "File too large", nil)
	}
	s, err := h.svc.SaveImage(c.Params("id"), kind, data)
	if err != nil {
		return h.serviceError(c, err)
	}
	return WriteSuccess(c, s)
}

func (h *SystemHandler) UploadIcon(c *fiber.Ctx) error {
	return h.uploadImage(c, "icon", 2<<20)
}

func (h *SystemHandler) UploadCover(c *fiber.Ctx) error {
	return h.uploadImage(c, "cover", 5<<20)
}

func (h *SystemHandler) action(c *fiber.Ctx, fn func(string) (*models.GameSystem, error)) error {
	s, err := fn(c.Params("id"))
	if err != nil {
		return h.serviceError(c, err)
	}
	return WriteSuccess(c, s)
}

func (h *SystemHandler) Publish(c *fiber.Ctx) error   { return h.action(c, h.svc.Publish) }
func (h *SystemHandler) Lock(c *fiber.Ctx) error      { return h.action(c, h.svc.Lock) }
func (h *SystemHandler) Unlock(c *fiber.Ctx) error     { return h.action(c, h.svc.Unlock) }
func (h *SystemHandler) Archive(c *fiber.Ctx) error    { return h.action(c, h.svc.Archive) }
func (h *SystemHandler) Restore(c *fiber.Ctx) error   { return h.action(c, h.svc.Restore) }
func (h *SystemHandler) Clone(c *fiber.Ctx) error      { return h.action(c, h.svc.Clone) }
func (h *SystemHandler) Fork(c *fiber.Ctx) error       { return h.action(c, h.svc.Fork) }

func (h *SystemHandler) SaveTemplate(c *fiber.Ctx) error {
	var req models.SaveTemplateRequest
	if err := c.BodyParser(&req); err != nil {
		return WriteError(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}
	if errs := ValidateStruct(req); len(errs) > 0 {
		return WriteError(c, fiber.StatusBadRequest, "Validation failed", errs)
	}
	if err := h.svc.SaveTemplate(c.Params("id"), req.TemplateName, req.TemplateDescription); err != nil {
		return h.serviceError(c, err)
	}
	return WriteSuccess(c, fiber.Map{"saved": true})
}

func (h *SystemHandler) Export(c *fiber.Ctx) error {
	data, err := h.svc.Export(c.Params("id"))
	if err != nil {
		return h.serviceError(c, err)
	}
	c.Set("Content-Type", "application/json")
	c.Set("Content-Disposition", `attachment; filename="system-export.json"`)
	return c.Send(data)
}

func (h *SystemHandler) Import(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return WriteError(c, fiber.StatusBadRequest, "No file uploaded", nil)
	}
	f, err := file.Open()
	if err != nil {
		return WriteError(c, fiber.StatusBadRequest, "Could not read file", nil)
	}
	defer f.Close()
	data, err := services.ReadUploadFile(f, 10<<20)
	if err != nil {
		return WriteError(c, fiber.StatusBadRequest, "File too large", nil)
	}
	s, err := h.svc.ImportJSON(data)
	if err != nil {
		return WriteError(c, fiber.StatusBadRequest, "Failed to import Template", []string{err.Error()})
	}
	return WriteSuccess(c, s)
}

// ServeUploads serves stored system images.
func ServeUploads(uploadDir string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		rel := c.Params("*")
		path := filepath.Join(uploadDir, filepath.Clean(rel))
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return c.SendStatus(fiber.StatusNotFound)
		}
		return c.SendFile(path)
	}
}
