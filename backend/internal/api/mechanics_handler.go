package api

import (
	"fmt"

	"github.com/gabriel/ttrpg-toolkit/backend/internal/models"
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

func (h *MechanicsHandler) placeholder(c *fiber.Ctx, feature string) error {
	return WriteError(c, fiber.StatusNotImplemented, feature+" not implemented yet", nil)
}

func (h *MechanicsHandler) GetMechanics(c *fiber.Ctx) error {
	return h.placeholder(c, "Mechanics")
}

func (h *MechanicsHandler) UpsertMechanics(c *fiber.Ctx) error {
	return h.placeholder(c, "Mechanics")
}

func (h *MechanicsHandler) ListAttributes(c *fiber.Ctx) error {
	return h.placeholder(c, "Attributes")
}

func (h *MechanicsHandler) CreateAttribute(c *fiber.Ctx) error {
	return h.placeholder(c, "Attributes")
}

func (h *MechanicsHandler) GetAttribute(c *fiber.Ctx) error {
	return h.placeholder(c, "Attributes")
}

func (h *MechanicsHandler) UpdateAttribute(c *fiber.Ctx) error {
	return h.placeholder(c, "Attributes")
}

func (h *MechanicsHandler) DeleteAttribute(c *fiber.Ctx) error {
	return h.placeholder(c, "Attributes")
}

func (h *MechanicsHandler) ListSkills(c *fiber.Ctx) error {
	return h.placeholder(c, "Skills")
}

func (h *MechanicsHandler) CreateSkill(c *fiber.Ctx) error {
	return h.placeholder(c, "Skills")
}

func (h *MechanicsHandler) GetSkill(c *fiber.Ctx) error {
	return h.placeholder(c, "Skills")
}

func (h *MechanicsHandler) UpdateSkill(c *fiber.Ctx) error {
	return h.placeholder(c, "Skills")
}

func (h *MechanicsHandler) DeleteSkill(c *fiber.Ctx) error {
	return h.placeholder(c, "Skills")
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
