package catalog

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"monstrolingo_backend/internal/db/models"
)

var ErrNotFound = errors.New("not found")

type Repository struct {
	db *gorm.DB
}

type categoryTableDBRow struct {
	ExternalKey       string
	SourceName        string
	SourceDescription string
	TargetName        string
	TargetDescription string
}

type translationRecord struct {
	Name          string
	Description   string
	FlavorText    string
	EffectSummary string
	Slug          string
	LanguageCode  string
}

type associatedJewelRow struct {
	DecorationExternalKey string
	Name                  string
	SlotSize              int16
	Rarity                int16
	SkillLevel            int16
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) LanguageExists(ctx context.Context, code string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&models.Language{}).
		Where("code = ? AND is_active = TRUE", code).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *Repository) ListLanguages(ctx context.Context) ([]LanguageOption, error) {
	rows := make([]models.Language, 0, 16)
	if err := r.db.WithContext(ctx).
		Model(&models.Language{}).
		Where("is_active = TRUE").
		Order("sort_order ASC").
		Order("code ASC").
		Find(&rows).Error; err != nil {
		return nil, err
	}

	out := make([]LanguageOption, 0, len(rows))
	for _, row := range rows {
		label := strings.TrimSpace(row.Label)
		if label == "" {
			label = row.Code
		}
		out = append(out, LanguageOption{
			Code:  row.Code,
			Label: label,
		})
	}
	return out, nil
}

func (r *Repository) ListCategoryTable(
	ctx context.Context,
	key CategoryKey,
	query normalizedTableQuery,
) (*CategoryTableResponse, error) {
	spec, err := getCategorySpec(key)
	if err != nil {
		return nil, err
	}

	var total int64
	if err := r.db.WithContext(ctx).
		Table(spec.CanonicalTable + " AS c").
		Where("c.deleted_at IS NULL").
		Count(&total).Error; err != nil {
		return nil, fmt.Errorf("count %s: %w", spec.CanonicalTable, err)
	}

	rows := make([]categoryTableDBRow, 0, query.Limit)
	selectSQL := strings.Join([]string{
		"c.external_key AS external_key",
		"COALESCE(src.name, '') AS source_name",
		"COALESCE(src.description, '') AS source_description",
		"COALESCE(tgt.name, '') AS target_name",
		"COALESCE(tgt.description, '') AS target_description",
	}, ", ")

	sourceJoin := fmt.Sprintf("LEFT JOIN %s AS src ON src.%s = c.id AND src.language_id = ls.id", spec.TranslationTable, spec.TranslationFK)
	targetJoin := fmt.Sprintf("LEFT JOIN %s AS tgt ON tgt.%s = c.id AND tgt.language_id = lt.id", spec.TranslationTable, spec.TranslationFK)
	orderBy := fmt.Sprintf("c.%s ASC, c.external_key ASC", spec.SortColumn)

	if err := r.db.WithContext(ctx).
		Table(spec.CanonicalTable+" AS c").
		Joins("JOIN languages AS ls ON ls.code = ? AND ls.is_active = TRUE", query.SourceLang).
		Joins("JOIN languages AS lt ON lt.code = ? AND lt.is_active = TRUE", query.TargetLang).
		Joins(sourceJoin).
		Joins(targetJoin).
		Where("c.deleted_at IS NULL").
		Select(selectSQL).
		Order(orderBy).
		Offset(query.Offset).
		Limit(query.Limit).
		Scan(&rows).Error; err != nil {
		return nil, fmt.Errorf("list %s: %w", spec.CanonicalTable, err)
	}

	out := make([]CategoryTableRow, 0, len(rows))
	for _, row := range rows {
		out = append(out, toTableRow(row, query.SourceLang, query.TargetLang))
	}

	return &CategoryTableResponse{
		Items:      out,
		Pagination: buildPagination(query.Page, query.Limit, total),
	}, nil
}

