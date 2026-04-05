package catalog

type categorySpec struct {
	Key              CategoryKey
	CanonicalTable   string
	TranslationTable string
	TranslationFK    string
	SortColumn       string
	HasFlavorText    bool
	HasEffectSummary bool
}

var categorySpecs = map[CategoryKey]categorySpec{
	CategoryItems: {
		Key:              CategoryItems,
		CanonicalTable:   "items",
		TranslationTable: "item_translations",
		TranslationFK:    "item_id",
		SortColumn:       "sort_order",
		HasFlavorText:    true,
	},
	CategoryWeapons: {
		Key:              CategoryWeapons,
		CanonicalTable:   "weapons",
		TranslationTable: "weapon_translations",
		TranslationFK:    "weapon_id",
		SortColumn:       "sort_order",
	},
	CategoryArmor: {
		Key:              CategoryArmor,
		CanonicalTable:   "armor_pieces",
		TranslationTable: "armor_piece_translations",
		TranslationFK:    "armor_piece_id",
		SortColumn:       "sort_order",
	},
	CategorySkills: {
		Key:              CategorySkills,
		CanonicalTable:   "skills",
		TranslationTable: "skill_translations",
		TranslationFK:    "skill_id",
		SortColumn:       "sort_order",
		HasEffectSummary: true,
	},
	CategoryDecorations: {
		Key:              CategoryDecorations,
		CanonicalTable:   "decorations",
		TranslationTable: "decorations_translations",
		TranslationFK:    "decoration_id",
		SortColumn:       "sort_order",
	},
	CategoryCharms: {
		Key:              CategoryCharms,
		CanonicalTable:   "charms",
		TranslationTable: "charm_translations",
		TranslationFK:    "charm_id",
		SortColumn:       "sort_order",
	},
	CategoryFoodSkills: {
		Key:              CategoryFoodSkills,
		CanonicalTable:   "food_skills",
		TranslationTable: "food_skill_translations",
		TranslationFK:    "food_skill_id",
		SortColumn:       "sort_order",
	},
	CategoryKinsects: {
		Key:              CategoryKinsects,
		CanonicalTable:   "kinsects",
		TranslationTable: "kinsect_translations",
		TranslationFK:    "kinsect_id",
		SortColumn:       "sort_order",
	},
}

func getCategorySpec(key CategoryKey) (categorySpec, error) {
	spec, ok := categorySpecs[key]
	if !ok {
		return categorySpec{}, invalidArgumentf("unsupported category: %s", key)
	}
	return spec, nil
}
