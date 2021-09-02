package usecase

import (
	"toychain/pkg/domain"

	badger "github.com/dgraph-io/badger/v3"
)

type TransactionUsecase struct {
	db          *badger.DB
	accountRepo domain.AccountRepository
}

func NewTransactionUsecase(db *badger.DB, accountRepo domain.AccountRepository) *TransactionUsecase {
	return &TransactionUsecase{
		db:          db,
		accountRepo: accountRepo,
	}
}
