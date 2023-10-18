package model

import (
	"errors"
	"time"

	"github.com/asaskevich/govalidator"
	uuid "github.com/satori/go.uuid"
)

const (
	TransactionPending   string = "PENDING"
	TransactionCompleted string = "COMPLETED"
	TransactionError     string = "ERROR"
	TransactionConfirmed string = "CONFIRMED"
)

type ITransactionRepository interface {
	Register(transaction *Transaction) error
	Save(transaction *Transaction) error
	Find(id string) (*Transaction, error)
}

type Transactions struct {
	Transaction []Transaction
}

type Transaction struct {
	Base              `valid:"required"`
	AccountFrom       *Account `valid:"-"`
	AccountIDFrom     *string  `gorm:"column:account_from_id;type:uuid;not null" valid:"notnull"`
	PixKeyTo          *PixKey  `valid:"-"`
	PixKeyIDTo        *string  `gorm:"column:pix_key_id_to;type:uuid;not null" valid:"notnull"`
	Value             int      `json:"value" gorm:"type:int" valid:"notnull"`
	Status            string   `json:"status" gorm:"type:varchar(20)" valid:"notnull"`
	Description       string   `json:"description" gorm:"type:varchar(255)" valid:"-"`
	CancelDescription string   `json:"cancel_description" gorm:"type:varchar(255)" valid:"-"`
}

func (transaction *Transaction) validate() error {
	_, err := govalidator.ValidateStruct(transaction)

	if transaction.Value <= 0 {
		return errors.New("invalid transaction value")
	}

	if transaction.Status != TransactionPending &&
		transaction.Status != TransactionCompleted &&
		transaction.Status != TransactionError &&
		transaction.Status != TransactionConfirmed {
		return errors.New("invalid transaction status")
	}

	if transaction.PixKeyTo.AccountID == transaction.AccountFrom.ID {
		return errors.New("the transaction cannot be to the same account")
	}

	return err
}

func NewTransaction(accountFrom *Account, pixKeyTo *PixKey, value int, description string) (*Transaction, error) {
	transaction := Transaction{
		AccountFrom: accountFrom,
		PixKeyTo:    pixKeyTo,
		Value:       value,
		Description: description,
		Status:      TransactionPending,
	}

	transaction.ID = uuid.NewV4().String()
	transaction.CreatedAt = time.Now()

	err := transaction.validate()
	if err != nil {
		return nil, err
	}

	return &transaction, nil
}

func (t *Transaction) Complete() error {
	t.Status = TransactionCompleted
	t.UpdatedAt = time.Now()
	return t.validate()
}

func (t *Transaction) Cancel(description string) error {
	t.Status = TransactionError
	t.CancelDescription = description
	t.UpdatedAt = time.Now()
	return t.validate()
}

func (t *Transaction) Confirm() error {
	t.Status = TransactionConfirmed
	t.UpdatedAt = time.Now()
	return t.validate()
}
