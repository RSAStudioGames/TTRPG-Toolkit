package services

import (
	"fmt"
	"strings"

	"github.com/gabriel/ttrpg-toolkit/backend/internal/models"
	"github.com/gabriel/ttrpg-toolkit/backend/internal/repository"
)

// MechanicsService handles game system mechanical framework logic.
type MechanicsService struct {
	repo        *repository.MechanicsRepository
	systems     *repository.SystemRepository
	formula     *FormulaService
}

// NewMechanicsService returns a mechanics service.
func NewMechanicsService(
	repo *repository.MechanicsRepository,
	systems *repository.SystemRepository,
	formula *FormulaService,
) *MechanicsService {
	return &MechanicsService{repo: repo, systems: systems, formula: formula}
}

func (svc *MechanicsService) GetMechanics(systemID string) (*models.MechanicsResponse, error) {
	row, err := svc.repo.GetMechanicsBySystemID(systemID)
	if err != nil {
		return nil, err
	}
	return mechanicsRowToResponse(row)
}

func (svc *MechanicsService) SaveResolutionConfig(systemID string, cfg models.ResolutionConfig) (*models.MechanicsResponse, error) {
	if svc.systems != nil {
		if _, err := svc.systems.GetByID(systemID); err != nil {
			return nil, err
		}
	}
	if err := validateResolutionConfig(&cfg); err != nil {
		return nil, err
	}
	sanitizeResolutionConfig(&cfg)
	if svc.formula != nil {
		valid, errs := svc.formula.ValidateFormula(cfg.RollExpression)
		if !valid {
			return nil, &InvalidFormulaError{Errors: errs}
		}
	}
	resJSON, err := models.MarshalConfig(cfg)
	if err != nil {
		return nil, err
	}
	row, err := svc.repo.UpsertResolutionConfig(systemID, resJSON)
	if err != nil {
		return nil, err
	}
	return mechanicsRowToResponse(row)
}

func (svc *MechanicsService) SaveProgressionConfig(systemID string, cfg models.ProgressionConfig) (*models.MechanicsResponse, error) {
	if err := svc.ensureSystemExists(systemID); err != nil {
		return nil, err
	}
	sanitizeProgressionConfig(&cfg)
	if err := validateProgressionConfig(&cfg); err != nil {
		return nil, err
	}
	progJSON, err := models.MarshalConfig(cfg)
	if err != nil {
		return nil, err
	}
	row, err := svc.repo.UpsertProgressionConfig(systemID, progJSON)
	if err != nil {
		return nil, err
	}
	return mechanicsRowToResponse(row)
}

func (svc *MechanicsService) SaveActionEconomyConfig(systemID string, cfg models.ActionEconomyConfig) (*models.MechanicsResponse, error) {
	if err := svc.ensureSystemExists(systemID); err != nil {
		return nil, err
	}
	sanitizeActionEconomyConfig(&cfg)
	if err := validateActionEconomyConfig(&cfg); err != nil {
		return nil, err
	}
	if svc.formula != nil {
		if initiativeNeedsFormula(cfg.InitiativeSystem) {
			valid, errs := svc.formula.ValidateFormula(cfg.InitiativeExpression)
			if !valid {
				return nil, &InvalidFormulaError{Errors: errs}
			}
		}
		if cfg.InitiativeSystem == models.InitiativeSystemStatic && expressionLike(cfg.StaticInitiativeValue) {
			valid, errs := svc.formula.ValidateFormula(cfg.StaticInitiativeValue)
			if !valid {
				return nil, &InvalidFormulaError{Errors: errs}
			}
		}
	}
	actionJSON, err := models.MarshalConfig(cfg)
	if err != nil {
		return nil, err
	}
	row, err := svc.repo.UpsertActionEconomyConfig(systemID, actionJSON)
	if err != nil {
		return nil, err
	}
	return mechanicsRowToResponse(row)
}

func expressionLike(s string) bool {
	s = strings.TrimSpace(s)
	return strings.ContainsAny(s, "{}+-*/")
}

