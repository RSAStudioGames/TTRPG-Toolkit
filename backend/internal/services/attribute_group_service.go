package services

import (
	"fmt"
	"strings"

	"github.com/gabriel/ttrpg-toolkit/backend/internal/models"
)

func validateDerivedNoGroup(cfg *models.AttributeConfig, groupID *string) error {
	if cfg == nil || !cfg.IsDerived {
		return nil
	}
	if groupID != nil && strings.TrimSpace(*groupID) != "" {
		return fmt.Errorf("%w: derived attributes cannot belong to a group", ErrInvalidAttribute)
	}
	return nil
}

func (svc *MechanicsService) resolveAttributeGroupName(
	systemID string,
	groupID *string,
) (*string, *string, error) {
	id := normalizeParentID(groupID)
	if id == nil {
		return nil, nil, nil
	}
	group, err := svc.repo.GetAttributeGroup(systemID, *id)
	if err != nil {
		return nil, nil, err
	}
	name := group.Name
	return id, &name, nil
}

func (svc *MechanicsService) ListAttributeGroups(systemID string) (*models.ListAttributeGroupsResponse, error) {
	if err := svc.ensureSystemExists(systemID); err != nil {
		return nil, err
	}
	rows, err := svc.repo.ListAttributeGroups(systemID)
	if err != nil {
		return nil, err
	}
	items := make([]models.AttributeGroupResponse, 0, len(rows))
	for i := range rows {
		items = append(items, attributeGroupRowToResponse(&rows[i]))
	}
	return &models.ListAttributeGroupsResponse{Items: items}, nil
}

func (svc *MechanicsService) CreateAttributeGroup(
	systemID string,
	req models.CreateAttributeGroupRequest,
) (*models.AttributeGroupResponse, error) {
	if err := svc.ensureSystemExists(systemID); err != nil {
		return nil, err
	}
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, fmt.Errorf("%w: group name is required", ErrInvalidAttribute)
	}
	exists, err := svc.repo.AttributeGroupNameExists(systemID, name, "")
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("%w: group name already in use", ErrInvalidAttribute)
	}
	row := &models.SystemAttributeGroup{
		SystemID: systemID, Name: name, SortOrder: req.SortOrder,
	}
	created, err := svc.repo.CreateAttributeGroup(row)
	if err != nil {
		return nil, err
	}
	resp := attributeGroupRowToResponse(created)
	return &resp, nil
}

func (svc *MechanicsService) UpdateAttributeGroup(
	systemID, groupID string,
	req models.UpdateAttributeGroupRequest,
) (*models.AttributeGroupResponse, error) {
	existing, err := svc.repo.GetAttributeGroup(systemID, groupID)
	if err != nil {
		return nil, err
	}
	row := *existing
	if req.Name != nil {
		row.Name = strings.TrimSpace(*req.Name)
		if row.Name == "" {
			return nil, fmt.Errorf("%w: group name is required", ErrInvalidAttribute)
		}
		exists, err := svc.repo.AttributeGroupNameExists(systemID, row.Name, groupID)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, fmt.Errorf("%w: group name already in use", ErrInvalidAttribute)
		}
	}
	if req.SortOrder != nil {
		row.SortOrder = *req.SortOrder
	}
	updated, err := svc.repo.UpdateAttributeGroup(&row)
	if err != nil {
		return nil, err
	}
	if req.Name != nil {
		if err := svc.repo.SyncAttributeGroupNames(systemID, groupID, updated.Name); err != nil {
			return nil, err
		}
	}
	resp := attributeGroupRowToResponse(updated)
	return &resp, nil
}

func (svc *MechanicsService) DeleteAttributeGroup(systemID, groupID string) error {
	if _, err := svc.repo.GetAttributeGroup(systemID, groupID); err != nil {
		return err
	}
	return svc.repo.DeleteAttributeGroup(systemID, groupID)
}

func attributeGroupRowToResponse(row *models.SystemAttributeGroup) models.AttributeGroupResponse {
	return models.AttributeGroupResponse{
		ID: row.ID, SystemID: row.SystemID, Name: row.Name, SortOrder: row.SortOrder,
	}
}
