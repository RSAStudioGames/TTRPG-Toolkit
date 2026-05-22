package services

import (
	"github.com/gabriel/ttrpg-toolkit/backend/internal/models"
	"github.com/gabriel/ttrpg-toolkit/backend/internal/repository"
)

// MechanicsService handles game system mechanical framework logic.
type MechanicsService struct {
	repo *repository.MechanicsRepository
}

// NewMechanicsService returns a mechanics service.
func NewMechanicsService(repo *repository.MechanicsRepository) *MechanicsService {
	return &MechanicsService{repo: repo}
}

func (svc *MechanicsService) GetMechanics(systemID string) (*models.MechanicsResponse, error) {
	_, err := svc.repo.GetMechanicsBySystemID(systemID)
	return nil, err
}

func (svc *MechanicsService) UpsertMechanics(systemID string, req models.UpsertMechanicsRequest) (*models.MechanicsResponse, error) {
	resJSON, _ := models.MarshalConfig(req.ResolutionConfig)
	progJSON, _ := models.MarshalConfig(req.ProgressionConfig)
	actionJSON, _ := models.MarshalConfig(req.ActionEconomyConfig)
	row := &models.SystemMechanics{
		SystemID:                systemID,
		ResolutionConfigJSON:  resJSON,
		ProgressionConfigJSON:   progJSON,
		ActionEconomyConfigJSON: actionJSON,
	}
	_, err := svc.repo.UpsertMechanics(row)
	return nil, err
}

func (svc *MechanicsService) ListAttributes(systemID string) (*models.ListAttributesResponse, error) {
	rows, err := svc.repo.ListAttributes(systemID)
	if err != nil {
		return nil, err
	}
	return &models.ListAttributesResponse{Items: attributeRowsToResponses(rows)}, nil
}

func (svc *MechanicsService) GetAttribute(systemID, attrID string) (*models.AttributeResponse, error) {
	row, err := svc.repo.GetAttribute(systemID, attrID)
	if err != nil {
		return nil, err
	}
	return attributeRowToResponse(row), nil
}

func (svc *MechanicsService) CreateAttribute(systemID string, req models.CreateAttributeRequest) (*models.AttributeResponse, error) {
	cfg, _ := models.MarshalConfig(req.Config)
	row := &models.SystemAttribute{
		SystemID: systemID, GroupName: req.GroupName, Name: req.Name,
		Type: req.Type, ConfigJSON: cfg, SortOrder: req.SortOrder,
	}
	created, err := svc.repo.CreateAttribute(row)
	if err != nil {
		return nil, err
	}
	return attributeRowToResponse(created), nil
}

func (svc *MechanicsService) UpdateAttribute(systemID, attrID string, req models.UpdateAttributeRequest) (*models.AttributeResponse, error) {
	_, err := svc.repo.GetAttribute(systemID, attrID)
	if err != nil {
		return nil, err
	}
	_ = req
	return nil, repository.ErrNotFound
}

func (svc *MechanicsService) DeleteAttribute(systemID, attrID string) error {
	return svc.repo.DeleteAttribute(systemID, attrID)
}

func (svc *MechanicsService) ListSkills(systemID string) (*models.ListSkillsResponse, error) {
	rows, err := svc.repo.ListSkills(systemID)
	if err != nil {
		return nil, err
	}
	return &models.ListSkillsResponse{Items: skillRowsToResponses(rows)}, nil
}

func (svc *MechanicsService) GetSkill(systemID, skillID string) (*models.SkillResponse, error) {
	row, err := svc.repo.GetSkill(systemID, skillID)
	if err != nil {
		return nil, err
	}
	return skillRowToResponse(row), nil
}

func (svc *MechanicsService) CreateSkill(systemID string, req models.CreateSkillRequest) (*models.SkillResponse, error) {
	cfg, _ := models.MarshalConfig(req.Config)
	row := &models.SystemSkill{
		SystemID: systemID, Name: req.Name, LinkedAttributeID: req.LinkedAttributeID,
		Type: req.Type, Category: req.Category, ConfigJSON: cfg, SortOrder: req.SortOrder,
	}
	created, err := svc.repo.CreateSkill(row)
	if err != nil {
		return nil, err
	}
	return skillRowToResponse(created), nil
}