func mechanicsRowToResponse(row *models.SystemMechanics) (*models.MechanicsResponse, error) {
	if row == nil {
		return nil, nil
	}
	res, err := models.UnmarshalConfig[models.ResolutionConfig](row.ResolutionConfigJSON)
	if err != nil {
		return nil, err
	}
	prog, err := models.UnmarshalConfig[models.ProgressionConfig](row.ProgressionConfigJSON)
	if err != nil {
		return nil, err
	}
	action, err := models.UnmarshalConfig[models.ActionEconomyConfig](row.ActionEconomyConfigJSON)
	if err != nil {
		return nil, err
	}
	return &models.MechanicsResponse{
		ID:                  row.ID,
		SystemID:            row.SystemID,
		ResolutionConfig:    res,
		ProgressionConfig:   prog,
		ActionEconomyConfig: action,
	}, nil
}

func validateResolutionConfig(cfg *models.ResolutionConfig) error {
	if cfg == nil {
		return fmt.Errorf("%w: config is required", ErrInvalidResolution)
	}
	if !models.IsAllowedResolutionType(cfg.ResolutionType) {
		return fmt.Errorf("%w: invalid resolution_type", ErrInvalidResolution)
	}
	if cfg.ResolutionType == models.ResolutionTypeCustom && strings.TrimSpace(cfg.CustomParadigmName) == "" {
		return fmt.Errorf("%w: custom_paradigm_name is required when resolution_type is custom", ErrInvalidResolution)
	}
	sd := cfg.SuccessDetermination
	if !models.IsAllowedSuccessMethod(sd.Method) {
		return fmt.Errorf("%w: invalid success_determination.method", ErrInvalidResolution)
	}
	if sd.Method == models.SuccessMethodSuccessThresholdLadder {
		if len(sd.ThresholdLadder) == 0 {
			return fmt.Errorf("%w: threshold_ladder is required when method is success_threshold_ladder", ErrInvalidResolution)
		}
	}
	for i, tier := range sd.ThresholdLadder {
		if tier.Operator != "" && !models.IsAllowedLadderOperator(tier.Operator) {
			return fmt.Errorf("%w: invalid threshold_ladder[%d].operator", ErrInvalidResolution, i)
		}
	}
	cm := cfg.CriticalMechanics
	if cm.EnableCritSuccess && strings.TrimSpace(cm.CritSuccessTrigger) == "" {
		return fmt.Errorf("%w: crit_success_trigger is required when enable_crit_success is true", ErrInvalidResolution)
	}
	if cm.EnableCritFailure && strings.TrimSpace(cm.CritFailureTrigger) == "" {
		return fmt.Errorf("%w: crit_failure_trigger is required when enable_crit_failure is true", ErrInvalidResolution)
	}
	for i, entry := range cfg.AdvantageDisadvantage {
		if strings.TrimSpace(entry.Name) == "" {
			return fmt.Errorf("%w: advantage_disadvantage[%d].name is required", ErrInvalidResolution, i)
		}
		if !models.IsAllowedMechanicType(entry.MechanicType) {
			return fmt.Errorf("%w: invalid advantage_disadvantage[%d].mechanic_type", ErrInvalidResolution, i)
		}
	}
	return nil
}

func sanitizeResolutionConfig(cfg *models.ResolutionConfig) {
	if cfg == nil {
		return
	}
	if cfg.ResolutionType != models.ResolutionTypeCustom {
		cfg.CustomParadigmName = ""
	}
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
	if err := svc.ensureSystemExists(systemID); err != nil {
		return nil, err
	}
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
	if err := svc.ensureSystemExists(systemID); err != nil {
		return nil, err
	}
	if err := svc.ensureUniqueAttributeName(systemID, req.Name, ""); err != nil {
		return nil, err
	}
	cfg := req.Config
	sanitizeAttributeConfig(&cfg)
	if err := validateAttributeConfig(req.Type, &cfg); err != nil {
		return nil, err
	}
	if err := svc.validateAttributeFormulas(&cfg); err != nil {
		return nil, err
	}
	if err := svc.validateParentAttribute(systemID, "", req.ParentAttributeID); err != nil {
		return nil, err
	}
	cfgJSON, err := models.MarshalConfig(cfg)
	if err != nil {
		return nil, err
	}
	row := &models.SystemAttribute{
		SystemID:          systemID,
		GroupName:         req.GroupName,
		ParentAttributeID: normalizeParentID(req.ParentAttributeID),
		Name:              req.Name,
		Type:              req.Type,
		ConfigJSON:        cfgJSON,
		SortOrder:         req.SortOrder,
	}
	created, err := svc.repo.CreateAttribute(row)
	if err != nil {
		return nil, err
	}
	return attributeRowToResponse(created), nil
}

