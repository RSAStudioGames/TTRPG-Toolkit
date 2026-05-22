package repository

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/gabriel/ttrpg-toolkit/backend/internal/models"
	"github.com/jmoiron/sqlx"
)

var ErrNotFound = errors.New("not found")
var ErrSlugConflict = errors.New("slug conflict")

type SystemRepository struct {
	db *sqlx.DB
}

func NewSystemRepository(db *sqlx.DB) *SystemRepository {
	return &SystemRepository{db: db}
}

const systemColumns = `id, name, slug, edition, publisher, description, license_type, version,
	playstyle, complexity, measurement_unit, currency_symbol, status, is_active, system_family,
	player_count_min, player_count_max, official_links, tags, core_rulebooks, icon_url, cover_url,
	parent_system_id, is_core, is_protected, created_at, updated_at`

func decodeSystem(row *models.GameSystem) error {
	if len(row.OfficialLinksJSON) > 0 {
		_ = json.Unmarshal(row.OfficialLinksJSON, &row.OfficialLinks)
	}
	if row.OfficialLinks == nil {
		row.OfficialLinks = []string{}
	}
	if len(row.TagsJSON) > 0 {
		_ = json.Unmarshal(row.TagsJSON, &row.Tags)
	}
	if row.Tags == nil {
		row.Tags = []string{}
	}
	if len(row.CoreRulebooksJSON) > 0 {
		_ = json.Unmarshal(row.CoreRulebooksJSON, &row.CoreRulebooks)
	}
	if row.CoreRulebooks == nil {
		row.CoreRulebooks = []string{}
	}
	return nil
}

func encodeJSONStrings(items []string) ([]byte, error) {
	if items == nil {
		items = []string{}
	}
	return json.Marshal(items)
}

