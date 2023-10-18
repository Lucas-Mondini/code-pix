package model

import (
	"time"

	"github.com/asaskevich/govalidator"
	uuid "github.com/satori/go.uuid"
)

type Account struct {
	Base      `valid:"required"`
	Bank      *Bank     `valid:"-"`
	BankID    string    `gorm:"column:bank_id;type:uuid;not null" valid:"-"`
	OwnerName string    `json:"owner_name" gorm:"type:varchar(255);not null" valid:"notnull"`
	Number    string    `json:"number" gorm:"type:varchar(20);not null" valid:"notnull"`
	PixKeys   []*PixKey `gorm:"ForeignKey:AccountID" valid:"-"`
}

func (account *Account) validate() error {
	_, err := govalidator.ValidateStruct(account)
	return err
}

func NewAccount(bank *Bank, ownerName string, number string) (*Account, error) {
	account := Account{
		Bank:      bank,
		OwnerName: ownerName,
		Number:    number,
	}
	account.ID = uuid.NewV4().String()
	account.CreatedAt = time.Now()

	err := account.validate()
	if err != nil {
		return nil, err
	}

	return &account, nil
}