func (svc *MechanicsService) UpdateAttribute(systemID, attrID string, req models.UpdateAttributeRequest) (*models.AttributeResponse, error) {
	existing, err := svc.repo.GetAttribute(systemID, attrID)
	if err != nil {
		return nil, err
	}
	row := *existing
	if req.GroupName != nil {
		row.GroupName = req.GroupName
	}
	if req.ParentAttributeID != nil {
		row.ParentAttributeID = normalizeParentID(req.ParentAttributeID)
	}
	if req.Name != nil {
		row.Name = *req.Name
	}
	if req.Type != nil {
		row.Type = *req.Type
	}
	if req.SortOrder != nil {
		row.SortOrder = *req.SortOrder
	}
	cfg, err := models.UnmarshalConfig[models.AttributeConfig](row.ConfigJSON)
	if err != nil {
		return nil, err
	}
	if req.Config != nil {
		cfg = *req.Config
	}
	sanitizeAttributeConfig(&cfg)
	if err := validateAttributeConfig(row.Type, &cfg); err != nil {
		return nil, err
	}
	if err := svc.validateAttributeFormulas(&cfg); err != nil {
		return nil, err
	}
	if err := svc.validateParentAttribute(systemID, attrID, row.ParentAttributeID); err != nil {
		return nil, err
	}
	if err := svc.ensureUniqueAttributeName(systemID, row.Name, attrID); err != nil {
		return nil, err
	}
	row.ConfigJSON, err = models.MarshalConfig(cfg)
	if err != nil {
		return nil, err
	}
	updated, err := svc.repo.UpdateAttribute(&row)
	if err != nil {
		return nil, err
	}
	return attributeRowToResponse(updated), nil
}

func (svc *MechanicsService) DeleteAttribute(systemID, attrID string) error {
	if _, err := svc.repo.GetAttribute(systemID, attrID); err != nil {
		return err
	}
	return svc.repo.DeleteAttribute(systemID, attrID)
}

func normalizeParentID(id *string) *string {
	if id == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*id)
	if trimmed == "" {
		return nil
	}
	return &trimmed
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
	if err := svc.ensureSystemExists(systemID); err != nil {
		return nil, err
	}
	name := strings.TrimSpace(req.Name)
	if err := svc.ensureUniqueSkillName(systemID, name, ""); err != nil {
		return nil, err
	}
	cfg := req.Config
	sanitizeSkillConfig(req.Type, &cfg)
	if err := validateSkillConfig(req.Type, &cfg); err != nil {
		return nil, err
	}
	if err := svc.validateLinkedAttribute(systemID, req.LinkedAttributeID); err != nil {
		return nil, err
	}
	cfgJSON, err := models.MarshalConfig(cfg)
	if err != nil {
		return nil, err
	}
	row := &models.SystemSkill{
		SystemID:          systemID,
		Name:              name,
		LinkedAttributeID: normalizeParentID(req.LinkedAttributeID),
		Type:              req.Type,
		Category:          req.Category,
		ConfigJSON:        cfgJSON,
		SortOrder:         req.SortOrder,
	}
	created, err := svc.repo.CreateSkill(row)
	if err != nil {
		return nil, err
	}
	return skillRowToResponse(created), nil
}

