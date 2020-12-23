package contract

import (
	"gorm.io/gorm"
)

//go:generate mockgen -source ./contract.go -destination ./mock_contract/main.go Repo
type (
	// Repository is the account repo implementation.
	Repository struct {
		db *gorm.DB
	}

	Repo interface {
	}
)

const ()

// New creates an instance of repository using the provided db.
func New(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}
