package catalog

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"sync"

	"monstrolingo_backend/internal/appenv"
	"monstrolingo_backend/internal/db"
)

type Service struct {
	repo *Repository
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
	return &Service{repo: NewRepository(gormDB)}, nil
}

func (s *Service) ListCategoryTable(ctx context.Context, key CategoryKey, req *CategoryTableRequest) (*CategoryTableResponse, error) {
	query, err := normalizeTableQuery(req)
	if err != nil {
		return nil, err
	}
	if err := s.ensureLanguageSupported(ctx, query.SourceLang); err != nil {
		return nil, err
	}
	if err := s.ensureLanguageSupported(ctx, query.TargetLang); err != nil {
		return nil, err
	}
	out, err := s.repo.ListCategoryTable(ctx, key, query)
	if err != nil {
		return nil, internalError(fmt.Sprintf("list %s table", key), err)
	}
	return out, nil
}

func (s *Service) normalizeAndValidateDetail(ctx context.Context, req *CategoryDetailRequest) (normalizedDetailQuery, error) {
	query, err := normalizeDetailQuery(req)
	if err != nil {
		return normalizedDetailQuery{}, err
	}
	if err := s.ensureLanguageSupported(ctx, query.TargetLang); err != nil {
		return normalizedDetailQuery{}, err
	}
	if err := s.ensureLanguageSupported(ctx, "en"); err != nil {
		return normalizedDetailQuery{}, err
	}
	return query, nil
}

func (s *Service) GetItemDetail(ctx context.Context, req *CategoryDetailRequest) (*ItemDetailResponse, error) {
	query, err := s.normalizeAndValidateDetail(ctx, req)
	if err != nil {
		return nil, err
	}
	out, err := s.repo.GetItemDetail(ctx, query.ExternalKey, query.TargetLang)
	return out, s.mapDetailError("item", query.ExternalKey, err)
}

func (s *Service) GetWeaponDetail(ctx context.Context, req *CategoryDetailRequest) (*WeaponDetailResponse, error) {
	query, err := s.normalizeAndValidateDetail(ctx, req)
	if err != nil {
		return nil, err
	}
	out, err := s.repo.GetWeaponDetail(ctx, query.ExternalKey, query.TargetLang)
	return out, s.mapDetailError("weapon", query.ExternalKey, err)
}

func (s *Service) GetArmorDetail(ctx context.Context, req *CategoryDetailRequest) (*ArmorDetailResponse, error) {
	query, err := s.normalizeAndValidateDetail(ctx, req)
	if err != nil {
		return nil, err
	}
	out, err := s.repo.GetArmorDetail(ctx, query.ExternalKey, query.TargetLang)
	return out, s.mapDetailError("armor", query.ExternalKey, err)
}

func (s *Service) GetSkillDetail(ctx context.Context, req *CategoryDetailRequest) (*SkillDetailResponse, error) {
	query, err := s.normalizeAndValidateDetail(ctx, req)
	if err != nil {
		return nil, err
	}
	out, err := s.repo.GetSkillDetail(ctx, query.ExternalKey, query.TargetLang)
	return out, s.mapDetailError("skill", query.ExternalKey, err)
}

func (s *Service) GetDecorationDetail(ctx context.Context, req *CategoryDetailRequest) (*DecorationDetailResponse, error) {
	query, err := s.normalizeAndValidateDetail(ctx, req)
	if err != nil {
		return nil, err
	}
	out, err := s.repo.GetDecorationDetail(ctx, query.ExternalKey, query.TargetLang)
	return out, s.mapDetailError("decoration", query.ExternalKey, err)
}

func (s *Service) GetCharmDetail(ctx context.Context, req *CategoryDetailRequest) (*CharmDetailResponse, error) {
	query, err := s.normalizeAndValidateDetail(ctx, req)
	if err != nil {
		return nil, err
	}
	out, err := s.repo.GetCharmDetail(ctx, query.ExternalKey, query.TargetLang)
	return out, s.mapDetailError("charm", query.ExternalKey, err)
}

func (s *Service) GetFoodSkillDetail(ctx context.Context, req *CategoryDetailRequest) (*FoodSkillDetailResponse, error) {
	query, err := s.normalizeAndValidateDetail(ctx, req)
	if err != nil {
		return nil, err
	}
	out, err := s.repo.GetFoodSkillDetail(ctx, query.ExternalKey, query.TargetLang)
	return out, s.mapDetailError("food skill", query.ExternalKey, err)
}

func (s *Service) GetKinsectDetail(ctx context.Context, req *CategoryDetailRequest) (*KinsectDetailResponse, error) {
	query, err := s.normalizeAndValidateDetail(ctx, req)
	if err != nil {
		return nil, err
	}
	out, err := s.repo.GetKinsectDetail(ctx, query.ExternalKey, query.TargetLang)
	return out, s.mapDetailError("kinsect", query.ExternalKey, err)
}

func (s *Service) ListLanguages(ctx context.Context) (*LanguagesResponse, error) {
	languages, err := s.repo.ListLanguages(ctx)
	if err != nil {
		return nil, internalError("list active languages", err)
	}
	return &LanguagesResponse{Languages: languages}, nil
}

func (s *Service) mapDetailError(entity string, key string, err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, ErrNotFound) {
		return notFoundf(entity, key)
	}
	return internalError(fmt.Sprintf("load %s detail", entity), err)
}

func (s *Service) ensureLanguageSupported(ctx context.Context, code string) error {
	ok, err := s.repo.LanguageExists(ctx, code)
	if err != nil {
		return internalError(fmt.Sprintf("check language %q", code), err)
	}
	if !ok {
		return unsupportedLanguage(code)
	}
	return nil
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