func (svc *MechanicsService) UpdateSkill(systemID, skillID string, req models.UpdateSkillRequest) (*models.SkillResponse, error) {
	existing, err := svc.repo.GetSkill(systemID, skillID)
	if err != nil {
		return nil, err
	}
	row := *existing
	if req.Name != nil {
		row.Name = strings.TrimSpace(*req.Name)
	}
	if req.LinkedAttributeID != nil {
		row.LinkedAttributeID = normalizeParentID(req.LinkedAttributeID)
	}
	if req.Type != nil {
		row.Type = *req.Type
	}
	if req.Category != nil {
		row.Category = req.Category
	}
	if req.SortOrder != nil {
		row.SortOrder = *req.SortOrder
	}
	cfg, err := models.UnmarshalConfig[models.SkillConfig](row.ConfigJSON)
	if err != nil {
		return nil, err
	}
	if req.Config != nil {
		cfg = *req.Config
	}
	sanitizeSkillConfig(row.Type, &cfg)
	if err := validateSkillConfig(row.Type, &cfg); err != nil {
		return nil, err
	}
	if err := svc.validateLinkedAttribute(systemID, row.LinkedAttributeID); err != nil {
		return nil, err
	}
	if err := svc.ensureUniqueSkillName(systemID, row.Name, skillID); err != nil {
		return nil, err
	}
	row.ConfigJSON, err = models.MarshalConfig(cfg)
	if err != nil {
		return nil, err
	}
	updated, err := svc.repo.UpdateSkill(&row)
	if err != nil {
		return nil, err
	}
	return skillRowToResponse(updated), nil
}

func (svc *MechanicsService) DeleteSkill(systemID, skillID string) error {
	if _, err := svc.repo.GetSkill(systemID, skillID); err != nil {
		return err
	}
	return svc.repo.DeleteSkill(systemID, skillID)
}

func (svc *MechanicsService) ListResources(systemID string) (*models.ListResourcesResponse, error) {
	if err := svc.ensureSystemExists(systemID); err != nil {
		return nil, err
	}
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
	if err := svc.ensureSystemExists(systemID); err != nil {
		return nil, err
	}
	name := strings.TrimSpace(req.Name)
	if err := svc.ensureUniqueResourceName(systemID, name, ""); err != nil {
		return nil, err
	}
	cfg := req.Config
	sanitizeResourceConfig(&cfg)
	if err := validateResourceConfig(&cfg); err != nil {
		return nil, err
	}
	cfgJSON, err := models.MarshalConfig(cfg)
	if err != nil {
		return nil, err
	}
	row := &models.SystemResource{
		SystemID: systemID, Name: name, Type: req.Type,
		ConfigJSON: cfgJSON, SortOrder: req.SortOrder,
	}
	created, err := svc.repo.CreateResource(row)
	if err != nil {
		return nil, err
	}
	return resourceRowToResponse(created), nil
}

func (svc *MechanicsService) UpdateResource(systemID, resourceID string, req models.UpdateResourceRequest) (*models.ResourceResponse, error) {
	existing, err := svc.repo.GetResource(systemID, resourceID)
	if err != nil {
		return nil, err
	}
	row := *existing
	if req.Name != nil {
		row.Name = strings.TrimSpace(*req.Name)
	}
	if req.Type != nil {
		row.Type = *req.Type
	}
	if req.SortOrder != nil {
		row.SortOrder = *req.SortOrder
	}
	cfg, err := models.UnmarshalConfig[models.ResourceConfig](row.ConfigJSON)
	if err != nil {
		return nil, err
	}
	if req.Config != nil {
		cfg = *req.Config
	}
	sanitizeResourceConfig(&cfg)
	if err := validateResourceConfig(&cfg); err != nil {
		return nil, err
	}
	if err := svc.ensureUniqueResourceName(systemID, row.Name, resourceID); err != nil {
		return nil, err
	}
	row.ConfigJSON, err = models.MarshalConfig(cfg)
	if err != nil {
		return nil, err
	}
	updated, err := svc.repo.UpdateResource(&row)
	if err != nil {
		return nil, err
	}
	return resourceRowToResponse(updated), nil
}

func (svc *MechanicsService) DeleteResource(systemID, resourceID string) error {
	if _, err := svc.repo.GetResource(systemID, resourceID); err != nil {
		return err
	}
	return svc.repo.DeleteResource(systemID, resourceID)
}

func attributeRowToResponse(row *models.SystemAttribute) *models.AttributeResponse {
	if row == nil {
		return nil
	}
	cfg, _ := models.UnmarshalConfig[models.AttributeConfig](row.ConfigJSON)
	return &models.AttributeResponse{
		ID: row.ID, SystemID: row.SystemID, GroupName: row.GroupName,
		ParentAttributeID: row.ParentAttributeID,
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