func (r *Repository) GetItemDetail(ctx context.Context, externalKey string, targetLang string) (*ItemDetailResponse, error) {
	var row models.Item
	if err := r.db.WithContext(ctx).
		Where("external_key = ?", externalKey).
		Where("deleted_at IS NULL").
		First(&row).Error; err != nil {
		return nil, wrapNotFound("item", externalKey, err)
	}
	spec, _ := getCategorySpec(CategoryItems)
	translation, err := r.loadDetailTranslation(ctx, spec, row.ID, targetLang)
	if err != nil {
		return nil, err
	}

	return &ItemDetailResponse{
		ExternalKey:        row.ExternalKey,
		SyncedAt:           row.SyncedAt,
		Category:           row.Category,
		Subcategory:        row.Subcategory,
		Rarity:             row.Rarity,
		CarryLimit:         row.CarryLimit,
		BuyPrice:           row.BuyPrice,
		SellPrice:          row.SellPrice,
		Points:             row.Points,
		IsCraftingMaterial: row.IsCraftingMaterial,
		IsConsumable:       row.IsConsumable,
		IconKey:            row.IconKey,
		SortOrder:          row.SortOrder,
		Translation:        translation,
	}, nil
}

func (r *Repository) GetWeaponDetail(ctx context.Context, externalKey string, targetLang string) (*WeaponDetailResponse, error) {
	var row models.Weapon
	if err := r.db.WithContext(ctx).
		Preload("Skills", func(tx *gorm.DB) *gorm.DB { return tx.Order("sort_order ASC") }).
		Where("external_key = ?", externalKey).
		Where("deleted_at IS NULL").
		First(&row).Error; err != nil {
		return nil, wrapNotFound("weapon", externalKey, err)
	}

	spec, _ := getCategorySpec(CategoryWeapons)
	translation, err := r.loadDetailTranslation(ctx, spec, row.ID, targetLang)
	if err != nil {
		return nil, err
	}

	skills := make([]SkillLinkDetail, 0, len(row.Skills))
	for _, rel := range row.Skills {
		skill, err := r.loadSkillLinkBase(ctx, rel.SkillID, targetLang)
		if err != nil {
			continue
		}
		skill.Level = rel.Level
		skill.SortOrder = rel.SortOrder
		skill.SourceType = rel.SourceType
		skills = append(skills, skill)
	}
	setSkills, regularSkills := splitSkillsBySetType(skills)

	return &WeaponDetailResponse{
		ExternalKey:     row.ExternalKey,
		SyncedAt:        row.SyncedAt,
		WeaponType:      row.WeaponType,
		Rarity:          row.Rarity,
		Attack:          row.Attack,
		AffinityPercent: row.AffinityPercent,
		DefenseBonus:    row.DefenseBonus,
		Slot1Level:      row.Slot1Level,
		Slot2Level:      row.Slot2Level,
		Slot3Level:      row.Slot3Level,
		ElementType:     row.ElementType,
		ElementValue:    row.ElementValue,
		AilmentType:     row.AilmentType,
		AilmentValue:    row.AilmentValue,
		SharpnessRed:    row.SharpnessRed,
		SharpnessOrange: row.SharpnessOrange,
		SharpnessYellow: row.SharpnessYellow,
		SharpnessGreen:  row.SharpnessGreen,
		SharpnessBlue:   row.SharpnessBlue,
		SharpnessWhite:  row.SharpnessWhite,
		SharpnessPurple: row.SharpnessPurple,
		CraftCost:       row.CraftCost,
		UpgradeCost:     row.UpgradeCost,
		TreeDepth:       row.TreeDepth,
		IsFinalUpgrade:  row.IsFinalUpgrade,
		SortOrder:       row.SortOrder,
		Translation:     translation,
		Skills:          skills,
		SetSkills:       setSkills,
		RegularSkills:   regularSkills,
	}, nil
}

