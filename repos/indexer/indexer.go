package indexer

import (
	"errors"
	"fmt"
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
		GetContractOperations(contract types.Address, blockLevel uint64, entrypoint string) ([]models.TransactionOperation, error)
		GetContractRevealOperation(contract types.Address) (models.RevealOperation, bool, error)
		GetContractOriginationOperation(txID string) (tx models.OriginationOperation, isFound bool, err error)

		GetContractStorage(address types.Address) (storage models.Storage, isFound bool, err error)
		GetContractStorageChange(address types.Address, level uint64) (storage []models.Storage, err error)
		GetContractsStoragesContainsKey(contracts []string, key string) ([]string, error)
		GetContractScript(address types.Address) (script models.Script, isFound bool, err error)
		GetAccount(address types.Address) (account models.Account, isFound bool, err error)
		GetAccountByID(id uint64) (account models.Account, isFound bool, err error)

		GetLastBlock() (block models.Block, err error)
		GetTezosQuote() (models.Quote, error)
	}
)

// New creates an instance of repository using the provided db.
func New(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetContractOperations(contract types.Address, blockLevel uint64, entrypoint string) (operations []models.TransactionOperation, err error) {

	db := r.db.Table("TransactionOps").
		Joins(`LEFT JOIN "Accounts" a on "TargetId" = a."Id"`).
		Where(`"Address" = ?`, contract.String()).
		Where(`"Level" > ?`, blockLevel)

	if len(entrypoint) > 0 {
		db = db.Where(`"Entrypoint" = ?`, entrypoint)
	}

	err = db.Order(`"TransactionOps"."Id" asc`).
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

func (r *Repository) GetContractOriginationOperation(txID string) (tx models.OriginationOperation, isFound bool, err error) {

	err = r.db.Select("*").
		Table("OriginationOps").
		Joins(`LEFT JOIN "Accounts" a on "ContractId" = a."Id"`).
		Where(`"OpHash" = ?`, txID).
		First(&tx).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return tx, false, nil
		}
		return tx, false, err
	}

	return tx, true, nil
}

func (r *Repository) GetContractStorage(address types.Address) (storage models.Storage, isFound bool, err error) {
	err = r.db.Select("*").
		Table("Storages").
		Joins(`LEFT JOIN "Accounts" a on "ContractId" = a."Id"`).
		Where(`"Address" = ? AND "Current" IS TRUE`, address.String()).
		First(&storage).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return storage, false, nil
		}
		return storage, false, err
	}

	return storage, true, nil
}

func (r *Repository) GetContractStorageChange(address types.Address, level uint64) (storages []models.Storage, err error) {
	err = r.db.Select("*").
		Table("Storages").
		Joins(`LEFT JOIN "Accounts" a on "ContractId" = a."Id"`).
		Where(`"Address" = ? AND "Level" <= ?`, address.String(), level).
		Order(`"Level" desc`).
		Limit(2).
		Find(&storages).Error
	if err != nil {
		return storages, err
	}

	return storages, nil
}

func (r *Repository) GetContractsStoragesContainsKey(contracts []string, key string) (resp []string, err error) {

	subQuery := r.db.
		Select(fmt.Sprintf(`"Accounts"."Address",("JsonValue" -> 'keys') ?| array['%s'] AS is_contain`, key)).
		Table("Storages").
		Joins(`LEFT JOIN "Accounts" on "Accounts"."Id" = "ContractId"`).
		Where(`"Current" IS TRUE AND "Address" IN (?)`, contracts)

	err = r.db.Select(`"Address" `).
		Table("(?) s", subQuery).
		Where("s.is_contain IS TRUE").
		Find(&resp).Error
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (r *Repository) GetContractScript(address types.Address) (script models.Script, isFound bool, err error) {
	err = r.db.Select("*").
		Table("Scripts").
		Joins(`LEFT JOIN "Accounts" a on "ContractId" = a."Id"`).
		Where(`"Address" = ? AND "Current" IS TRUE`, address.String()).
		First(&script).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return script, false, nil
		}
		return script, false, err
	}

	return script, true, nil
}

func (r *Repository) GetAccount(address types.Address) (account models.Account, isFound bool, err error) {
	err = r.db.
		Table("Accounts").
		Where(`"Address" = ?`, address.String()).
		First(&account).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return account, false, nil
		}
		return account, false, err
	}

	return account, true, nil
}

func (r *Repository) GetAccountByID(id uint64) (account models.Account, isFound bool, err error) {
	err = r.db.
		Table("Accounts").
		Where(`"Id" = ?`, id).
		First(&account).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return account, false, nil
		}
		return account, false, err
	}

	return account, true, nil
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

func (r *Repository) GetLastBlock() (block models.Block, err error) {
	err = r.db.Table("Blocks").
		Order(`"Blocks"."Id" desc`).
		First(&block).Error
	if err != nil {
		return block, err
	}

	return block, nil
}
