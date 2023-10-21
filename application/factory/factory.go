package factory

import (
	"github.com/Lucas-Mondini/code-pix/application/usecase"
	"github.com/Lucas-Mondini/code-pix/infrastructure/repository"
	"github.com/jinzhu/gorm"
)

func TransactionUseCaseFactory(database *gorm.DB) usecase.TransactionUseCase {
	pixRepository := repository.PixKeyRepositoryDB{DB: database}
	transactionRepository := repository.TransactionRepositoryDB{DB: database}

	return usecase.TransactionUseCase{
		TransactionRepository: &transactionRepository,
		PixKeyRepository:      &pixRepository,
	}
}
