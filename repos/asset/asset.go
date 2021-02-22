package asset

import (
	"tezosign/models"

	"gorm.io/gorm"
)

//go:generate mockgen -source ./asset.go -destination ./mock_asset/main.go Repo
type (
	// Repository is the account repo implementation.
	Repository struct {
		db *gorm.DB
	}

	Repo interface {
		GetAssetsList(limit, offset int) (assets []models.Asset, err error)
	}
)

// New creates an instance of repository using the provided db.
func New(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetAssetsList(limit, offset int) (assets []models.Asset, err error) {
	err = r.db.Model(models.Asset{}).
		Limit(limit).
		Offset(offset).
		Find(&assets).Error
	if err != nil {
		return assets, err
	}
	return assets, nil
}
