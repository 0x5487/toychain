package domain

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"time"

	badger "github.com/dgraph-io/badger/v3"
)

type BlockHeader struct {
	Height          uint64 `json:"height"`
	Timestamp       uint64 `json:"timestamp"`
	PrevHash        string `json:"prev_hash"`
	Validator       string `json:"validator"`
	TransactionRoot string `json:"transaction_root"`
}

func (h *BlockHeader) ID() Identifier {
	b, _ := h.Serialize()
	return sha256.Sum256(b)
}

func (b *BlockHeader) Serialize() ([]byte, error) {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)

	err := encoder.Encode(b)
	if err != nil {
		return nil, err
	}

	return res.Bytes(), nil
}

func DeserializeBlock(data []byte) (*Block, error) {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(data))

	err := decoder.Decode(&block)
	if err != nil {
		return nil, err
	}
	return &block, nil
}

type Block struct {
	Header       *BlockHeader   `json:"header"`
	Transactions []*Transaction `json:"transactions"`
}

func NewBlock(prevHash string, height uint64) *Block {
	block := Block{
		Header: &BlockHeader{
			Height:    height,
			PrevHash:  prevHash,
			Timestamp: uint64(time.Now().Unix()),
		},
	}

	return &block
}

func (b *Block) ID() Identifier {
	return b.Header.ID()
}

func (b *Block) Serialize() ([]byte, error) {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)

	err := encoder.Encode(b)
	if err != nil {
		return nil, err
	}

	return res.Bytes(), nil
}

func (b *Block) AddTransaction(tx *Transaction) {
	b.Transactions = append(b.Transactions, tx)

	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)
	_ = encoder.Encode(b.Transactions)

	aa := sha256.Sum256(res.Bytes())
	b.Header.TransactionRoot = hex.EncodeToString(aa[:]) // TODO: use merkle tree later
}

type BlockUseucase interface {
	//BlockByID(ctx context.Context, id string) (*Block, error)
	BlockByHeight(ctx context.Context, height uint64) (*Block, error)
}

type BlockRepository interface {
	BlockTX(ctx context.Context, key []byte, txn *badger.Txn) (*Block, error)
	BlockByHeightTX(ctx context.Context, height uint64, txn *badger.Txn) (*Block, error)
	StoreTX(ctx context.Context, block *Block, tx *badger.Txn) error
}
