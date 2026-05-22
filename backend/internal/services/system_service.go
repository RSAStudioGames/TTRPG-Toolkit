package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/disintegration/imaging"
	"github.com/gabriel/ttrpg-toolkit/backend/internal/models"
	"github.com/gabriel/ttrpg-toolkit/backend/internal/repository"
	"github.com/google/uuid"
	"github.com/gosimple/slug"
)

var (
	ErrLocked           = errors.New("system is locked")
	ErrProtected        = errors.New("system is deletion protected")
	ErrInvalidTransition = errors.New("invalid status transition")
	ErrSlugInUse        = errors.New("URL Identifier already in use")
	ErrCampaignsActive  = errors.New("cannot unlock: active Campaigns reference this system")
)

var slugPattern = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

type SystemService struct {
	repo      *repository.SystemRepository
	uploadDir string
}

func NewSystemService(repo *repository.SystemRepository, uploadDir string) *SystemService {
	return &SystemService{repo: repo, uploadDir: uploadDir}
}

func newID() string {
	return uuid.Must(uuid.NewV7()).String()
}

func normalizeSlug(s string) string {
	s = strings.TrimSpace(strings.ToLower(s))
	return slug.Make(s)
}

func (svc *SystemService) ensureSlug(name, requested string, excludeID string) (string, error) {
	s := normalizeSlug(requested)
	if s == "" {
		s = normalizeSlug(name)
	}
	if s == "" || !slugPattern.MatchString(s) {
		return "", fmt.Errorf("URL Identifier must be lowercase alphanumeric with hyphens")
	}
	exists, err := svc.repo.SlugExists(s, excludeID)
	if err != nil {
		return "", err
	}
	if exists {
		return "", ErrSlugInUse
	}
	return s, nil
}

func parseVersion(v string) (string, error) {
	if v == "" {
		v = "0.1.0"
	}
	if _, err := semver.NewVersion(v); err != nil {
		return "", fmt.Errorf("Version must be valid semantic version (e.g. 0.1.0)")
	}
	return v, nil
}

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func (svc *SystemService) Create(req models.CreateSystemRequest) (*models.GameSystem, error) {
	slugVal, err := svc.ensureSlug(req.Name, req.Slug, "")
	if err != nil {
		return nil, err
	}
	ver, err := parseVersion(req.Version)
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	active := true
	if req.IsActive != nil {
		active = *req.IsActive
	}
	isCore := true
	if req.IsCore != nil {
		isCore = *req.IsCore
	}
	var parent *string
	if req.ParentSystemID != "" {
		parent = &req.ParentSystemID
	}

	s := &models.GameSystem{
		ID: newID(), Name: req.Name, Slug: slugVal,
		Edition: strPtr(req.Edition), Publisher: strPtr(req.Publisher),
		Description: strPtr(req.Description), LicenseType: strPtr(req.LicenseType),
		Version: ver, Playstyle: strPtr(req.Playstyle), Complexity: req.Complexity,
		MeasurementUnit: strPtr(req.MeasurementUnit), CurrencySymbol: strPtr(req.CurrencySymbol),
		Status: models.StatusDraft, IsActive: active, SystemFamily: strPtr(req.SystemFamily),
		PlayerCountMin: req.PlayerCountMin, PlayerCountMax: req.PlayerCountMax,
		OfficialLinks: req.OfficialLinks, Tags: req.Tags, CoreRulebooks: req.CoreRulebooks,
		IsCore: isCore, ParentSystemID: parent, IsProtected: false,
		CreatedAt: now, UpdatedAt: now,
	}
	if s.OfficialLinks == nil {
		s.OfficialLinks = []string{}
	}
	if s.Tags == nil {
		s.Tags = []string{}
	}
	if s.CoreRulebooks == nil {
		s.CoreRulebooks = []string{}
	}
	if err := svc.repo.Create(s); err != nil {
		if errors.Is(err, repository.ErrSlugConflict) {
			return nil, ErrSlugInUse
		}
		return nil, err
	}
	_ = svc.repo.SyncTags(s.ID, s.Tags)
	_ = svc.repo.SyncRulebooks(s.ID, s.CoreRulebooks)
	_ = svc.repo.SyncLinks(s.ID, s.OfficialLinks)
	return svc.repo.GetByID(s.ID)
}

