package usecase

import (
	"context"
	"toychain/pkg/domain"

	badger "github.com/dgraph-io/badger/v3"
)

type BlockUseucase struct {
	db        *badger.DB
	blockRepo domain.BlockRepository
}

func NewBlockUsecase(db *badger.DB, blockRepo domain.BlockRepository) *BlockUseucase {
	return &BlockUseucase{
		db:        db,
		blockRepo: blockRepo,
	}
}

func (uc *BlockUseucase) BlockByHeight(ctx context.Context, height uint64) (*domain.Block, error) {
	var block *domain.Block
	var err error

	err = uc.db.View(func(txn *badger.Txn) error {
		block, err = uc.blockRepo.BlockByHeightTX(ctx, height, txn)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return block, nil
	}

	return block, nil
}
