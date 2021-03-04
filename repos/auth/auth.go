package auth

import (
	"errors"
	"tezosign/models"
	"tezosign/types"

	"gorm.io/gorm"
)

//go:generate mockgen -source ./auth.go -destination ./mock_auth/main.go Repo
type (
	// Repository is the account repo implementation.
	Repository struct {
		db *gorm.DB
	}

	Repo interface {
		CreateAuthToken(authToken models.AuthToken) (err error)
		GetAuthToken(data string) (authToken models.AuthToken, isFound bool, err error)
		GetActiveTokenByPubKeyAndType(address types.PubKey, tokenType models.TokenType) (authToken models.AuthToken, isFound bool, err error)
		MarkAsUsedAuthToken(id uint64) (err error)
	}
)

// New creates an instance of repository using the provided db.
func New(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetAuthToken(data string) (authToken models.AuthToken, isFound bool, err error) {
	err = r.db.Model(models.AuthToken{}).
		Where("atn_data = ?", data).
		First(&authToken).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return authToken, false, nil
		}
		return authToken, false, err
	}

	return authToken, true, nil
}

func (r *Repository) GetActiveTokenByPubKeyAndType(address types.PubKey, tokenType models.TokenType) (authToken models.AuthToken, isFound bool, err error) {
	err = r.db.Model(models.AuthToken{}).
		Where("atn_pubkey = ? and atn_type = ? and atn_is_used = false and atn_expires_at < now()", address, tokenType).
		First(&authToken).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return authToken, false, nil
		}
		return authToken, false, err
	}

	return authToken, true, nil
}

func (r *Repository) CreateAuthToken(authToken models.AuthToken) (err error) {
	err = r.db.
		Model(models.AuthToken{}).
		Create(&authToken).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) MarkAsUsedAuthToken(id uint64) (err error) {
	err = r.db.
		Model(models.AuthToken{}).
		Where("atn_id = ?", id).
		Update("atn_is_used", true).Error
	if err != nil {
		return err
	}

	return nil
}
