package usecase

import (
	"github.com/Lucas-Mondini/code-pix/domain/model"
)

type PixUseCase struct {
	PixKeyRepository model.IPixKeyRepository
}

func (p *PixUseCase) RegisterKey(key string, kind string, accountID string) (*model.PixKey, error) {
	account, err := p.PixKeyRepository.FindAccountByID(accountID)
	if err != nil {
		return nil, err
	}

	pixKey, err := model.NewPixKey(kind, account, key)

	if err != nil {
		return nil, err
	}

	registeredPixKey, err := p.PixKeyRepository.RegisterKey(pixKey)

	if err != nil {
		return nil, err
	}

	return registeredPixKey, nil
}

func (p *PixUseCase) FindKey(key string, kind string) (*model.PixKey, error) {
	return p.PixKeyRepository.FindKeyByKind(key, kind)
}
