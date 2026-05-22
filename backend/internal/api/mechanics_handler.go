package api

import (
	"errors"
	"fmt"

	"github.com/gabriel/ttrpg-toolkit/backend/internal/models"
	"github.com/gabriel/ttrpg-toolkit/backend/internal/repository"
	"github.com/gabriel/ttrpg-toolkit/backend/internal/services"
	"github.com/gofiber/fiber/v2"
)

// MechanicsHandler serves system mechanics API routes.
type MechanicsHandler struct {
	svc     *services.MechanicsService
	formula *services.FormulaService
}

// NewMechanicsHandler returns a mechanics handler.
func NewMechanicsHandler(svc *services.MechanicsService, formula *services.FormulaService) *MechanicsHandler {
	return &MechanicsHandler{svc: svc, formula: formula}
}

// ValidateFormula handles POST /api/systems/:id/validate-formula.
func (h *MechanicsHandler) ValidateFormula(c *fiber.Ctx) error {
	id := c.Params("id")
	var req models.ValidateFormulaRequest
	if err := c.BodyParser(&req); err != nil {
		return WriteError(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}
	if errs := ValidateStruct(req); len(errs) > 0 {
		return WriteError(c, fiber.StatusBadRequest, "Validation failed", errs)
	}
	if h.formula == nil {
		return WriteError(c, fiber.StatusInternalServerError, "Formula validation unavailable", nil)
	}
	valid, errs := h.formula.ValidateFormula(req.Formula)
	if !valid {
		if len(errs) == 0 {
			return WriteError(c, fiber.StatusInternalServerError, "Something went wrong",
				[]string{fmt.Errorf("validating formula for system %s: validation failed without details", id).Error()})
		}
		return WriteError(c, fiber.StatusBadRequest, "Invalid formula", errs)
	}
	return WriteSuccess(c, models.ValidateFormulaResponse{Valid: true})
}

func (h *MechanicsHandler) mechanicsServiceError(c *fiber.Ctx, err error) error {
	if err == nil {
		return WriteError(c, fiber.StatusInternalServerError, "Something went wrong", nil)
	}
	var formulaErr *services.InvalidFormulaError
	var attrFormulaErr *services.AttributeFormulaError
	switch {
	case errors.Is(err, repository.ErrNotFound):
		return WriteError(c, fiber.StatusNotFound, "Not found", nil)
	case errors.As(err, &formulaErr):
		return WriteError(c, fiber.StatusBadRequest, "Invalid formula", formulaErr.Errors)
	case errors.As(err, &attrFormulaErr):
		return WriteError(c, fiber.StatusBadRequest, attrFormulaErr.Error(), attrFormulaErr.Errors)
	case errors.Is(err, services.ErrInvalidResolution):
		return WriteError(c, fiber.StatusBadRequest, err.Error(), []string{err.Error()})
	case errors.Is(err, services.ErrInvalidAttribute):
		return WriteError(c, fiber.StatusBadRequest, err.Error(), []string{err.Error()})
	case errors.Is(err, services.ErrInvalidSkill):
		return WriteError(c, fiber.StatusBadRequest, err.Error(), []string{err.Error()})
	default:
		if msg := err.Error(); msg != "" {
			return WriteError(c, fiber.StatusBadRequest, msg, []string{msg})
		}
		return WriteError(c, fiber.StatusInternalServerError, "Something went wrong", nil)
	}
}

func (h *MechanicsHandler) placeholder(c *fiber.Ctx, feature string) error {
	return WriteError(c, fiber.StatusNotImplemented, feature+" not implemented yet", nil)
}

func (h *MechanicsHandler) GetMechanics(c *fiber.Ctx) error {
	if h.svc == nil {
		return WriteError(c, fiber.StatusInternalServerError, "Mechanics unavailable", nil)
	}
	resp, err := h.svc.GetMechanics(c.Params("id"))
	if err != nil {
		return h.mechanicsServiceError(c, err)
	}
	return WriteSuccess(c, resp)
}

func (h *MechanicsHandler) SaveResolutionConfig(c *fiber.Ctx) error {
	if h.svc == nil {
		return WriteError(c, fiber.StatusInternalServerError, "Mechanics unavailable", nil)
	}
	var req models.SaveResolutionConfigRequest
	if err := c.BodyParser(&req); err != nil {
		return WriteError(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}
	if errs := ValidateStruct(req); len(errs) > 0 {
		return WriteError(c, fiber.StatusBadRequest, "Validation failed", errs)
	}
	resp, err := h.svc.SaveResolutionConfig(c.Params("id"), req.ToResolutionConfig())
	if err != nil {
		return h.mechanicsServiceError(c, err)
	}
	return WriteSuccess(c, resp)
}

func (h *MechanicsHandler) UpsertMechanics(c *fiber.Ctx) error {
	return h.placeholder(c, "Mechanics")
}

func (h *MechanicsHandler) ListAttributes(c *fiber.Ctx) error {
	if h.svc == nil {
		return WriteError(c, fiber.StatusInternalServerError, "Mechanics unavailable", nil)
	}
	resp, err := h.svc.ListAttributes(c.Params("id"))
	if err != nil {
		return h.mechanicsServiceError(c, err)
	}
	return WriteSuccess(c, resp)
}

func (h *MechanicsHandler) CreateAttribute(c *fiber.Ctx) error {
	if h.svc == nil {
		return WriteError(c, fiber.StatusInternalServerError, "Mechanics unavailable", nil)
	}
	var req models.CreateAttributeRequest
	if err := c.BodyParser(&req); err != nil {
		return WriteError(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}
	if errs := ValidateStruct(req); len(errs) > 0 {
		return WriteError(c, fiber.StatusBadRequest, "Validation failed", errs)
	}
	resp, err := h.svc.CreateAttribute(c.Params("id"), req)
	if err != nil {
		return h.mechanicsServiceError(c, err)
	}
	return WriteSuccess(c, resp)
}

func (h *MechanicsHandler) GetAttribute(c *fiber.Ctx) error {
	if h.svc == nil {
		return WriteError(c, fiber.StatusInternalServerError, "Mechanics unavailable", nil)
	}
	resp, err := h.svc.GetAttribute(c.Params("id"), c.Params("attrId"))
	if err != nil {
		return h.mechanicsServiceError(c, err)
	}
	return WriteSuccess(c, resp)
}

func (h *MechanicsHandler) UpdateAttribute(c *fiber.Ctx) error {
	if h.svc == nil {
		return WriteError(c, fiber.StatusInternalServerError, "Mechanics unavailable", nil)
	}
	var req models.UpdateAttributeRequest
	if err := c.BodyParser(&req); err != nil {
		return WriteError(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}
	resp, err := h.svc.UpdateAttribute(c.Params("id"), c.Params("attrId"), req)
	if err != nil {
		return h.mechanicsServiceError(c, err)
	}
	return WriteSuccess(c, resp)
}

func (h *MechanicsHandler) DeleteAttribute(c *fiber.Ctx) error {
	if h.svc == nil {
		return WriteError(c, fiber.StatusInternalServerError, "Mechanics unavailable", nil)
	}
	if err := h.svc.DeleteAttribute(c.Params("id"), c.Params("attrId")); err != nil {
		return h.mechanicsServiceError(c, err)
	}
	return WriteSuccess(c, map[string]string{"deleted": c.Params("attrId")})
}

func (h *MechanicsHandler) ListSkills(c *fiber.Ctx) error {
	if h.svc == nil {
		return WriteError(c, fiber.StatusInternalServerError, "Mechanics unavailable", nil)
	}
	resp, err := h.svc.ListSkills(c.Params("id"))
	if err != nil {
		return h.mechanicsServiceError(c, err)
	}
	return WriteSuccess(c, resp)
}

func (h *MechanicsHandler) CreateSkill(c *fiber.Ctx) error {
	if h.svc == nil {
		return WriteError(c, fiber.StatusInternalServerError, "Mechanics unavailable", nil)
	}
	var req models.CreateSkillRequest
	if err := c.BodyParser(&req); err != nil {
		return WriteError(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}
	if errs := ValidateStruct(req); len(errs) > 0 {
		return WriteError(c, fiber.StatusBadRequest, "Validation failed", errs)
	}
	resp, err := h.svc.CreateSkill(c.Params("id"), req)
	if err != nil {
		return h.mechanicsServiceError(c, err)
	}
	return WriteSuccess(c, resp)
}

func (h *MechanicsHandler) GetSkill(c *fiber.Ctx) error {
	if h.svc == nil {
		return WriteError(c, fiber.StatusInternalServerError, "Mechanics unavailable", nil)
	}
	resp, err := h.svc.GetSkill(c.Params("id"), c.Params("skillId"))
	if err != nil {
		return h.mechanicsServiceError(c, err)
	}
	return WriteSuccess(c, resp)
}

func (h *MechanicsHandler) UpdateSkill(c *fiber.Ctx) error {
	if h.svc == nil {
		return WriteError(c, fiber.StatusInternalServerError, "Mechanics unavailable", nil)
	}
	var req models.UpdateSkillRequest
	if err := c.BodyParser(&req); err != nil {
		return WriteError(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}
	resp, err := h.svc.UpdateSkill(c.Params("id"), c.Params("skillId"), req)
	if err != nil {
		return h.mechanicsServiceError(c, err)
	}
	return WriteSuccess(c, resp)
}

func (h *MechanicsHandler) DeleteSkill(c *fiber.Ctx) error {
	if h.svc == nil {
		return WriteError(c, fiber.StatusInternalServerError, "Mechanics unavailable", nil)
	}
	if err := h.svc.DeleteSkill(c.Params("id"), c.Params("skillId")); err != nil {
		return h.mechanicsServiceError(c, err)
	}
	return WriteSuccess(c, map[string]string{"deleted": c.Params("skillId")})
}

func (h *MechanicsHandler) ListResources(c *fiber.Ctx) error {
	return h.placeholder(c, "Resources")
}

func (h *MechanicsHandler) CreateResource(c *fiber.Ctx) error {
	return h.placeholder(c, "Resources")
}

func (h *MechanicsHandler) GetResource(c *fiber.Ctx) error {
	return h.placeholder(c, "Resources")
}

func (h *MechanicsHandler) UpdateResource(c *fiber.Ctx) error {
	return h.placeholder(c, "Resources")
}

func (h *MechanicsHandler) DeleteResource(c *fiber.Ctx) error {
	return h.placeholder(c, "Resources")
}