func (svc *SystemService) List(f repository.ListFilter) (*models.SystemListResponse, error) {
	items, total, err := svc.repo.List(f)
	if err != nil {
		return nil, err
	}
	perPage := f.PerPage
	if perPage < 1 {
		perPage = 20
	}
	pages := (total + perPage - 1) / perPage
	return &models.SystemListResponse{
		Items: items, Page: f.Page, PerPage: perPage, Total: total, TotalPages: pages,
	}, nil
}

func (svc *SystemService) Get(id string) (*models.GameSystem, error) {
	return svc.repo.GetByID(id)
}

func (svc *SystemService) Update(id string, req models.UpdateSystemRequest) (*models.GameSystem, error) {
	s, err := svc.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if s.Status == models.StatusLocked {
		return nil, ErrLocked
	}

	if req.Name != nil {
		s.Name = *req.Name
	}
	if req.Slug != nil {
		slugVal, err := svc.ensureSlug(s.Name, *req.Slug, id)
		if err != nil {
			return nil, err
		}
		s.Slug = slugVal
	}
	if req.Edition != nil {
		s.Edition = req.Edition
	}
	if req.Publisher != nil {
		s.Publisher = req.Publisher
	}
	if req.Description != nil {
		s.Description = req.Description
	}
	if req.LicenseType != nil {
		s.LicenseType = req.LicenseType
	}
	if req.Version != nil {
		ver, err := parseVersion(*req.Version)
		if err != nil {
			return nil, err
		}
		s.Version = ver
	}
	if req.Playstyle != nil {
		s.Playstyle = req.Playstyle
	}
	if req.Complexity != nil {
		s.Complexity = req.Complexity
	}
	if req.MeasurementUnit != nil {
		s.MeasurementUnit = req.MeasurementUnit
	}
	if req.CurrencySymbol != nil {
		s.CurrencySymbol = req.CurrencySymbol
	}
	if req.IsActive != nil {
		s.IsActive = *req.IsActive
	}
	if req.SystemFamily != nil {
		s.SystemFamily = req.SystemFamily
	}
	if req.PlayerCountMin != nil {
		s.PlayerCountMin = req.PlayerCountMin
	}
	if req.PlayerCountMax != nil {
		s.PlayerCountMax = req.PlayerCountMax
	}
	if req.OfficialLinks != nil {
		s.OfficialLinks = req.OfficialLinks
	}
	if req.Tags != nil {
		s.Tags = req.Tags
	}
	if req.CoreRulebooks != nil {
		s.CoreRulebooks = req.CoreRulebooks
	}
	if req.IsCore != nil {
		s.IsCore = *req.IsCore
	}
	if req.ParentSystemID != nil {
		if *req.ParentSystemID == "" {
			s.ParentSystemID = nil
		} else {
			s.ParentSystemID = req.ParentSystemID
		}
	}
	if req.IsProtected != nil {
		s.IsProtected = *req.IsProtected
	}

	s.UpdatedAt = time.Now().UTC()
	if err := svc.repo.Update(s); err != nil {
		if errors.Is(err, repository.ErrSlugConflict) {
			return nil, ErrSlugInUse
		}
		return nil, err
	}
	if req.Tags != nil {
		_ = svc.repo.SyncTags(id, s.Tags)
	}
	if req.CoreRulebooks != nil {
		_ = svc.repo.SyncRulebooks(id, s.CoreRulebooks)
	}
	if req.OfficialLinks != nil {
		_ = svc.repo.SyncLinks(id, s.OfficialLinks)
	}
	return svc.repo.GetByID(id)
}

func (svc *SystemService) Delete(id string) error {
	s, err := svc.repo.GetByID(id)
	if err != nil {
		return err
	}
	if s.IsProtected {
		return ErrProtected
	}
	return svc.repo.Delete(id)
}