func (r *Repository) GetArmorDetail(ctx context.Context, externalKey string, targetLang string) (*ArmorDetailResponse, error) {
	var row models.ArmorPiece
	if err := r.db.WithContext(ctx).
		Preload("Skills", func(tx *gorm.DB) *gorm.DB { return tx.Order("sort_order ASC") }).
		Where("external_key = ?", externalKey).
		Where("deleted_at IS NULL").
		First(&row).Error; err != nil {
		return nil, wrapNotFound("armor", externalKey, err)
	}

	spec, _ := getCategorySpec(CategoryArmor)
	translation, err := r.loadDetailTranslation(ctx, spec, row.ID, targetLang)
	if err != nil {
		return nil, err
	}

	skills := make([]SkillLinkDetail, 0, len(row.Skills))
	for _, rel := range row.Skills {
		skill, err := r.loadSkillLinkBase(ctx, rel.SkillID, targetLang)
		if err != nil {
			continue
		}
		skill.Level = rel.Level
		skill.SortOrder = rel.SortOrder
		skills = append(skills, skill)
	}
	setSkills, regularSkills := splitSkillsBySetType(skills)

	return &ArmorDetailResponse{
		ExternalKey:         row.ExternalKey,
		SyncedAt:            row.SyncedAt,
		ArmorSetKey:         row.ArmorSetKey,
		ArmorSetName:        row.ArmorSetName,
		ArmorSetVariant:     row.ArmorSetVariant,
		PieceType:           row.PieceType,
		Rank:                row.Rank,
		Gender:              row.Gender,
		Rarity:              row.Rarity,
		DefenseBase:         row.DefenseBase,
		DefenseMax:          row.DefenseMax,
		DefenseAugmentedMax: row.DefenseAugmentedMax,
		FireRes:             row.FireRes,
		WaterRes:            row.WaterRes,
		ThunderRes:          row.ThunderRes,
		IceRes:              row.IceRes,
		DragonRes:           row.DragonRes,
		Slot1Level:          row.Slot1Level,
		Slot2Level:          row.Slot2Level,
		Slot3Level:          row.Slot3Level,
		IsLayered:           row.IsLayered,
		SortOrder:           row.SortOrder,
		Translation:         translation,
		Skills:              skills,
		SetSkills:           setSkills,
		RegularSkills:       regularSkills,
	}, nil
}

func (r *Repository) GetSkillDetail(ctx context.Context, externalKey string, targetLang string) (*SkillDetailResponse, error) {
	var row models.Skill
	if err := r.db.WithContext(ctx).
		Preload("Levels", func(tx *gorm.DB) *gorm.DB { return tx.Order("level ASC") }).
		Where("external_key = ?", externalKey).
		Where("deleted_at IS NULL").
		First(&row).Error; err != nil {
		return nil, wrapNotFound("skill", externalKey, err)
	}

	spec, _ := getCategorySpec(CategorySkills)
	translation, err := r.loadDetailTranslation(ctx, spec, row.ID, targetLang)
	if err != nil {
		return nil, err
	}

	levels := make([]SkillLevelDetail, 0, len(row.Levels))
	for _, lvl := range row.Levels {
		levels = append(levels, SkillLevelDetail{
			Level:           lvl.Level,
			EffectValueText: lvl.EffectValueText,
			Description:     lvl.Description,
		})
	}
	associatedJewels, err := r.loadAssociatedJewelsForSkill(ctx, row.ID, targetLang)
	if err != nil {
		return nil, err
	}

	return &SkillDetailResponse{
		ExternalKey:      row.ExternalKey,
		SyncedAt:         row.SyncedAt,
		SkillKind:        row.SkillKind,
		MaxLevel:         row.MaxLevel,
		IsBinary:         row.IsBinary,
		IsSetBonusSkill:  row.IsSetBonusSkill,
		SortOrder:        row.SortOrder,
		Translation:      translation,
		Levels:           levels,
		AssociatedJewels: associatedJewels,
	}, nil
}

