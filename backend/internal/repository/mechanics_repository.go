package repository

import (
	"database/sql"
	"encoding/json"
	"errors"

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

const mechanicsColumns = `id, system_id, resolution_config, progression_config, action_economy_config, attributes_config`

func normalizeMechanicsJSON(row *models.SystemMechanics) {
	if len(row.ResolutionConfigJSON) == 0 {
		row.ResolutionConfigJSON = models.EmptyJSONObject()
	}
	if len(row.ProgressionConfigJSON) == 0 {
		row.ProgressionConfigJSON = models.EmptyJSONObject()
	}
	if len(row.ActionEconomyConfigJSON) == 0 {
		row.ActionEconomyConfigJSON = models.EmptyJSONObject()
	}
	if len(row.AttributesConfigJSON) == 0 {
		row.AttributesConfigJSON = models.EmptyJSONObject()
	}
}

// GetMechanicsBySystemID loads the mechanics row for a system.
func (r *MechanicsRepository) GetMechanicsBySystemID(systemID string) (*models.SystemMechanics, error) {
	var row models.SystemMechanics
	q := `SELECT ` + mechanicsColumns + ` FROM system_mechanics WHERE system_id = $1`
	if err := r.db.Get(&row, q, systemID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	normalizeMechanicsJSON(&row)
	return &row, nil
}

// UpsertResolutionConfig inserts or updates resolution_config for a system.
func (r *MechanicsRepository) UpsertResolutionConfig(systemID string, resolutionJSON json.RawMessage) (*models.SystemMechanics, error) {
	if len(resolutionJSON) == 0 {
		resolutionJSON = models.EmptyJSONObject()
	}
	var row models.SystemMechanics
	q := `
INSERT INTO system_mechanics (system_id, resolution_config, progression_config, action_economy_config, attributes_config)
VALUES ($1, $2, '{}', '{}', '{}')
ON CONFLICT (system_id) DO UPDATE
  SET resolution_config = EXCLUDED.resolution_config
RETURNING ` + mechanicsColumns
	if err := r.db.Get(&row, q, systemID, resolutionJSON); err != nil {
		return nil, err
	}
	normalizeMechanicsJSON(&row)
	return &row, nil
}

// UpsertProgressionConfig inserts or updates progression_config for a system.
func (r *MechanicsRepository) UpsertProgressionConfig(systemID string, progressionJSON json.RawMessage) (*models.SystemMechanics, error) {
	if len(progressionJSON) == 0 {
		progressionJSON = models.EmptyJSONObject()
	}
	var row models.SystemMechanics
	q := `
INSERT INTO system_mechanics (system_id, resolution_config, progression_config, action_economy_config, attributes_config)
VALUES ($1, '{}', $2, '{}', '{}')
ON CONFLICT (system_id) DO UPDATE
  SET progression_config = EXCLUDED.progression_config
RETURNING ` + mechanicsColumns
	if err := r.db.Get(&row, q, systemID, progressionJSON); err != nil {
		return nil, err
	}
	normalizeMechanicsJSON(&row)
	return &row, nil
}

// UpsertActionEconomyConfig inserts or updates action_economy_config for a system.
func (r *MechanicsRepository) UpsertActionEconomyConfig(systemID string, actionJSON json.RawMessage) (*models.SystemMechanics, error) {
	if len(actionJSON) == 0 {
		actionJSON = models.EmptyJSONObject()
	}
	var row models.SystemMechanics
	q := `
INSERT INTO system_mechanics (system_id, resolution_config, progression_config, action_economy_config, attributes_config)
VALUES ($1, '{}', '{}', $2, '{}')
ON CONFLICT (system_id) DO UPDATE
  SET action_economy_config = EXCLUDED.action_economy_config
RETURNING ` + mechanicsColumns
	if err := r.db.Get(&row, q, systemID, actionJSON); err != nil {
		return nil, err
	}
	normalizeMechanicsJSON(&row)
	return &row, nil
}

// UpsertAttributesConfig inserts or updates attributes_config for a system.
func (r *MechanicsRepository) UpsertAttributesConfig(systemID string, attrsJSON json.RawMessage) (*models.SystemMechanics, error) {
	if len(attrsJSON) == 0 {
		attrsJSON = models.EmptyJSONObject()
	}
	var row models.SystemMechanics
	q := `
INSERT INTO system_mechanics (system_id, resolution_config, progression_config, action_economy_config, attributes_config)
VALUES ($1, '{}', '{}', '{}', $2)
ON CONFLICT (system_id) DO UPDATE
  SET attributes_config = EXCLUDED.attributes_config
RETURNING ` + mechanicsColumns
	if err := r.db.Get(&row, q, systemID, attrsJSON); err != nil {
		return nil, err
	}
	normalizeMechanicsJSON(&row)
	return &row, nil
}

// UpsertMechanics creates or updates mechanics for a system.
func (r *MechanicsRepository) UpsertMechanics(_ *models.SystemMechanics) (*models.SystemMechanics, error) {
	return nil, ErrNotFound
}

const attributeColumns = `id, system_id, group_name, attribute_group_id, name, type, config, sort_order, parent_attribute_id`

const attributeGroupColumns = `id, system_id, name, sort_order`

func normalizeAttributeJSON(row *models.SystemAttribute) {
	if len(row.ConfigJSON) == 0 {
		row.ConfigJSON = models.EmptyJSONObject()
	}
}

// ListAttributes returns attributes for a system ordered by sort_order.
func (r *MechanicsRepository) ListAttributes(systemID string) ([]models.SystemAttribute, error) {
	var rows []models.SystemAttribute
	q := `SELECT ` + attributeColumns + ` FROM system_attributes WHERE system_id = $1 ORDER BY sort_order, name`
	if err := r.db.Select(&rows, q, systemID); err != nil {
		return nil, err
	}
	for i := range rows {
		normalizeAttributeJSON(&rows[i])
	}
	return rows, nil
}

// GetAttribute loads one attribute by id.
func (r *MechanicsRepository) GetAttribute(systemID, attrID string) (*models.SystemAttribute, error) {
	var row models.SystemAttribute
	q := `SELECT ` + attributeColumns + ` FROM system_attributes WHERE system_id = $1 AND id = $2`
	if err := r.db.Get(&row, q, systemID, attrID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	normalizeAttributeJSON(&row)
	return &row, nil
}

// AttributeNameExists reports whether another attribute in the system has the same name (case-insensitive).
func (r *MechanicsRepository) AttributeNameExists(systemID, name, excludeID string) (bool, error) {
	var exists bool
	q := `SELECT EXISTS(
		SELECT 1 FROM system_attributes
		WHERE system_id = $1 AND lower(name) = lower($2) AND ($3 = '' OR id::text != $3)
	)`
	if err := r.db.Get(&exists, q, systemID, name, excludeID); err != nil {
		return false, err
	}
	return exists, nil
}

// CreateAttribute inserts a new attribute.
func (r *MechanicsRepository) CreateAttribute(row *models.SystemAttribute) (*models.SystemAttribute, error) {
	if len(row.ConfigJSON) == 0 {
		row.ConfigJSON = models.EmptyJSONObject()
	}
	var created models.SystemAttribute
	q := `
INSERT INTO system_attributes (system_id, group_name, attribute_group_id, name, type, config, sort_order, parent_attribute_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING ` + attributeColumns
	if err := r.db.Get(&created, q,
		row.SystemID, row.GroupName, row.AttributeGroupID, row.Name, row.Type, row.ConfigJSON, row.SortOrder, row.ParentAttributeID,
	); err != nil {
		return nil, err
	}
	normalizeAttributeJSON(&created)
	return &created, nil
}

// UpdateAttribute updates an existing attribute.
func (r *MechanicsRepository) UpdateAttribute(row *models.SystemAttribute) (*models.SystemAttribute, error) {
	if len(row.ConfigJSON) == 0 {
		row.ConfigJSON = models.EmptyJSONObject()
	}
	var updated models.SystemAttribute
	q := `
UPDATE system_attributes
SET group_name = $3, attribute_group_id = $4, name = $5, type = $6, config = $7, sort_order = $8, parent_attribute_id = $9
WHERE system_id = $1 AND id = $2
RETURNING ` + attributeColumns
	if err := r.db.Get(&updated, q,
		row.SystemID, row.ID, row.GroupName, row.AttributeGroupID, row.Name, row.Type, row.ConfigJSON, row.SortOrder, row.ParentAttributeID,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	normalizeAttributeJSON(&updated)
	return &updated, nil
}

// DeleteAttribute removes an attribute.
func (r *MechanicsRepository) DeleteAttribute(systemID, attrID string) error {
	res, err := r.db.Exec(`DELETE FROM system_attributes WHERE system_id = $1 AND id = $2`, systemID, attrID)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

const skillColumns = `id, system_id, name, linked_attribute_id, type, category, config, sort_order`

func normalizeSkillJSON(row *models.SystemSkill) {
	if len(row.ConfigJSON) == 0 {
		row.ConfigJSON = models.EmptyJSONObject()
	}
}

// ListSkills returns skills for a system.
func (r *MechanicsRepository) ListSkills(systemID string) ([]models.SystemSkill, error) {
	var rows []models.SystemSkill
	q := `SELECT ` + skillColumns + ` FROM system_skills WHERE system_id = $1 ORDER BY sort_order, name`
	if err := r.db.Select(&rows, q, systemID); err != nil {
		return nil, err
	}
	for i := range rows {
		normalizeSkillJSON(&rows[i])
	}
	return rows, nil
}

// GetSkill loads one skill by id.
func (r *MechanicsRepository) GetSkill(systemID, skillID string) (*models.SystemSkill, error) {
	var row models.SystemSkill
	q := `SELECT ` + skillColumns + ` FROM system_skills WHERE system_id = $1 AND id = $2`
	if err := r.db.Get(&row, q, systemID, skillID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	normalizeSkillJSON(&row)
	return &row, nil
}

// SkillNameExists reports whether another skill in the system has the same name (case-insensitive).
func (r *MechanicsRepository) SkillNameExists(systemID, name, excludeID string) (bool, error) {
	var exists bool
	q := `SELECT EXISTS(
		SELECT 1 FROM system_skills
		WHERE system_id = $1 AND lower(name) = lower($2) AND ($3 = '' OR id::text != $3)
	)`
	if err := r.db.Get(&exists, q, systemID, name, excludeID); err != nil {
		return false, err
	}
	return exists, nil
}

// CreateSkill inserts a new skill.
func (r *MechanicsRepository) CreateSkill(row *models.SystemSkill) (*models.SystemSkill, error) {
	if len(row.ConfigJSON) == 0 {
		row.ConfigJSON = models.EmptyJSONObject()
	}
	var created models.SystemSkill
	q := `
INSERT INTO system_skills (system_id, name, linked_attribute_id, type, category, config, sort_order)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING ` + skillColumns
	if err := r.db.Get(&created, q,
		row.SystemID, row.Name, row.LinkedAttributeID, row.Type, row.Category, row.ConfigJSON, row.SortOrder,
	); err != nil {
		return nil, err
	}
	normalizeSkillJSON(&created)
	return &created, nil
}

// UpdateSkill updates an existing skill.
func (r *MechanicsRepository) UpdateSkill(row *models.SystemSkill) (*models.SystemSkill, error) {
	if len(row.ConfigJSON) == 0 {
		row.ConfigJSON = models.EmptyJSONObject()
	}
	var updated models.SystemSkill
	q := `
UPDATE system_skills
SET name = $3, linked_attribute_id = $4, type = $5, category = $6, config = $7, sort_order = $8
WHERE system_id = $1 AND id = $2
RETURNING ` + skillColumns
	if err := r.db.Get(&updated, q,
		row.SystemID, row.ID, row.Name, row.LinkedAttributeID, row.Type, row.Category, row.ConfigJSON, row.SortOrder,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	normalizeSkillJSON(&updated)
	return &updated, nil
}

// DeleteSkill removes a skill.
func (r *MechanicsRepository) DeleteSkill(systemID, skillID string) error {
	res, err := r.db.Exec(`DELETE FROM system_skills WHERE system_id = $1 AND id = $2`, systemID, skillID)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

const resourceColumns = `id, system_id, name, type, config, sort_order`

func normalizeResourceJSON(row *models.SystemResource) {
	if len(row.ConfigJSON) == 0 {
		row.ConfigJSON = models.EmptyJSONObject()
	}
}

// ListResources returns resources for a system.
func (r *MechanicsRepository) ListResources(systemID string) ([]models.SystemResource, error) {
	var rows []models.SystemResource
	q := `SELECT ` + resourceColumns + ` FROM system_resources WHERE system_id = $1 ORDER BY sort_order, name`
	if err := r.db.Select(&rows, q, systemID); err != nil {
		return nil, err
	}
	for i := range rows {
		normalizeResourceJSON(&rows[i])
	}
	return rows, nil
}

// GetResource loads one resource by id.
func (r *MechanicsRepository) GetResource(systemID, resourceID string) (*models.SystemResource, error) {
	var row models.SystemResource
	q := `SELECT ` + resourceColumns + ` FROM system_resources WHERE system_id = $1 AND id = $2`
	if err := r.db.Get(&row, q, systemID, resourceID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	normalizeResourceJSON(&row)
	return &row, nil
}

// ResourceNameExists reports whether another resource in the system has the same name (case-insensitive).
func (r *MechanicsRepository) ResourceNameExists(systemID, name, excludeID string) (bool, error) {
	var exists bool
	q := `SELECT EXISTS(
		SELECT 1 FROM system_resources
		WHERE system_id = $1 AND lower(name) = lower($2) AND ($3 = '' OR id::text != $3)
	)`
	if err := r.db.Get(&exists, q, systemID, name, excludeID); err != nil {
		return false, err
	}
	return exists, nil
}

// CreateResource inserts a new resource.
func (r *MechanicsRepository) CreateResource(row *models.SystemResource) (*models.SystemResource, error) {
	if len(row.ConfigJSON) == 0 {
		row.ConfigJSON = models.EmptyJSONObject()
	}
	var created models.SystemResource
	q := `
INSERT INTO system_resources (system_id, name, type, config, sort_order)
VALUES ($1, $2, $3, $4, $5)
RETURNING ` + resourceColumns
	if err := r.db.Get(&created, q,
		row.SystemID, row.Name, row.Type, row.ConfigJSON, row.SortOrder,
	); err != nil {
		return nil, err
	}
	normalizeResourceJSON(&created)
	return &created, nil
}

// UpdateResource updates an existing resource.
func (r *MechanicsRepository) UpdateResource(row *models.SystemResource) (*models.SystemResource, error) {
	if len(row.ConfigJSON) == 0 {
		row.ConfigJSON = models.EmptyJSONObject()
	}
	var updated models.SystemResource
	q := `
UPDATE system_resources
SET name = $3, type = $4, config = $5, sort_order = $6
WHERE system_id = $1 AND id = $2
RETURNING ` + resourceColumns
	if err := r.db.Get(&updated, q,
		row.SystemID, row.ID, row.Name, row.Type, row.ConfigJSON, row.SortOrder,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	normalizeResourceJSON(&updated)
	return &updated, nil
}

// DeleteResource removes a resource.
func (r *MechanicsRepository) DeleteResource(systemID, resourceID string) error {
	res, err := r.db.Exec(`DELETE FROM system_resources WHERE system_id = $1 AND id = $2`, systemID, resourceID)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

// ListAttributeGroups returns groups for a system ordered by sort_order.
func (r *MechanicsRepository) ListAttributeGroups(systemID string) ([]models.SystemAttributeGroup, error) {
	var rows []models.SystemAttributeGroup
	q := `SELECT ` + attributeGroupColumns + ` FROM system_attribute_groups WHERE system_id = $1 ORDER BY sort_order, name`
	if err := r.db.Select(&rows, q, systemID); err != nil {
		return nil, err
	}
	return rows, nil
}

// GetAttributeGroup loads one group by id.
func (r *MechanicsRepository) GetAttributeGroup(systemID, groupID string) (*models.SystemAttributeGroup, error) {
	var row models.SystemAttributeGroup
	q := `SELECT ` + attributeGroupColumns + ` FROM system_attribute_groups WHERE system_id = $1 AND id = $2`
	if err := r.db.Get(&row, q, systemID, groupID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &row, nil
}

// AttributeGroupNameExists reports duplicate group names (case-insensitive).
func (r *MechanicsRepository) AttributeGroupNameExists(systemID, name, excludeID string) (bool, error) {
	var exists bool
	q := `SELECT EXISTS(
		SELECT 1 FROM system_attribute_groups
		WHERE system_id = $1 AND lower(name) = lower($2) AND ($3 = '' OR id::text != $3)
	)`
	if err := r.db.Get(&exists, q, systemID, name, excludeID); err != nil {
		return false, err
	}
	return exists, nil
}

// CreateAttributeGroup inserts a new group.
func (r *MechanicsRepository) CreateAttributeGroup(row *models.SystemAttributeGroup) (*models.SystemAttributeGroup, error) {
	var created models.SystemAttributeGroup
	q := `
INSERT INTO system_attribute_groups (system_id, name, sort_order)
VALUES ($1, $2, $3)
RETURNING ` + attributeGroupColumns
	if err := r.db.Get(&created, q, row.SystemID, row.Name, row.SortOrder); err != nil {
		return nil, err
	}
	return &created, nil
}

// UpdateAttributeGroup updates a group row.
func (r *MechanicsRepository) UpdateAttributeGroup(row *models.SystemAttributeGroup) (*models.SystemAttributeGroup, error) {
	var updated models.SystemAttributeGroup
	q := `
UPDATE system_attribute_groups
SET name = $3, sort_order = $4
WHERE system_id = $1 AND id = $2
RETURNING ` + attributeGroupColumns
	if err := r.db.Get(&updated, q, row.SystemID, row.ID, row.Name, row.SortOrder); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &updated, nil
}

// SyncAttributeGroupNames updates group_name on all attributes in a group.
func (r *MechanicsRepository) SyncAttributeGroupNames(systemID, groupID, name string) error {
	_, err := r.db.Exec(
		`UPDATE system_attributes SET group_name = $3 WHERE system_id = $1 AND attribute_group_id = $2`,
		systemID, groupID, name,
	)
	return err
}

// ClearAttributeGroupMembers clears group assignment on attributes (Ungrouped).
func (r *MechanicsRepository) ClearAttributeGroupMembers(systemID, groupID string) error {
	_, err := r.db.Exec(
		`UPDATE system_attributes SET attribute_group_id = NULL, group_name = NULL WHERE system_id = $1 AND attribute_group_id = $2`,
		systemID, groupID,
	)
	return err
}

// DeleteAttributeGroup removes a group after clearing members.
func (r *MechanicsRepository) DeleteAttributeGroup(systemID, groupID string) error {
	if err := r.ClearAttributeGroupMembers(systemID, groupID); err != nil {
		return err
	}
	res, err := r.db.Exec(`DELETE FROM system_attribute_groups WHERE system_id = $1 AND id = $2`, systemID, groupID)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}
