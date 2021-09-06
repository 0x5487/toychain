package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"toychain/pkg/domain"

	badger "github.com/dgraph-io/badger/v3"
)

type AccountRepo struct {
	db *badger.DB
}

func NewAccountRepo(db *badger.DB) *AccountRepo {
	return &AccountRepo{
		db: db,
	}
}

func (repo *AccountRepo) StoreTX(ctx context.Context, account *domain.Account, txn *badger.Txn) error {
	val, err := account.Serialize()
	if err != nil {
		return err
	}

	return txn.Set(account.Key(), val)
}

func (repo *AccountRepo) AccountTX(ctx context.Context, key []byte, txn *badger.Txn) (*domain.Account, error) {
	account := domain.Account{}

	item, err := txn.Get(key)
	if err != nil {
		if errors.Is(err, badger.ErrKeyNotFound) {
			return nil, fmt.Errorf("repo: can't find account key. err: %w", domain.ErrNotFound)
		}
		return nil, err
	}

	err = item.Value(func(val []byte) error {
		decoder := json.NewDecoder(bytes.NewReader(val))
		return decoder.Decode(&account)
	})

	if err != nil {
		return nil, err
	}

	return &account, nil
}
