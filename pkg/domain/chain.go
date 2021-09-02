package domain

import (
	"context"

	badger "github.com/dgraph-io/badger/v3"
)

const (
	GodAddress = "LteoBMKhjHjV14rEcLF154CPs9BY6JYqexDwkMQc2cGEWDvAv"
)

var (
	lastBlockHeader = []byte("last_block_header")
)

type Chain struct {
	id             uint32
	blockChan      chan Block
	tokenRewardNum uint64
	validatorFee   uint64
}

func NewChain(chainID uint32) *Chain {
	return &Chain{
		id:             chainID,
		blockChan:      make(chan Block, 10),
		tokenRewardNum: 1,
		validatorFee:   3,
	}
}

func LastBlockHeaderKey() []byte {
	return lastBlockHeader
}

type ChainUsecase interface {
	Initialize(ctx context.Context) error
	LastBlockHeader(ctx context.Context) (*BlockHeader, error)
	AddPendingTransaction(ctx context.Context, tx *Transaction) (string, error)
}

type ChainRepository interface {
	LastBlockHeader(ctx context.Context) (*BlockHeader, error)
	StoreLastBlockHeaderTX(ctx context.Context, header *BlockHeader, txn *badger.Txn) error
}
