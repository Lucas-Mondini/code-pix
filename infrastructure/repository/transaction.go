package repository

import (
	"fmt"

	"github.com/Lucas-Mondini/code-pix/domain/model"
	"github.com/jinzhu/gorm"
)

// Register(transaction *Transaction) error
// Save(transaction *Transaction) error
// Find(id string) (*Transaction, error)

type TransactionRepositoryDB struct {
	DB *gorm.DB
}

func (r *TransactionRepositoryDB) Register(t *model.Transaction) error {
	return r.DB.Create(t).Error
}

func (r *TransactionRepositoryDB) Save(t *model.Transaction) error {
	return r.DB.Save(t).Error
}

func (r *TransactionRepositoryDB) Find(id string) (*model.Transaction, error) {
	var transaction model.Transaction
	r.DB.Preload("AccountFrom.Bank").First(&transaction, "id = ?", id)

	if transaction.ID == "" {
		return nil, fmt.Errorf("no transaction found")
	}
	return &transaction, nil
}
