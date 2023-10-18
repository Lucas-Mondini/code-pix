package usecase

import (
	"errors"
	"log"

	"github.com/Lucas-Mondini/code-pix/domain/model"
)

type TransactionUseCase struct {
	TransactionRepository model.ITransactionRepository
	PixKeyRepository      model.IPixKeyRepository
}

func (t *TransactionUseCase) Register(accountID string, value int, pixKeyTo string, pixKeyKindTo string, description string) (*model.Transaction, error) {
	account, err := t.PixKeyRepository.FindAccountByID(accountID)
	if err != nil {
		return nil, err
	}

	pixKey, err := t.PixKeyRepository.FindKeyByKind(pixKeyTo, pixKeyKindTo)

	if err != nil {
		return nil, err
	}

	transaction, err := model.NewTransaction(account, pixKey, value, description)

	if err != nil {
		return nil, err
	}

	t.TransactionRepository.Register(transaction)

	if transaction.ID == "" {
		return nil, errors.New("cannot process transaction")
	}
	return transaction, nil
}

func (t *TransactionUseCase) Confirm(transactionID string) (*model.Transaction, error) {
	transaction, err := t.TransactionRepository.Find(transactionID)

	if err != nil {
		log.Println("Transaction not found -> ", transactionID)
		return nil, err
	}

	transaction.Status = model.TransactionConfirmed

	err = t.TransactionRepository.Save(transaction)

	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func (t *TransactionUseCase) Complete(transactionID string) (*model.Transaction, error) {
	transaction, err := t.TransactionRepository.Find(transactionID)

	if err != nil {
		log.Println("Transaction not found -> ", transactionID)
		return nil, err
	}

	transaction.Status = model.TransactionCompleted

	err = t.TransactionRepository.Save(transaction)

	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func (t *TransactionUseCase) Error(transactionID string, description string) (*model.Transaction, error) {
	transaction, err := t.TransactionRepository.Find(transactionID)

	if err != nil {
		log.Println("Transaction not found -> ", transactionID)
		return nil, err
	}

	transaction.Status = model.TransactionError
	transaction.CancelDescription = description

	err = t.TransactionRepository.Save(transaction)

	if err != nil {
		return nil, err
	}

	return transaction, nil
}
