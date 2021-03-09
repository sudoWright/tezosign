package indexer

import (
	"errors"
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
		GetContractOperations(contract types.Address, blockLevel uint64) ([]models.TransactionOperation, error)
		GetContractRevealOperation(contract types.Address) (models.RevealOperation, bool, error)
		GetTezosQuote() (models.Quote, error)
	}
)

// New creates an instance of repository using the provided db.
func New(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetContractOperations(contract types.Address, blockLevel uint64) (operations []models.TransactionOperation, err error) {

	err = r.db.Table("TransactionOps").
		Joins(`LEFT JOIN "Accounts" a on "TargetId" = a."Id"`).
		Where(`"Address" = ?`, contract.String()).
		Where(`"Level" > ?`, blockLevel).
		Order(`"TransactionOps"."Id" asc`).
		Find(&operations).Error
	if err != nil {
		return operations, err
	}

	return operations, nil
}

func (r *Repository) GetContractRevealOperation(address types.Address) (tx models.RevealOperation, isFound bool, err error) {
	//TODO use single Account table
	err = r.db.Select("*").
		Table("RevealOps").
		Joins(`LEFT JOIN "Accounts" a on "SenderId" = a."Id"`).
		Where(`"Address" = ?`, address.String()).
		First(&tx).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return tx, false, nil
		}
		return tx, false, err
	}

	return tx, true, nil
}

//TODO WIP
type Storage struct {
	Level     uint64 `gorm:"column:Level"`
	Current   bool   `gorm:"column:Current"`
	RawValue  []byte `gorm:"column:RawValue"`
	JsonValue string `gorm:"column:JsonValue"`
}

func (r *Repository) GetContractStorage(address types.Address) (storage Storage, isFound bool, err error) {
	err = r.db.Select("*").
		Table("Storages").
		Joins(`LEFT JOIN "Accounts" a on "ContractId" = a."Id"`).
		Where(`"Address" = ?`, address.String()).
		First(&storage).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return storage, false, nil
		}
		return storage, false, err
	}

	return storage, true, nil
}

type Script struct {
	Current         bool   `gorm:"column:Current"`
	ParameterSchema []byte `gorm:"column:ParameterSchema"`
	StorageSchema   []byte `gorm:"column:StorageSchema"`
	CodeSchema      []byte `gorm:"column:CodeSchema"`
}

func (r *Repository) GetContractScript(address types.Address) (script Script, isFound bool, err error) {
	err = r.db.Select("*").
		Table("Scripts").
		Joins(`LEFT JOIN "Accounts" a on "ContractId" = a."Id"`).
		Where(`"Address" = ?`, address.String()).
		First(&script).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return script, false, nil
		}
		return script, false, err
	}

	return script, true, nil
}

func (r *Repository) GetTezosQuote() (quote models.Quote, err error) {
	err = r.db.Table("Quotes").
		Order(`"Quotes"."Id" desc`).
		First(&quote).Error
	if err != nil {
		return quote, err
	}

	return quote, nil
}
