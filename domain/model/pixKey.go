package model

import (
	"errors"
	"time"

	"github.com/asaskevich/govalidator"
	uuid "github.com/satori/go.uuid"
)

type IPixKeyRepository interface {
	RegisterKey(pixKey *PixKey) (*PixKey, error)
	FindKeyByKind(key string, kind string) (*PixKey, error)
	AddBank(bank *Bank) error
	AddAccount(account *Account) error
	FindAccountByID(id string) (*Account, error)
}

type PixKey struct {
	Base      `valid:"required"`
	Account   *Account `valid:"-"`
	Key       string   `json:"key" gorm:"type:varchar(255)" valid:"notnull"`
	Kind      string   `json:"kind" gorm:"type:varchar(20)" valid:"notnull"`
	AccountID string   `gorm:"column:account_id;type:uuid;not null" valid:"-"`
	Status    string   `json:"status" gorm:"type:varchar(8)" valid:"notnull"`
}

func (pixKey *PixKey) validate() error {
	_, err := govalidator.ValidateStruct(pixKey)

	if pixKey.Kind != "email" && pixKey.Kind != "CPF" {
		return errors.New("invalid pix key type")
	}

	if pixKey.Status != "active" && pixKey.Status != "inactive" {
		return errors.New("invalid pix key status")
	}
	return err
}

func NewPixKey(kind string, account *Account, key string) (*PixKey, error) {
	pixKey := PixKey{
		Kind:    kind,
		Account: account,
		Key:     key,
		Status:  "active",
	}
	pixKey.ID = uuid.NewV4().String()
	pixKey.CreatedAt = time.Now()

	err := pixKey.validate()
	if err != nil {
		return nil, err
	}

	return &pixKey, nil
}
