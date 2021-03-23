package vesting

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
		GetVestingsList(contract uint64) (assets []models.Vesting, err error)
		GetVesting(contract uint64, vestingAddress types.Address) (vesting models.Vesting, isFound bool, err error)
		CreateVesting(asset models.Vesting) (err error)
		UpdateVesting(asset models.Vesting) (err error)
		DeleteContractVesting(vestingID uint64) (err error)
	}
)

// New creates an instance of repository using the provided db.
func New(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) CreateVesting(vesting models.Vesting) (err error) {
	err = r.db.
		Model(models.Vesting{}).
		Create(&vesting).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) UpdateVesting(vesting models.Vesting) (err error) {
	err = r.db.Model(&models.Vesting{ID: vesting.ID}).
		Updates(models.Vesting{
			Name: vesting.Name,
		}).
		Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetVestingsList(contractID uint64) (vesting []models.Vesting, err error) {

	err = r.db.Model(models.Vesting{}).
		Where("ctr_id = ?", contractID).
		Find(&vesting).Error
	if err != nil {
		return vesting, err
	}
	return vesting, nil
}

func (r *Repository) DeleteContractVesting(vestingID uint64) (err error) {
	err = r.db.
		Model(models.Vesting{}).
		Delete(&models.Vesting{ID: vestingID}).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetVesting(contract uint64, vestingAddress types.Address) (vesting models.Vesting, isFound bool, err error) {
	err = r.db.Model(models.Vesting{}).
		Where("ctr_id = ? AND vst_address = ? ", contract, vestingAddress).
		First(&vesting).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return vesting, false, nil
		}
		return vesting, false, err
	}

	return vesting, true, nil
}
