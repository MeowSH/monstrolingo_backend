package catalogcore

import "monstrolingo_backend/catalog"

type CategoryKey = catalog.CategoryKey

const (
	CategoryItems       CategoryKey = catalog.CategoryItems
	CategoryWeapons     CategoryKey = catalog.CategoryWeapons
	CategoryArmor       CategoryKey = catalog.CategoryArmor
	CategorySkills      CategoryKey = catalog.CategorySkills
	CategoryDecorations CategoryKey = catalog.CategoryDecorations
	CategoryCharms      CategoryKey = catalog.CategoryCharms
	CategoryFoodSkills  CategoryKey = catalog.CategoryFoodSkills
	CategoryKinsects    CategoryKey = catalog.CategoryKinsects
)

type CategoryTableRequest = catalog.CategoryTableRequest
type CategoryDetailRequest = catalog.CategoryDetailRequest
type TargetLanguageRequest = catalog.TargetLanguageRequest

type TableTranslation = catalog.TableTranslation
type CategoryTableRow = catalog.CategoryTableRow
type Pagination = catalog.Pagination
type CategoryTableResponse = catalog.CategoryTableResponse

type DetailTranslation = catalog.DetailTranslation
type SkillLinkDetail = catalog.SkillLinkDetail
type SkillLevelDetail = catalog.SkillLevelDetail
type FoodSkillLevelDetail = catalog.FoodSkillLevelDetail

type ItemDetailResponse = catalog.ItemDetailResponse
type WeaponDetailResponse = catalog.WeaponDetailResponse
type ArmorDetailResponse = catalog.ArmorDetailResponse
type SkillDetailResponse = catalog.SkillDetailResponse
type DecorationDetailResponse = catalog.DecorationDetailResponse
type CharmDetailResponse = catalog.CharmDetailResponse
type FoodSkillDetailResponse = catalog.FoodSkillDetailResponse
type KinsectDetailResponse = catalog.KinsectDetailResponse

type LanguageOption = catalog.LanguageOption
type LanguagesResponse = catalog.LanguagesResponse
