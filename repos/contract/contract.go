package contract

import (
	"gorm.io/gorm"
	"msig/models"
	"msig/types"
)

//go:generate mockgen -source ./contract.go -destination ./mock_contract/main.go Repo
type (
	// Repository is the account repo implementation.
	Repository struct {
		db *gorm.DB
	}

	Repo interface {
		GetOrCreateContract(address types.Address) (contract models.Contract, err error)
		GetContractByID(id uint64) (contract models.Contract, err error)
		SavePayload(request models.Request) error
		SavePayloadSignature(sign models.Sig) error
		GetSignaturesCount(id uint64) (count uint64, err error)
		GetPayload(id string) (models.Request, bool, error)
		GetSignatures(id string) ([]models.Sig, error)
	}
)

const ()

// New creates an instance of repository using the provided db.
func New(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) SavePayloadSignature(sign models.Sig) error {
	return nil
}

func (r *Repository) GetSignaturesCount(id uint64) (count uint64, err error) {
	return 0, nil
}

func (r *Repository) GetPayload(id string) (models.Request, bool, error) {
	return models.Request{}, false, nil
}

func (r *Repository) SavePayload(request models.Request) error {
	return nil
}

func (r *Repository) GetSignatures(id string) ([]models.Sig, error) {
	return nil, nil
}

func (r *Repository) GetOrCreateContract(address types.Address) (contract models.Contract, err error) {
	return contract, nil
}

func (r *Repository) GetContractByID(id uint64) (contract models.Contract, err error) {
	return contract, nil
}
