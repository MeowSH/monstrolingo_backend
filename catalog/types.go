package catalog

import "time"

type CategoryKey string

const (
	CategoryItems       CategoryKey = "items"
	CategoryWeapons     CategoryKey = "weapons"
	CategoryArmor       CategoryKey = "armor"
	CategorySkills      CategoryKey = "skills"
	CategoryDecorations CategoryKey = "decorations"
	CategoryCharms      CategoryKey = "charms"
	CategoryFoodSkills  CategoryKey = "food-skills"
	CategoryKinsects    CategoryKey = "kinsects"
)

type CategoryTableRequest struct {
	SourceLang string `query:"source_lang"`
	TargetLang string `query:"target_lang"`
	Page       int    `query:"page"`
	Limit      int    `query:"limit"`
}

type CategoryDetailRequest struct {
	ExternalKey string `path:"external_key"`
	TargetLang  string `query:"target_lang"`
}

type TargetLanguageRequest struct {
	TargetLang string `query:"target_lang"`
}

type normalizedTableQuery struct {
	SourceLang string
	TargetLang string
	Page       int
	Limit      int
	Offset     int
}

type normalizedDetailQuery struct {
	ExternalKey string
	TargetLang  string
}

type TableTranslation struct {
	Language    string `json:"language"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CategoryTableRow struct {
	ExternalKey string           `json:"external_key"`
	Source      TableTranslation `json:"source"`
	Target      TableTranslation `json:"target"`
}

type Pagination struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
	HasNext    bool  `json:"has_next"`
}

type CategoryTableResponse struct {
	Items      []CategoryTableRow `json:"items"`
	Pagination Pagination         `json:"pagination"`
}

type DetailTranslation struct {
	LanguageCode  string `json:"language_code"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	FlavorText    string `json:"flavor_text,omitempty"`
	EffectSummary string `json:"effect_summary,omitempty"`
	Slug          string `json:"slug"`
}

type SkillLinkDetail struct {
	SkillExternalKey string `json:"skill_external_key"`
	Name             string `json:"name"`
	Description      string `json:"description"`
	EffectSummary    string `json:"effect_summary,omitempty"`
	SkillKind        string `json:"skill_kind,omitempty"`
	IsSetBonusSkill  bool   `json:"is_set_bonus_skill"`
	Level            int16  `json:"level"`
	SortOrder        int16  `json:"sort_order"`
	SourceType       string `json:"source_type,omitempty"`
	RankRequired     int16  `json:"rank_required,omitempty"`
}

type AssociatedJewelDetail struct {
	DecorationExternalKey string `json:"decoration_external_key"`
	Name                  string `json:"name"`
	SlotSize              int16  `json:"slot_size"`
	Rarity                int16  `json:"rarity"`
	SkillLevel            int16  `json:"skill_level"`
}

type SkillLevelDetail struct {
	Level           int16  `json:"level"`
	EffectValueText string `json:"effect_value_text"`
	Description     string `json:"description"`
}

type FoodSkillLevelDetail struct {
	Level             int16  `json:"level"`
	DurationSeconds   *int32 `json:"duration_seconds"`
	ActivationPercent *int16 `json:"activation_percent"`
	EffectValueText   string `json:"effect_value_text"`
}

type ItemDetailResponse struct {
	ExternalKey        string            `json:"external_key"`
	SyncedAt           time.Time         `json:"synced_at"`
	Category           string            `json:"category"`
	Subcategory        string            `json:"subcategory"`
	Rarity             int16             `json:"rarity"`
	CarryLimit         *int16            `json:"carry_limit"`
	BuyPrice           *int32            `json:"buy_price"`
	SellPrice          *int32            `json:"sell_price"`
	Points             *int32            `json:"points"`
	IsCraftingMaterial bool              `json:"is_crafting_material"`
	IsConsumable       bool              `json:"is_consumable"`
	IconKey            string            `json:"icon_key"`
	SortOrder          int32             `json:"sort_order"`
	Translation        DetailTranslation `json:"translation"`
}

type WeaponDetailResponse struct {
	ExternalKey     string            `json:"external_key"`
	SyncedAt        time.Time         `json:"synced_at"`
	WeaponType      string            `json:"weapon_type"`
	Rarity          int16             `json:"rarity"`
	Attack          int32             `json:"attack"`
	AffinityPercent int16             `json:"affinity_percent"`
	DefenseBonus    int16             `json:"defense_bonus"`
	Slot1Level      int16             `json:"slot1_level"`
	Slot2Level      int16             `json:"slot2_level"`
	Slot3Level      int16             `json:"slot3_level"`
	ElementType     string            `json:"element_type"`
	ElementValue    int32             `json:"element_value"`
	AilmentType     string            `json:"ailment_type"`
	AilmentValue    int32             `json:"ailment_value"`
	SharpnessRed    int16             `json:"sharpness_red"`
	SharpnessOrange int16             `json:"sharpness_orange"`
	SharpnessYellow int16             `json:"sharpness_yellow"`
	SharpnessGreen  int16             `json:"sharpness_green"`
	SharpnessBlue   int16             `json:"sharpness_blue"`
	SharpnessWhite  int16             `json:"sharpness_white"`
	SharpnessPurple int16             `json:"sharpness_purple"`
	CraftCost       *int32            `json:"craft_cost"`
	UpgradeCost     *int32            `json:"upgrade_cost"`
	TreeDepth       int16             `json:"tree_depth"`
	IsFinalUpgrade  bool              `json:"is_final_upgrade"`
	SortOrder       int32             `json:"sort_order"`
	Translation     DetailTranslation `json:"translation"`
	Skills          []SkillLinkDetail `json:"skills"`
	SetSkills       []SkillLinkDetail `json:"set_skills,omitempty"`
	RegularSkills   []SkillLinkDetail `json:"regular_skills,omitempty"`
}

