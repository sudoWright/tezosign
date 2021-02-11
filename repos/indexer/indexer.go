package indexer

import (
	"tezosign/models"
	"tezosign/types"

	"gorm.io/gorm"
)

//go:generate mockgen -source ./indexer.go -destination ./mock_indexer/main.go Repo
type (
	// Repository is the account repo implementation.
	Repository struct {
		db *gorm.DB
	}

	Repo interface {
		GetContractOperations(contract types.Address, blockLevel uint64) ([]models.TezosOperation, error)
	}
)

// New creates an instance of repository using the provided db.
func New(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetContractOperations(contract types.Address, blockLevel uint64) (operations []models.TezosOperation, err error) {

	err = r.db.Table("TransactionOps").
		Joins(`LEFT JOIN "Accounts" a on "TargetId" = a."Id"`).
		Where(`"Address" = ?`, contract.String()).
		Where(`"Level" > ?`, blockLevel).
		Order(`"Id" asc`).
		Find(&operations).Error
	if err != nil {
		return operations, err
	}

	return operations, nil
}
