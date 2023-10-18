package repository

import (
	"fmt"

	"github.com/Lucas-Mondini/code-pix/domain/model"
	"github.com/jinzhu/gorm"
)

// RegisterKey(pixKey *PixKey) (*PixKey, error)
// FindKeyByKind(key string, kind string) (*PixKey, error)
// AddBank(bank *Bank) error
// AddAccount(account *Account) error
// FindAccountByID(id string) (*Account, error)

type PixKeyRepositoryDB struct {
	DB *gorm.DB
}

func (r *PixKeyRepositoryDB) AddBank(bank *model.Bank) error {
	err := r.DB.Create(bank).Error
	return err
}
func (r *PixKeyRepositoryDB) AddAccount(account *model.Account) error {
	err := r.DB.Create(account).Error
	return err
}
func (r *PixKeyRepositoryDB) RegisterKey(pk *model.PixKey) (*model.PixKey, error) {
	err := r.DB.Create(pk).Error
	if err != nil {
		return nil, err
	}
	return pk, nil
}

func (r *PixKeyRepositoryDB) FindKeyByKind(key string, kind string) (*model.PixKey, error) {
	var pk model.PixKey
	r.DB.Preload("Account.Bank").First(&pk, "kind = ? and key = ?", kind, key)

	if pk.ID == "" {
		return nil, fmt.Errorf("no key found")
	}
	return &pk, nil
}

func (r *PixKeyRepositoryDB) FindAccountByID(id string) (*model.Account, error) {
	var acc model.Account
	r.DB.Preload("Bank").First(&acc, "id = ?", id)

	if acc.ID == "" {
		return nil, fmt.Errorf("no account found")
	}
	return &acc, nil
}

func (r *PixKeyRepositoryDB) FindBankByID(id string) (*model.Bank, error) {
	var bank model.Bank
	r.DB.First(&bank, "id = ?", id)

	if bank.ID == "" {
		return nil, fmt.Errorf("no bank found")
	}
	return &bank, nil
}
