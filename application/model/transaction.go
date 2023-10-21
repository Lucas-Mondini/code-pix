package model

import (
	"encoding/json"
	"fmt"

	"github.com/go-playground/validator/v10"
)

type Transaction struct {
	ID           string `json:"id" validate:"required,uuid4"`
	AccountID    string `json:"accountId" validate:"required,uuid4"`
	Value        int    `json:"value" validate:"required,numeric"`
	PixKeyTo     string `json:"pixKeyTo" validate:"required"`
	PixKeyKindTo string `json:"pixKeyKindTo" validate:"required"`
	Description  string `json:"description" validate:"required"`
	Status       string `json:"status" validate:"required"`
	Error        string `json:"description"`
}

func (t *Transaction) validate() error {
	v := validator.New()
	err := v.Struct(t)
	if err != nil {
		fmt.Errorf("Error during Transaction Validation: %s", err.Error())
		return err
	}
	return nil
}

func (t *Transaction) ParseJson(data []byte) error {
	err := json.Unmarshal(data, t)
	if err != nil {
		fmt.Errorf("Error parsing Transaction", err.Error())
	}
	t.validate()
	return err
}

func (t *Transaction) ToJson() ([]byte, error) {
	err := t.validate()
	if err != nil {
		return nil, err
	}
	result, err := json.Marshal(t)
	if err != nil {
		fmt.Errorf("Error parsing Transaction", err.Error())
		return nil, err
	}
	return result, nil
}

func NewTransaction() *Transaction {
	return &Transaction{}
}
