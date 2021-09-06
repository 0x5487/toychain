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

type ChainRepo struct {
	db *badger.DB
}

func NewChainRepo(db *badger.DB) *ChainRepo {
	return &ChainRepo{
		db: db,
	}
}

func (repo *ChainRepo) LastBlockHeader(ctx context.Context) (*domain.BlockHeader, error) {
	var blockHeader *domain.BlockHeader
	var err error

	err = repo.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(domain.LastBlockHeaderKey())
		if err != nil {
			if errors.Is(err, badger.ErrKeyNotFound) {
				return fmt.Errorf("repo: can't find last block header. err: %w", domain.ErrNotFound)
			}
			return err
		}

		return item.Value(func(val []byte) error {
			decoder := json.NewDecoder(bytes.NewReader(val))
			return decoder.Decode(&blockHeader)
		})
	})

	if err != nil {
		return nil, err
	}

	return blockHeader, nil
}

func (repo *ChainRepo) StoreLastBlockHeaderTX(ctx context.Context, header *domain.BlockHeader, txn *badger.Txn) error {
	val, err := header.Serialize()
	if err != nil {
		return err
	}

	err = txn.Set(domain.LastBlockHeaderKey(), val)
	if err != nil {
		return err
	}

	return nil
}
