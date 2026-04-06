package simbuildcore

import (
	"context"
	"os"
	"strconv"
	"strings"
	"sync"

	"gorm.io/gorm"

	"monstrolingo_backend/internal/appenv"
	"monstrolingo_backend/internal/db"
	"monstrolingo_backend/internal/db/models"
)

type Service struct {
	repo *Repository
}

type Repository struct {
	db *gorm.DB
}

var (
	serviceOnce sync.Once
	serviceInst *Service
	serviceErr  error
)

func GetService() (*Service, error) {
	serviceOnce.Do(func() {
		serviceInst, serviceErr = newService()
	})
	if serviceErr != nil {
		return nil, serviceErr
	}
	return serviceInst, nil
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func newService() (*Service, error) {
	appenv.Load()

	port, err := envInt("POSTGRES_PORT", 5435)
	if err != nil {
		return nil, invalidArgumentf("invalid POSTGRES_PORT: %v", err)
	}
	cfg := db.Config{
		Host:     envOr("POSTGRES_HOST", "localhost"),
		Port:     port,
		User:     envOr("POSTGRES_USER", "monstrolingo"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		DBName:   envOr("POSTGRES_DB", "monstrolingo"),
		SSLMode:  envOr("POSTGRES_SSLMODE", "disable"),
	}
	gormDB, err := db.Open(cfg)
	if err != nil {
		return nil, internalError("open database", err)
	}

	return &Service{
		repo: NewRepository(gormDB),
	}, nil
}

func (r *Repository) LanguageExists(ctx context.Context, code string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&models.Language{}).
		Where("LOWER(code) = LOWER(?) AND is_active = TRUE", code).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *Repository) ListActiveLanguageCodes(ctx context.Context) ([]string, error) {
	rows := make([]models.Language, 0, 16)
	if err := r.db.WithContext(ctx).
		Model(&models.Language{}).
		Where("is_active = TRUE").
		Order("sort_order ASC").
		Order("code ASC").
		Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]string, 0, len(rows))
	for _, row := range rows {
		out = append(out, normalizeLanguageCode(row.Code))
	}
	return out, nil
}

func (r *Repository) ListSkillTranslations(ctx context.Context) ([]skillTranslationRow, error) {
	rows := make([]skillTranslationRow, 0, 2048)
	if err := r.db.WithContext(ctx).
		Table("skills AS s").
		Select(strings.Join([]string{
			"s.id AS skill_id",
			"s.external_key AS external_key",
			"s.max_level AS max_level",
			"s.is_set_bonus_skill AS is_set_bonus",
			"l.code AS language_code",
			"st.name AS name",
			"st.effect_summary AS effect_summary",
		}, ", ")).
		Joins("JOIN skill_translations AS st ON st.skill_id = s.id").
		Joins("JOIN languages AS l ON l.id = st.language_id").
		Where("s.deleted_at IS NULL").
		Where("l.is_active = TRUE").
		Find(&rows).Error; err != nil {
		return nil, err
	}

	for i := range rows {
		rows[i].LanguageCode = normalizeLanguageCode(rows[i].LanguageCode)
	}
	return rows, nil
}

func envOr(key string, fallback string) string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	return v
}

func envInt(key string, fallback int) (int, error) {
	v := os.Getenv(key)
	if v == "" {
		return fallback, nil
	}
	return strconv.Atoi(v)
}