func (r *Repository) loadAssociatedJewelsForSkill(
	ctx context.Context,
	skillID uuid.UUID,
	targetLang string,
) ([]AssociatedJewelDetail, error) {
	rows := make([]associatedJewelRow, 0, 16)
	selectSQL := strings.Join([]string{
		"d.external_key AS decoration_external_key",
		"d.slot_size AS slot_size",
		"d.rarity AS rarity",
		"ds.level AS skill_level",
		"COALESCE(NULLIF(TRIM(tgt.name), ''), NULLIF(TRIM(en.name), ''), d.external_key) AS name",
	}, ", ")
	if err := r.db.WithContext(ctx).
		Table("decoration_skills AS ds").
		Select(selectSQL).
		Joins("JOIN decorations AS d ON d.id = ds.decoration_id").
		Joins("LEFT JOIN languages AS lt ON lt.code = ? AND lt.is_active = TRUE", targetLang).
		Joins("LEFT JOIN decorations_translations AS tgt ON tgt.decoration_id = d.id AND tgt.language_id = lt.id").
		Joins("LEFT JOIN languages AS le ON le.code = ? AND le.is_active = TRUE", "en").
		Joins("LEFT JOIN decorations_translations AS en ON en.decoration_id = d.id AND en.language_id = le.id").
		Where("ds.skill_id = ?", skillID).
		Where("d.deleted_at IS NULL").
		Order("d.slot_size ASC").
		Order("ds.level DESC").
		Order("COALESCE(NULLIF(TRIM(tgt.name), ''), NULLIF(TRIM(en.name), ''), d.external_key) ASC").
		Order("d.external_key ASC").
		Scan(&rows).Error; err != nil {
		return nil, err
	}

	out := make([]AssociatedJewelDetail, 0, len(rows))
	for _, row := range rows {
		out = append(out, AssociatedJewelDetail{
			DecorationExternalKey: row.DecorationExternalKey,
			Name:                  strings.TrimSpace(row.Name),
			SlotSize:              row.SlotSize,
			Rarity:                row.Rarity,
			SkillLevel:            row.SkillLevel,
		})
	}
	return out, nil
}

func (r *Repository) GetDecorationDetail(ctx context.Context, externalKey string, targetLang string) (*DecorationDetailResponse, error) {
	var row models.Decoration
	if err := r.db.WithContext(ctx).
		Preload("Skills", func(tx *gorm.DB) *gorm.DB { return tx.Order("sort_order ASC") }).
		Where("external_key = ?", externalKey).
		Where("deleted_at IS NULL").
		First(&row).Error; err != nil {
		return nil, wrapNotFound("decoration", externalKey, err)
	}

	spec, _ := getCategorySpec(CategoryDecorations)
	translation, err := r.loadDetailTranslation(ctx, spec, row.ID, targetLang)
	if err != nil {
		return nil, err
	}

	skills := make([]SkillLinkDetail, 0, len(row.Skills))
	for _, rel := range row.Skills {
		skill, err := r.loadSkillLinkBase(ctx, rel.SkillID, targetLang)
		if err != nil {
			continue
		}
		skill.Level = rel.Level
		skill.SortOrder = rel.SortOrder
		skills = append(skills, skill)
	}

	return &DecorationDetailResponse{
		ExternalKey: row.ExternalKey,
		SyncedAt:    row.SyncedAt,
		SlotSize:    row.SlotSize,
		Rarity:      row.Rarity,
		IsCraftable: row.IsCraftable,
		CraftCost:   row.CraftCost,
		SortOrder:   row.SortOrder,
		Translation: translation,
		Skills:      skills,
	}, nil
}