func (svc *SystemService) DeletePreview(id string) (*models.DeletePreview, error) {
	return svc.repo.DeletePreview(id)
}

func (svc *SystemService) transition(id, target string) (*models.GameSystem, error) {
	s, err := svc.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	allowed := map[string]map[string]bool{
		models.StatusDraft:     {models.StatusPublished: true, models.StatusArchived: true},
		models.StatusPublished: {models.StatusLocked: true, models.StatusArchived: true},
		models.StatusLocked:    {models.StatusPublished: true},
		models.StatusArchived:  {models.StatusDraft: true},
	}
	if !allowed[s.Status][target] {
		return nil, ErrInvalidTransition
	}
	if target == models.StatusPublished && s.Status == models.StatusLocked {
		n, err := svc.repo.CountActiveCampaignsBySystemID(id)
		if err != nil {
			return nil, err
		}
		if n > 0 {
			return nil, ErrCampaignsActive
		}
	}
	s.Status = target
	s.UpdatedAt = time.Now().UTC()
	if err := svc.repo.Update(s); err != nil {
		return nil, err
	}
	return svc.repo.GetByID(id)
}

func (svc *SystemService) Publish(id string) (*models.GameSystem, error) {
	return svc.transition(id, models.StatusPublished)
}
func (svc *SystemService) Lock(id string) (*models.GameSystem, error) {
	return svc.transition(id, models.StatusLocked)
}
func (svc *SystemService) Unlock(id string) (*models.GameSystem, error) {
	return svc.transition(id, models.StatusPublished)
}
func (svc *SystemService) Archive(id string) (*models.GameSystem, error) {
	return svc.transition(id, models.StatusArchived)
}
func (svc *SystemService) Restore(id string) (*models.GameSystem, error) {
	return svc.transition(id, models.StatusDraft)
}

func (svc *SystemService) Clone(id string) (*models.GameSystem, error) {
	src, err := svc.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	copyName := src.Name + " (Copy)"
	copySlug, err := svc.ensureSlug(copyName, "", "")
	if err != nil {
		base := normalizeSlug(copyName)
		for i := 2; i < 100; i++ {
			copySlug, err = svc.ensureSlug(copyName, fmt.Sprintf("%s-%d", base, i), "")
			if err == nil {
				break
			}
		}
		if err != nil {
			return nil, err
		}
	}
	dup := *src
	dup.ID = newID()
	dup.Name = copyName
	dup.Slug = copySlug
	dup.Status = models.StatusDraft
	dup.CreatedAt = now
	dup.UpdatedAt = now
	dup.IconURL = src.IconURL
	dup.CoverURL = src.CoverURL
	if err := svc.repo.Create(&dup); err != nil {
		return nil, err
	}
	_ = svc.repo.SyncTags(dup.ID, dup.Tags)
	_ = svc.repo.SyncRulebooks(dup.ID, dup.CoreRulebooks)
	_ = svc.repo.SyncLinks(dup.ID, dup.OfficialLinks)
	return svc.repo.GetByID(dup.ID)
}

func (svc *SystemService) Fork(id string) (*models.GameSystem, error) {
	src, err := svc.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	forkName := src.Name + " (Fork)"
	forkSlug, err := svc.ensureSlug(forkName, "", "")
	if err != nil {
		base := normalizeSlug(forkName)
		for i := 2; i < 100; i++ {
			forkSlug, err = svc.ensureSlug(forkName, fmt.Sprintf("%s-%d", base, i), "")
			if err == nil {
				break
			}
		}
		if err != nil {
			return nil, err
		}
	}
	parent := src.ID
	dup := *src
	dup.ID = newID()
	dup.Name = forkName
	dup.Slug = forkSlug
	dup.Status = models.StatusDraft
	dup.IsCore = false
	dup.ParentSystemID = &parent
	dup.CreatedAt = now
	dup.UpdatedAt = now
	if err := svc.repo.Create(&dup); err != nil {
		return nil, err
	}
	_ = svc.repo.SyncTags(dup.ID, dup.Tags)
	_ = svc.repo.SyncRulebooks(dup.ID, dup.CoreRulebooks)
	_ = svc.repo.SyncLinks(dup.ID, dup.OfficialLinks)
	return svc.repo.GetByID(dup.ID)
}

