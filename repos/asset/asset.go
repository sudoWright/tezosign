package asset

import (
	"errors"
	"tezosign/models"
	"tezosign/types"

	"gorm.io/gorm"
)

//go:generate mockgen -source ./asset.go -destination ./mock_asset/main.go Repo
type (
	// Repository is the account repo implementation.
	Repository struct {
		db *gorm.DB
	}

	Repo interface {
		GetAssetsList(contract uint64, isOwner bool) (assets []models.Asset, err error)
		GetAsset(contract uint64, assetAddress types.Address) (assets models.Asset, isFound bool, err error)
		CreateAsset(asset models.Asset) (err error)
		UpdateAsset(asset models.Asset) (err error)
		DeleteContractAsset(assetID uint64) (err error)
	}
)

// New creates an instance of repository using the provided db.
func New(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) CreateAsset(asset models.Asset) (err error) {
	err = r.db.
		Model(models.Asset{}).
		Create(&asset).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) UpdateAsset(asset models.Asset) (err error) {
	err = r.db.Model(&models.Asset{ID: asset.ID}).
		Updates(models.Asset{
			Name:                    asset.Name,
			Scale:                   asset.Scale,
			Ticker:                  asset.Ticker,
			LastOperationBlockLevel: asset.LastOperationBlockLevel,
		}).
		Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetAssetsList(contractID uint64, isOwner bool) (assets []models.Asset, err error) {

	db := r.db.Model(models.Asset{}).
		Where("ctr_id IS NULL")

	if isOwner {
		db = db.Or("ctr_id = ?", contractID)
	}

	err = db.Find(&assets).Error
	if err != nil {
		return assets, err
	}
	return assets, nil
}

func (r *Repository) DeleteContractAsset(assetID uint64) (err error) {
	err = r.db.
		Model(models.Asset{}).
		Delete(&models.Asset{ID: assetID}).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetAsset(contract uint64, assetAddress types.Address) (asset models.Asset, isFound bool, err error) {
	err = r.db.Model(models.Asset{}).
		Where("(ctr_id = ?  OR ctr_id IS NULL) AND ast_address = ? ", contract, assetAddress).
		First(&asset).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return asset, false, nil
		}
		return asset, false, err
	}

	return asset, true, nil
}