func (r *Repository) GetCharmDetail(ctx context.Context, externalKey string, targetLang string) (*CharmDetailResponse, error) {
	var row models.Charm
	if err := r.db.WithContext(ctx).
		Preload("Skills", func(tx *gorm.DB) *gorm.DB { return tx.Order("rank_required ASC, sort_order ASC") }).
		Where("external_key = ?", externalKey).
		Where("deleted_at IS NULL").
		First(&row).Error; err != nil {
		return nil, wrapNotFound("charm", externalKey, err)
	}

	spec, _ := getCategorySpec(CategoryCharms)
	translation, err := r.loadDetailTranslation(ctx, spec, row.ID, targetLang)
	if err != nil {
		return nil, err
	}

	skills := make([]SkillLinkDetail, 0, len(row.Skills))
	for _, rel := range row.Skills {
		skill, err := r.loadSkillLinkBase(ctx, rel.SkillID, targetLang)
		if err != nil {
			continue
		}
		skill.Level = rel.Level
		skill.SortOrder = rel.SortOrder
		skill.RankRequired = rel.RankRequired
		skills = append(skills, skill)
	}

	return &CharmDetailResponse{
		ExternalKey: row.ExternalKey,
		SyncedAt:    row.SyncedAt,
		Rarity:      row.Rarity,
		MaxRank:     row.MaxRank,
		CraftCost:   row.CraftCost,
		SortOrder:   row.SortOrder,
		Translation: translation,
		Skills:      skills,
	}, nil
}

func (r *Repository) GetFoodSkillDetail(ctx context.Context, externalKey string, targetLang string) (*FoodSkillDetailResponse, error) {
	var row models.FoodSkill
	if err := r.db.WithContext(ctx).
		Preload("Levels", func(tx *gorm.DB) *gorm.DB { return tx.Order("level ASC") }).
		Where("external_key = ?", externalKey).
		Where("deleted_at IS NULL").
		First(&row).Error; err != nil {
		return nil, wrapNotFound("food skill", externalKey, err)
	}

	spec, _ := getCategorySpec(CategoryFoodSkills)
	translation, err := r.loadDetailTranslation(ctx, spec, row.ID, targetLang)
	if err != nil {
		return nil, err
	}

	levels := make([]FoodSkillLevelDetail, 0, len(row.Levels))
	for _, lvl := range row.Levels {
		levels = append(levels, FoodSkillLevelDetail{
			Level:             lvl.Level,
			DurationSeconds:   lvl.DurationSeconds,
			ActivationPercent: lvl.ActivationPercent,
			EffectValueText:   lvl.EffectValueText,
		})
	}

	return &FoodSkillDetailResponse{
		ExternalKey:           row.ExternalKey,
		SyncedAt:              row.SyncedAt,
		FoodCategory:          row.FoodCategory,
		MaxLevel:              row.MaxLevel,
		BaseDurationSeconds:   row.BaseDurationSeconds,
		BaseActivationPercent: row.BaseActivationPercent,
		SortOrder:             row.SortOrder,
		Translation:           translation,
		Levels:                levels,
	}, nil
}

func (r *Repository) GetKinsectDetail(ctx context.Context, externalKey string, targetLang string) (*KinsectDetailResponse, error) {
	var row models.Kinsect
	if err := r.db.WithContext(ctx).
		Where("external_key = ?", externalKey).
		Where("deleted_at IS NULL").
		First(&row).Error; err != nil {
		return nil, wrapNotFound("kinsect", externalKey, err)
	}

	spec, _ := getCategorySpec(CategoryKinsects)
	translation, err := r.loadDetailTranslation(ctx, spec, row.ID, targetLang)
	if err != nil {
		return nil, err
	}

	return &KinsectDetailResponse{
		ExternalKey:           row.ExternalKey,
		SyncedAt:              row.SyncedAt,
		KinsectType:           row.KinsectType,
		AttackType:            row.AttackType,
		PowderType:            row.PowderType,
		KinsectBonusPrimary:   row.KinsectBonusPrimary,
		KinsectBonusSecondary: row.KinsectBonusSecondary,
		Rarity:                row.Rarity,
		PowerValue:            row.PowerValue,
		SpeedValue:            row.SpeedValue,
		HealValue:             row.HealValue,
		StaminaValue:          row.StaminaValue,
		ElementType:           row.ElementType,
		ElementValue:          row.ElementValue,
		SortOrder:             row.SortOrder,
		Translation:           translation,
	}, nil
}

