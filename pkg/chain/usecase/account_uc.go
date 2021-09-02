package usecase

import (
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"toychain/internal/pkg/crypto/base58"
	"toychain/pkg/domain"

	badger "github.com/dgraph-io/badger/v3"
)

type AccountUsecase struct {
	db          *badger.DB
	accountRepo domain.AccountRepository
}

func NewAccountUsecase(db *badger.DB, accountRepo domain.AccountRepository) *AccountUsecase {
	return &AccountUsecase{
		db:          db,
		accountRepo: accountRepo,
	}
}

func (uc *AccountUsecase) GenerateAccount(ctx context.Context) (*domain.Account, error) {
	pubKey, privKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		return nil, err
	}

	account := domain.Account{
		Address:    base58.EncodeCheck(pubKey),
		PublicKey:  hex.EncodeToString(pubKey),
		PrivateKey: hex.EncodeToString(privKey),
	}

	return &account, nil
}

func (uc *AccountUsecase) Account(ctx context.Context, address string) (*domain.Account, error) {
	var account *domain.Account
	var err error

	err = uc.db.View(func(txn *badger.Txn) error {

		accountKey := []byte(address)
		account, err = uc.accountRepo.AccountTX(ctx, accountKey, txn)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return account, nil
	}

	return account, nil
}