type ArmorDetailResponse struct {
	ExternalKey         string            `json:"external_key"`
	SyncedAt            time.Time         `json:"synced_at"`
	ArmorSetKey         string            `json:"armor_set_key"`
	ArmorSetName        string            `json:"armor_set_name"`
	ArmorSetVariant     string            `json:"armor_set_variant"`
	PieceType           string            `json:"piece_type"`
	Rank                string            `json:"rank"`
	Gender              string            `json:"gender"`
	Rarity              int16             `json:"rarity"`
	DefenseBase         int16             `json:"defense_base"`
	DefenseMax          int16             `json:"defense_max"`
	DefenseAugmentedMax int16             `json:"defense_augmented_max"`
	FireRes             int16             `json:"fire_res"`
	WaterRes            int16             `json:"water_res"`
	ThunderRes          int16             `json:"thunder_res"`
	IceRes              int16             `json:"ice_res"`
	DragonRes           int16             `json:"dragon_res"`
	Slot1Level          int16             `json:"slot1_level"`
	Slot2Level          int16             `json:"slot2_level"`
	Slot3Level          int16             `json:"slot3_level"`
	IsLayered           bool              `json:"is_layered"`
	SortOrder           int32             `json:"sort_order"`
	Translation         DetailTranslation `json:"translation"`
	Skills              []SkillLinkDetail `json:"skills"`
	SetSkills           []SkillLinkDetail `json:"set_skills,omitempty"`
	RegularSkills       []SkillLinkDetail `json:"regular_skills,omitempty"`
}

type SkillDetailResponse struct {
	ExternalKey      string                  `json:"external_key"`
	SyncedAt         time.Time               `json:"synced_at"`
	SkillKind        string                  `json:"skill_kind"`
	MaxLevel         int16                   `json:"max_level"`
	IsBinary         bool                    `json:"is_binary"`
	IsSetBonusSkill  bool                    `json:"is_set_bonus_skill"`
	SortOrder        int32                   `json:"sort_order"`
	Translation      DetailTranslation       `json:"translation"`
	Levels           []SkillLevelDetail      `json:"levels"`
	AssociatedJewels []AssociatedJewelDetail `json:"associated_jewels,omitempty"`
}

type DecorationDetailResponse struct {
	ExternalKey string            `json:"external_key"`
	SyncedAt    time.Time         `json:"synced_at"`
	SlotSize    int16             `json:"slot_size"`
	Rarity      int16             `json:"rarity"`
	IsCraftable bool              `json:"is_craftable"`
	CraftCost   *int32            `json:"craft_cost"`
	SortOrder   int32             `json:"sort_order"`
	Translation DetailTranslation `json:"translation"`
	Skills      []SkillLinkDetail `json:"skills"`
}

type CharmDetailResponse struct {
	ExternalKey string            `json:"external_key"`
	SyncedAt    time.Time         `json:"synced_at"`
	Rarity      int16             `json:"rarity"`
	MaxRank     int16             `json:"max_rank"`
	CraftCost   *int32            `json:"craft_cost"`
	SortOrder   int32             `json:"sort_order"`
	Translation DetailTranslation `json:"translation"`
	Skills      []SkillLinkDetail `json:"skills"`
}

type FoodSkillDetailResponse struct {
	ExternalKey           string                 `json:"external_key"`
	SyncedAt              time.Time              `json:"synced_at"`
	FoodCategory          string                 `json:"food_category"`
	MaxLevel              int16                  `json:"max_level"`
	BaseDurationSeconds   *int32                 `json:"base_duration_seconds"`
	BaseActivationPercent *int16                 `json:"base_activation_percent"`
	SortOrder             int32                  `json:"sort_order"`
	Translation           DetailTranslation      `json:"translation"`
	Levels                []FoodSkillLevelDetail `json:"levels"`
}

type KinsectDetailResponse struct {
	ExternalKey           string            `json:"external_key"`
	SyncedAt              time.Time         `json:"synced_at"`
	KinsectType           string            `json:"kinsect_type"`
	AttackType            string            `json:"attack_type"`
	PowderType            string            `json:"powder_type"`
	KinsectBonusPrimary   string            `json:"kinsect_bonus_primary"`
	KinsectBonusSecondary string            `json:"kinsect_bonus_secondary"`
	Rarity                int16             `json:"rarity"`
	PowerValue            int16             `json:"power_value"`
	SpeedValue            int16             `json:"speed_value"`
	HealValue             int16             `json:"heal_value"`
	StaminaValue          int16             `json:"stamina_value"`
	ElementType           string            `json:"element_type"`
	ElementValue          int32             `json:"element_value"`
	SortOrder             int32             `json:"sort_order"`
	Translation           DetailTranslation `json:"translation"`
}

type LanguageOption struct {
	Code  string `json:"code"`
	Label string `json:"label"`
}

type LanguagesResponse struct {
	Languages []LanguageOption `json:"languages"`
}