func (r *Repository) loadSkillLinkBase(ctx context.Context, skillID uuid.UUID, targetLang string) (SkillLinkDetail, error) {
	var row models.Skill
	if err := r.db.WithContext(ctx).
		Where("id = ?", skillID).
		Where("deleted_at IS NULL").
		First(&row).Error; err != nil {
		return SkillLinkDetail{}, wrapNotFound("skill", skillID.String(), err)
	}
	spec, _ := getCategorySpec(CategorySkills)
	tr, err := r.findTranslationByEntityID(ctx, spec, row.ID, targetLang)
	if err != nil && !errors.Is(err, ErrNotFound) {
		return SkillLinkDetail{}, err
	}
	return SkillLinkDetail{
		SkillExternalKey: row.ExternalKey,
		Name:             tr.Name,
		Description:      tr.Description,
		EffectSummary:    tr.EffectSummary,
		SkillKind:        row.SkillKind,
		IsSetBonusSkill:  row.IsSetBonusSkill,
	}, nil
}

func splitSkillsBySetType(skills []SkillLinkDetail) ([]SkillLinkDetail, []SkillLinkDetail) {
	if len(skills) == 0 {
		return nil, nil
	}
	setSkills := make([]SkillLinkDetail, 0, len(skills))
	regularSkills := make([]SkillLinkDetail, 0, len(skills))
	for _, skill := range skills {
		if skill.IsSetBonusSkill {
			setSkills = append(setSkills, skill)
			continue
		}
		regularSkills = append(regularSkills, skill)
	}
	return setSkills, regularSkills
}

func (r *Repository) loadDetailTranslation(ctx context.Context, spec categorySpec, entityID uuid.UUID, targetLang string) (DetailTranslation, error) {
	row, err := r.findTranslationByEntityID(ctx, spec, entityID, targetLang)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return DetailTranslation{
				LanguageCode: targetLang,
			}, nil
		}
		return DetailTranslation{}, err
	}
	return toDetailTranslation(row), nil
}

func (r *Repository) findTranslationByEntityID(
	ctx context.Context,
	spec categorySpec,
	entityID uuid.UUID,
	targetLang string,
) (translationRecord, error) {
	row, err := r.findTranslationForLanguage(ctx, spec, entityID, targetLang)
	if err == nil {
		return row, nil
	}
	if !errors.Is(err, ErrNotFound) {
		return translationRecord{}, err
	}
	if targetLang == "en" {
		return translationRecord{}, ErrNotFound
	}
	return r.findTranslationForLanguage(ctx, spec, entityID, "en")
}

func (r *Repository) findTranslationForLanguage(
	ctx context.Context,
	spec categorySpec,
	entityID uuid.UUID,
	languageCode string,
) (translationRecord, error) {
	selectParts := []string{
		"COALESCE(t.name, '') AS name",
		"COALESCE(t.description, '') AS description",
		"COALESCE(t.slug, '') AS slug",
		"l.code AS language_code",
	}
	if spec.HasFlavorText {
		selectParts = append(selectParts, "COALESCE(t.flavor_text, '') AS flavor_text")
	} else {
		selectParts = append(selectParts, "'' AS flavor_text")
	}
	if spec.HasEffectSummary {
		selectParts = append(selectParts, "COALESCE(t.effect_summary, '') AS effect_summary")
	} else {
		selectParts = append(selectParts, "'' AS effect_summary")
	}

	var out translationRecord
	tx := r.db.WithContext(ctx).
		Table(spec.TranslationTable+" AS t").
		Select(strings.Join(selectParts, ", ")).
		Joins("JOIN languages AS l ON l.id = t.language_id").
		Where("t."+spec.TranslationFK+" = ?", entityID).
		Where("l.code = ?", languageCode).
		Limit(1).
		Scan(&out)
	if tx.Error != nil {
		return translationRecord{}, tx.Error
	}
	if tx.RowsAffected == 0 {
		return translationRecord{}, ErrNotFound
	}
	return out, nil
}

func wrapNotFound(entity string, key string, err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("%s not found for key %q: %w", entity, key, ErrNotFound)
	}
	return err
}