func (svc *MechanicsService) UpdateSkill(systemID, skillID string, req models.UpdateSkillRequest) (*models.SkillResponse, error) {
	_, err := svc.repo.GetSkill(systemID, skillID)
	if err != nil {
		return nil, err
	}
	_ = req
	return nil, repository.ErrNotFound
}

func (svc *MechanicsService) DeleteSkill(systemID, skillID string) error {
	return svc.repo.DeleteSkill(systemID, skillID)
}

func (svc *MechanicsService) ListResources(systemID string) (*models.ListResourcesResponse, error) {
	rows, err := svc.repo.ListResources(systemID)
	if err != nil {
		return nil, err
	}
	return &models.ListResourcesResponse{Items: resourceRowsToResponses(rows)}, nil
}

func (svc *MechanicsService) GetResource(systemID, resourceID string) (*models.ResourceResponse, error) {
	row, err := svc.repo.GetResource(systemID, resourceID)
	if err != nil {
		return nil, err
	}
	return resourceRowToResponse(row), nil
}

func (svc *MechanicsService) CreateResource(systemID string, req models.CreateResourceRequest) (*models.ResourceResponse, error) {
	cfg, _ := models.MarshalConfig(req.Config)
	row := &models.SystemResource{
		SystemID: systemID, Name: req.Name, Type: req.Type,
		ConfigJSON: cfg, SortOrder: req.SortOrder,
	}
	created, err := svc.repo.CreateResource(row)
	if err != nil {
		return nil, err
	}
	return resourceRowToResponse(created), nil
}

func (svc *MechanicsService) UpdateResource(systemID, resourceID string, req models.UpdateResourceRequest) (*models.ResourceResponse, error) {
	_, err := svc.repo.GetResource(systemID, resourceID)
	if err != nil {
		return nil, err
	}
	_ = req
	return nil, repository.ErrNotFound
}

func (svc *MechanicsService) DeleteResource(systemID, resourceID string) error {
	return svc.repo.DeleteResource(systemID, resourceID)
}

func attributeRowToResponse(row *models.SystemAttribute) *models.AttributeResponse {
	if row == nil {
		return nil
	}
	cfg, _ := models.UnmarshalConfig[models.AttributeConfig](row.ConfigJSON)
	return &models.AttributeResponse{
		ID: row.ID, SystemID: row.SystemID, GroupName: row.GroupName,
		Name: row.Name, Type: row.Type, Config: cfg, SortOrder: row.SortOrder,
	}
}

func attributeRowsToResponses(rows []models.SystemAttribute) []models.AttributeResponse {
	out := make([]models.AttributeResponse, 0, len(rows))
	for i := range rows {
		if r := attributeRowToResponse(&rows[i]); r != nil {
			out = append(out, *r)
		}
	}
	return out
}

func skillRowToResponse(row *models.SystemSkill) *models.SkillResponse {
	if row == nil {
		return nil
	}
	cfg, _ := models.UnmarshalConfig[models.SkillConfig](row.ConfigJSON)
	return &models.SkillResponse{
		ID: row.ID, SystemID: row.SystemID, Name: row.Name,
		LinkedAttributeID: row.LinkedAttributeID, Type: row.Type,
		Category: row.Category, Config: cfg, SortOrder: row.SortOrder,
	}
}

func skillRowsToResponses(rows []models.SystemSkill) []models.SkillResponse {
	out := make([]models.SkillResponse, 0, len(rows))
	for i := range rows {
		if r := skillRowToResponse(&rows[i]); r != nil {
			out = append(out, *r)
		}
	}
	return out
}

func resourceRowToResponse(row *models.SystemResource) *models.ResourceResponse {
	if row == nil {
		return nil
	}
	cfg, _ := models.UnmarshalConfig[models.ResourceConfig](row.ConfigJSON)
	return &models.ResourceResponse{
		ID: row.ID, SystemID: row.SystemID, Name: row.Name,
		Type: row.Type, Config: cfg, SortOrder: row.SortOrder,
	}
}

func resourceRowsToResponses(rows []models.SystemResource) []models.ResourceResponse {
	out := make([]models.ResourceResponse, 0, len(rows))
	for i := range rows {
		if r := resourceRowToResponse(&rows[i]); r != nil {
			out = append(out, *r)
		}
	}
	return out
}
