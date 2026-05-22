package repository

import (
	"github.com/gabriel/ttrpg-toolkit/backend/internal/models"
	"github.com/jmoiron/sqlx"
)

// MechanicsRepository persists system mechanics definitions.
type MechanicsRepository struct {
	db *sqlx.DB
}

// NewMechanicsRepository returns a mechanics repository.
func NewMechanicsRepository(db *sqlx.DB) *MechanicsRepository {
	return &MechanicsRepository{db: db}
}

// GetMechanicsBySystemID loads the mechanics row for a system.
func (r *MechanicsRepository) GetMechanicsBySystemID(systemID string) (*models.SystemMechanics, error) {
	return nil, ErrNotFound
}

// UpsertMechanics creates or updates mechanics for a system.
func (r *MechanicsRepository) UpsertMechanics(_ *models.SystemMechanics) (*models.SystemMechanics, error) {
	return nil, ErrNotFound
}

// ListAttributes returns attributes for a system ordered by sort_order.
func (r *MechanicsRepository) ListAttributes(systemID string) ([]models.SystemAttribute, error) {
	_ = systemID
	return []models.SystemAttribute{}, nil
}

// GetAttribute loads one attribute by id.
func (r *MechanicsRepository) GetAttribute(systemID, attrID string) (*models.SystemAttribute, error) {
	_, _ = systemID, attrID
	return nil, ErrNotFound
}

// CreateAttribute inserts a new attribute.
func (r *MechanicsRepository) CreateAttribute(_ *models.SystemAttribute) (*models.SystemAttribute, error) {
	return nil, ErrNotFound
}

// UpdateAttribute updates an existing attribute.
func (r *MechanicsRepository) UpdateAttribute(_ *models.SystemAttribute) (*models.SystemAttribute, error) {
	return nil, ErrNotFound
}

// DeleteAttribute removes an attribute.
func (r *MechanicsRepository) DeleteAttribute(systemID, attrID string) error {
	_, _ = systemID, attrID
	return ErrNotFound
}

// ListSkills returns skills for a system.
func (r *MechanicsRepository) ListSkills(systemID string) ([]models.SystemSkill, error) {
	_ = systemID
	return []models.SystemSkill{}, nil
}

// GetSkill loads one skill by id.
func (r *MechanicsRepository) GetSkill(systemID, skillID string) (*models.SystemSkill, error) {
	_, _ = systemID, skillID
	return nil, ErrNotFound
}

// CreateSkill inserts a new skill.
func (r *MechanicsRepository) CreateSkill(_ *models.SystemSkill) (*models.SystemSkill, error) {
	return nil, ErrNotFound
}

// UpdateSkill updates an existing skill.
func (r *MechanicsRepository) UpdateSkill(_ *models.SystemSkill) (*models.SystemSkill, error) {
	return nil, ErrNotFound
}

// DeleteSkill removes a skill.
func (r *MechanicsRepository) DeleteSkill(systemID, skillID string) error {
	_, _ = systemID, skillID
	return ErrNotFound
}

// ListResources returns resources for a system.
func (r *MechanicsRepository) ListResources(systemID string) ([]models.SystemResource, error) {
	_ = systemID
	return []models.SystemResource{}, nil
}

// GetResource loads one resource by id.
func (r *MechanicsRepository) GetResource(systemID, resourceID string) (*models.SystemResource, error) {
	_, _ = systemID, resourceID
	return nil, ErrNotFound
}

// CreateResource inserts a new resource.
func (r *MechanicsRepository) CreateResource(_ *models.SystemResource) (*models.SystemResource, error) {
	return nil, ErrNotFound
}

// UpdateResource updates an existing resource.
func (r *MechanicsRepository) UpdateResource(_ *models.SystemResource) (*models.SystemResource, error) {
	return nil, ErrNotFound
}

// DeleteResource removes a resource.
func (r *MechanicsRepository) DeleteResource(systemID, resourceID string) error {
	_, _ = systemID, resourceID
	return ErrNotFound
}
