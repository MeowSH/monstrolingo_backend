package catalogcore

import (
	"monstrolingo_backend/catalog"

	"gorm.io/gorm"
)

type Repository = catalog.Repository

func NewRepository(db *gorm.DB) *Repository {
	return catalog.NewRepository(db)
}