func (r *SystemRepository) GetByID(id string) (*models.GameSystem, error) {
	var s models.GameSystem
	q := `SELECT ` + systemColumns + ` FROM systems WHERE id = $1`
	if err := r.db.Get(&s, q, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	_ = decodeSystem(&s)
	return &s, nil
}

func (r *SystemRepository) GetBySlug(slug string) (*models.GameSystem, error) {
	var s models.GameSystem
	q := `SELECT ` + systemColumns + ` FROM systems WHERE slug = $1`
	if err := r.db.Get(&s, q, slug); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	_ = decodeSystem(&s)
	return &s, nil
}

type ListFilter struct {
	Page     int
	PerPage  int
	Status   string
	IsActive *bool
}

func (r *SystemRepository) List(f ListFilter) ([]models.GameSystem, int, error) {
	if f.Page < 1 {
		f.Page = 1
	}
	if f.PerPage < 1 || f.PerPage > 100 {
		f.PerPage = 20
	}
	offset := (f.Page - 1) * f.PerPage

	var conditions []string
	var args []interface{}
	n := 1
	if f.Status != "" {
		conditions = append(conditions, fmt.Sprintf("status = $%d", n))
		args = append(args, f.Status)
		n++
	}
	if f.IsActive != nil {
		conditions = append(conditions, fmt.Sprintf("is_active = $%d", n))
		args = append(args, *f.IsActive)
		n++
	}
	where := ""
	if len(conditions) > 0 {
		where = "WHERE " + strings.Join(conditions, " AND ")
	}

	countQ := `SELECT COUNT(*) FROM systems ` + where
	var total int
	if err := r.db.Get(&total, countQ, args...); err != nil {
		return nil, 0, err
	}

	listArgs := append(args, f.PerPage, offset)
	q := `SELECT ` + systemColumns + ` FROM systems ` + where +
		fmt.Sprintf(` ORDER BY updated_at DESC LIMIT $%d OFFSET $%d`, n, n+1)
	rows := []models.GameSystem{}
	if err := r.db.Select(&rows, q, listArgs...); err != nil {
		return nil, 0, err
	}
	for i := range rows {
		_ = decodeSystem(&rows[i])
	}
	return rows, total, nil
}

func (r *SystemRepository) Create(s *models.GameSystem) error {
	links, _ := encodeJSONStrings(s.OfficialLinks)
	tags, _ := encodeJSONStrings(s.Tags)
	books, _ := encodeJSONStrings(s.CoreRulebooks)

	q := `INSERT INTO systems (` + systemColumns + `) VALUES (
		:id, :name, :slug, :edition, :publisher, :description, :license_type, :version,
		:playstyle, :complexity, :measurement_unit, :currency_symbol, :status, :is_active, :system_family,
		:player_count_min, :player_count_max, :official_links, :tags, :core_rulebooks, :icon_url, :cover_url,
		:parent_system_id, :is_core, :is_protected, :created_at, :updated_at)`

	params := map[string]interface{}{
		"id": s.ID, "name": s.Name, "slug": s.Slug, "edition": s.Edition, "publisher": s.Publisher,
		"description": s.Description, "license_type": s.LicenseType, "version": s.Version,
		"playstyle": s.Playstyle, "complexity": s.Complexity, "measurement_unit": s.MeasurementUnit,
		"currency_symbol": s.CurrencySymbol, "status": s.Status, "is_active": s.IsActive,
		"system_family": s.SystemFamily, "player_count_min": s.PlayerCountMin, "player_count_max": s.PlayerCountMax,
		"official_links": links, "tags": tags, "core_rulebooks": books,
		"icon_url": s.IconURL, "cover_url": s.CoverURL, "parent_system_id": s.ParentSystemID,
		"is_core": s.IsCore, "is_protected": s.IsProtected, "created_at": s.CreatedAt, "updated_at": s.UpdatedAt,
	}
	_, err := r.db.NamedExec(q, params)
	if err != nil && strings.Contains(err.Error(), "duplicate key") && strings.Contains(err.Error(), "slug") {
		return ErrSlugConflict
	}
	return err
}

func (r *SystemRepository) Update(s *models.GameSystem) error {
	links, _ := encodeJSONStrings(s.OfficialLinks)
	tags, _ := encodeJSONStrings(s.Tags)
	books, _ := encodeJSONStrings(s.CoreRulebooks)

	q := `UPDATE systems SET
		name = :name, slug = :slug, edition = :edition, publisher = :publisher, description = :description,
		license_type = :license_type, version = :version, playstyle = :playstyle, complexity = :complexity,
		measurement_unit = :measurement_unit, currency_symbol = :currency_symbol, status = :status,
		is_active = :is_active, system_family = :system_family, player_count_min = :player_count_min,
		player_count_max = :player_count_max, official_links = :official_links, tags = :tags,
		core_rulebooks = :core_rulebooks, icon_url = :icon_url, cover_url = :cover_url,
		parent_system_id = :parent_system_id, is_core = :is_core, is_protected = :is_protected,
		updated_at = :updated_at
		WHERE id = :id`

	params := map[string]interface{}{
		"id": s.ID, "name": s.Name, "slug": s.Slug, "edition": s.Edition, "publisher": s.Publisher,
		"description": s.Description, "license_type": s.LicenseType, "version": s.Version,
		"playstyle": s.Playstyle, "complexity": s.Complexity, "measurement_unit": s.MeasurementUnit,
		"currency_symbol": s.CurrencySymbol, "status": s.Status, "is_active": s.IsActive,
		"system_family": s.SystemFamily, "player_count_min": s.PlayerCountMin, "player_count_max": s.PlayerCountMax,
		"official_links": links, "tags": tags, "core_rulebooks": books,
		"icon_url": s.IconURL, "cover_url": s.CoverURL, "parent_system_id": s.ParentSystemID,
		"is_core": s.IsCore, "is_protected": s.IsProtected, "updated_at": s.UpdatedAt,
	}
	res, err := r.db.NamedExec(q, params)
	if err != nil && strings.Contains(err.Error(), "duplicate key") && strings.Contains(err.Error(), "slug") {
		return ErrSlugConflict
	}
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *SystemRepository) Delete(id string) error {
	res, err := r.db.Exec(`DELETE FROM systems WHERE id = $1`, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *SystemRepository) SlugExists(slug string, excludeID string) (bool, error) {
	var count int
	q := `SELECT COUNT(*) FROM systems WHERE slug = $1`
	args := []interface{}{slug}
	if excludeID != "" {
		q += ` AND id != $2`
		args = append(args, excludeID)
	}
	if err := r.db.Get(&count, q, args...); err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *SystemRepository) CountActiveCampaignsBySystemID(_ string) (int, error) {
	return 0, nil
}

func (r *SystemRepository) CountChildren(parentID string) (int, error) {
	var n int
	err := r.db.Get(&n, `SELECT COUNT(*) FROM systems WHERE parent_system_id = $1`, parentID)
	return n, err
}

func (r *SystemRepository) DeletePreview(id string) (*models.DeletePreview, error) {
	var tagCount, bookCount, linkCount, childCount int
	_ = r.db.Get(&tagCount, `SELECT COUNT(*) FROM system_tags WHERE system_id = $1`, id)
	_ = r.db.Get(&bookCount, `SELECT COUNT(*) FROM system_rulebooks WHERE system_id = $1`, id)
	_ = r.db.Get(&linkCount, `SELECT COUNT(*) FROM system_links WHERE system_id = $1`, id)
	_ = r.db.Get(&childCount, `SELECT COUNT(*) FROM systems WHERE parent_system_id = $1`, id)
	// Fallback JSONB counts if relational empty
	var s models.GameSystem
	if err := r.db.Get(&s, `SELECT official_links, tags, core_rulebooks FROM systems WHERE id = $1`, id); err == nil {
		_ = decodeSystem(&s)
		if tagCount == 0 {
			tagCount = len(s.Tags)
		}
		if bookCount == 0 {
			bookCount = len(s.CoreRulebooks)
		}
		if linkCount == 0 {
			linkCount = len(s.OfficialLinks)
		}
	}
	total := tagCount + bookCount + linkCount + childCount
	return &models.DeletePreview{
		TagCount: tagCount, RulebookCount: bookCount, LinkCount: linkCount,
		ChildCount: childCount, TotalAssociated: total,
	}, nil
}

func (r *SystemRepository) SyncTags(systemID string, tags []string) error {
	_, _ = r.db.Exec(`DELETE FROM system_tags WHERE system_id = $1`, systemID)
	for _, t := range tags {
		if t == "" {
			continue
		}
		_, err := r.db.Exec(`INSERT INTO system_tags (system_id, tag) VALUES ($1, $2) ON CONFLICT DO NOTHING`, systemID, t)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *SystemRepository) SyncRulebooks(systemID string, books []string) error {
	_, _ = r.db.Exec(`DELETE FROM system_rulebooks WHERE system_id = $1`, systemID)
	for _, b := range books {
		if b == "" {
			continue
		}
		_, err := r.db.Exec(`INSERT INTO system_rulebooks (system_id, rulebook) VALUES ($1, $2) ON CONFLICT DO NOTHING`, systemID, b)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *SystemRepository) SyncLinks(systemID string, urls []string) error {
	_, _ = r.db.Exec(`DELETE FROM system_links WHERE system_id = $1`, systemID)
	for i, u := range urls {
		if u == "" {
			continue
		}
		_, err := r.db.Exec(`INSERT INTO system_links (system_id, url, sort_order) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`, systemID, u, i)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *SystemRepository) SaveTemplate(id, name, description string, schemaJSON []byte) error {
	_, err := r.db.Exec(
		`INSERT INTO system_templates (id, name, description, schema_json) VALUES ($1, $2, $3, $4)`,
		id, name, description, schemaJSON,
	)
	return err
}
