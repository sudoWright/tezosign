package contract

import (
	"database/sql"
	"errors"
	"tezosign/models"
	"tezosign/types"

	"gorm.io/gorm"
)

//go:generate mockgen -source ./contract.go -destination ./mock_contract/main.go Repo
type (
	// Repository is the account repo implementation.
	Repository struct {
		db *gorm.DB
	}

	Repo interface {
		GetOrCreateContract(address types.Address) (contract models.Contract, err error)
		UpdateContractLastOperationBlock(contractID, blockLevel uint64) (err error)
		GetContractByID(id uint64) (contract models.Contract, err error)
		GetContract(address types.Address) (contract models.Contract, isFound bool, err error)
		GetContractsList(limit, offset int) (contracts []models.Contract, err error)
		GetMaxContractPendingNone(contractID uint64) (sql.NullInt64, error)
		SavePayload(request models.Request) error
		UpdatePayload(request models.Request) error
		GetPayloadByContractAndCounter(contractID uint64, counter int64) (models.Request, bool, error)
		GetPayloadByHash(id string) (models.Request, bool, error)
		GetPayloadsReportByContractID(id uint64, isOwner bool, limit, offset int) ([]models.RequestReport, error)
		GetSignaturesByPayloadID(id uint64, signatureType models.PayloadType) ([]models.Signature, error)
		SavePayloadSignature(signature models.Signature) error
		GetPayloadSignature(sig types.Signature) (signature models.Signature, isFound bool, err error)
		GetSignaturesCount(id uint64) (count int64, err error)
	}
)

// New creates an instance of repository using the provided db.
func New(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetPayloadSignature(sig types.Signature) (signature models.Signature, isFound bool, err error) {
	err = r.db.Model(models.Signature{}).
		Where("sig_data = ?", sig).
		First(&signature).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return signature, false, nil
		}
		return signature, false, err
	}

	return signature, true, nil
}

func (r *Repository) SavePayloadSignature(sign models.Signature) (err error) {
	err = r.db.
		Model(models.Signature{}).
		Create(&sign).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetSignaturesCount(id uint64) (count int64, err error) {
	err = r.db.Model(models.Signature{}).
		Where("req_id = ?", id).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *Repository) GetSignaturesByPayloadID(id uint64, signatureType models.PayloadType) (signatures []models.Signature, err error) {
	err = r.db.Model(models.Signature{}).
		Where("req_id = ? and sig_type = ?", id, signatureType).
		Find(&signatures).Error
	if err != nil {
		return signatures, err
	}
	return signatures, nil
}

func (r *Repository) GetOrCreateContract(address types.Address) (contract models.Contract, err error) {
	err = r.db.Model(models.Contract{}).
		FirstOrCreate(&contract, models.Contract{Address: address}).
		Error
	if err != nil {
		return contract, err
	}

	return contract, nil
}

func (r *Repository) UpdateContractLastOperationBlock(contractID, blockLevel uint64) (err error) {
	err = r.db.Model(&models.Contract{ID: contractID}).
		Updates(models.Contract{
			LastOperationBlockLevel: blockLevel,
		}).
		Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetContractByID(id uint64) (contract models.Contract, err error) {
	err = r.db.Model(models.Contract{}).
		Where("ctr_id = ?", id).
		First(&contract).Error
	if err != nil {
		return contract, err
	}
	return contract, nil
}

func (r *Repository) GetContract(address types.Address) (contract models.Contract, isFound bool, err error) {
	err = r.db.Model(models.Contract{}).
		Where("ctr_address = ?", address).
		First(&contract).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return contract, false, nil
		}
		return contract, false, err
	}
	return contract, true, nil
}

func (r *Repository) GetContractsList(limit, offset int) (contracts []models.Contract, err error) {
	db := r.db.Model(models.Contract{})

	if limit > 0 {
		db = db.Limit(limit)
	}

	err = db.Offset(offset).
		Find(&contracts).Error
	if err != nil {
		return contracts, err
	}
	return contracts, nil
}

func (r *Repository) GetMaxContractPendingNone(contractID uint64) (max sql.NullInt64, err error) {
	m := struct {
		M sql.NullInt64
	}{}

	err = r.db.Select("max(req_counter) as m").Table("requests").
		Where("ctr_id = ?", contractID).
		Take(&m).Error
	if err != nil {
		return max, err
	}

	return m.M, nil
}