func (svc *SystemService) SaveImage(systemID, kind string, data []byte) (*models.GameSystem, error) {
	s, err := svc.repo.GetByID(systemID)
	if err != nil {
		return nil, err
	}
	if s.Status == models.StatusLocked {
		return nil, ErrLocked
	}

	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("invalid image file")
	}

	var out image.Image
	var filename string
	dir := filepath.Join(svc.uploadDir, "systems", systemID)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}

	switch kind {
	case "icon":
		out = imaging.Fit(img, 512, 512, imaging.Lanczos)
		filename = "icon.png"
	case "cover":
		out = imaging.Fit(img, 1920, 1080, imaging.Lanczos)
		filename = "cover.jpg"
	default:
		return nil, fmt.Errorf("unknown image kind")
	}

	path := filepath.Join(dir, filename)
	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	if kind == "icon" {
		err = imaging.Encode(f, out, imaging.PNG)
	} else {
		err = imaging.Encode(f, out, imaging.JPEG)
	}
	f.Close()
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("/api/uploads/systems/%s/%s", systemID, filename)
	if kind == "icon" {
		s.IconURL = &url
	} else {
		s.CoverURL = &url
	}
	s.UpdatedAt = time.Now().UTC()
	if err := svc.repo.Update(s); err != nil {
		return nil, err
	}
	return svc.repo.GetByID(systemID)
}

func (svc *SystemService) Export(id string) ([]byte, error) {
	s, err := svc.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return json.MarshalIndent(s, "", "  ")
}

func (svc *SystemService) ImportJSON(data []byte) (*models.GameSystem, error) {
	if err := validateSystemImportJSON(data); err != nil {
		return nil, err
	}
	var payload models.GameSystem
	if err := json.Unmarshal(data, &payload); err != nil {
		return nil, fmt.Errorf("invalid JSON")
	}
	req := models.CreateSystemRequest{
		Name:            payload.Name,
		Slug:            payload.Slug,
		Version:         payload.Version,
		OfficialLinks:   payload.OfficialLinks,
		Tags:            payload.Tags,
		CoreRulebooks:   payload.CoreRulebooks,
	}
	if payload.Edition != nil {
		req.Edition = *payload.Edition
	}
	if payload.Publisher != nil {
		req.Publisher = *payload.Publisher
	}
	if payload.Description != nil {
		req.Description = *payload.Description
	}
	if payload.LicenseType != nil {
		req.LicenseType = *payload.LicenseType
	}
	if payload.Playstyle != nil {
		req.Playstyle = *payload.Playstyle
	}
	req.Complexity = payload.Complexity
	if payload.MeasurementUnit != nil {
		req.MeasurementUnit = *payload.MeasurementUnit
	}
	if payload.CurrencySymbol != nil {
		req.CurrencySymbol = *payload.CurrencySymbol
	}
	active := payload.IsActive
	req.IsActive = &active
	if payload.SystemFamily != nil {
		req.SystemFamily = *payload.SystemFamily
	}
	req.PlayerCountMin = payload.PlayerCountMin
	req.PlayerCountMax = payload.PlayerCountMax
	return svc.Create(req)
}

func (svc *SystemService) SaveTemplate(systemID, name, description string) error {
	data, err := svc.Export(systemID)
	if err != nil {
		return err
	}
	return svc.repo.SaveTemplate(newID(), name, description, data)
}

func (svc *SystemService) UploadDir() string {
	return svc.uploadDir
}

func ReadUploadFile(r io.Reader, maxBytes int64) ([]byte, error) {
	return io.ReadAll(io.LimitReader(r, maxBytes))
}
