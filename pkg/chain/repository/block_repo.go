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

type BlockRepo struct {
	db *badger.DB
}

func NewBlockRepo(db *badger.DB) *BlockRepo {
	return &BlockRepo{
		db: db,
	}
}

func (repo *BlockRepo) StoreTX(ctx context.Context, block *domain.Block, tx *badger.Txn) error {
	val, err := block.Serialize()
	if err != nil {
		return err
	}

	key := block.ID().ToByte()
	err = tx.Set(key, val)
	if err != nil {
		return err
	}

	// block height key
	blockHeightKey := []byte(fmt.Sprintf("/block_height/%d", block.Header.Height))
	err = tx.Set(blockHeightKey, key)
	if err != nil {
		return err
	}

	return nil
}

func (repo *BlockRepo) BlockByHeightTX(ctx context.Context, height uint64, txn *badger.Txn) (*domain.Block, error) {
	blockHeightKey := []byte(fmt.Sprintf("/block_height/%d", height))
	item, err := txn.Get(blockHeightKey)
	if err != nil {
		if errors.Is(err, badger.ErrKeyNotFound) {
			return nil, fmt.Errorf("repo: can't find block height key. %w", domain.ErrNotFound)
		}
		return nil, err
	}

	var blockKey []byte
	err = item.Value(func(val []byte) error {
		blockKey = append([]byte{}, val...)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return repo.BlockTX(ctx, blockKey, txn)
}

func (repo *BlockRepo) BlockTX(ctx context.Context, key []byte, txn *badger.Txn) (*domain.Block, error) {
	block := domain.Block{}

	item, err := txn.Get(key)
	if err != nil {
		if errors.Is(err, badger.ErrKeyNotFound) {
			return nil, fmt.Errorf("repo: can't find block key. err: %w", domain.ErrNotFound)
		}
		return nil, err
	}

	err = item.Value(func(val []byte) error {
		decoder := json.NewDecoder(bytes.NewReader(val))
		return decoder.Decode(&block)
	})

	if err != nil {
		return nil, err
	}

	return &block, nil
}
