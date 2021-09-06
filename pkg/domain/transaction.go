package domain

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"time"
)

type Payload struct {
	To         string `json:"to"`
	Value      uint64 `json:"value"`
	GasLimited uint64 `json:"gas_limited"`
	GasPrice   uint64 `json:"gas_price"`
	Nonce      uint64 `json:"nonce"`
}

func (p *Payload) Serialize() ([]byte, error) {
	result, err := json.Marshal(&p)
	if err != nil {
		return nil, err
	}

	return result, nil
}

type Transaction struct {
	Type      uint32  `json:"type"` // stake, transaction, reward
	Timestamp uint64  `json:"timestamp"`
	From      string  `json:"from"`
	Payload   Payload `json:"payload"`
	Signature string  `json:"signature"`
}

func NewTransaction(from, to string, value uint64) *Transaction {
	return &Transaction{
		Type:      1,
		Timestamp: uint64(time.Now().Unix()),
		From:      from,
		Payload: Payload{
			To:    to,
			Value: value,
			Nonce: 0,
		},
	}
}

func (tx *Transaction) ID() Identifier {
	b, _ := tx.Serialize()
	return sha256.Sum256(b)
}

func (tx *Transaction) Serialize() ([]byte, error) {
	result, err := json.Marshal(&tx)
	if err != nil {
		return nil, err
	}

	return result, nil
}

type Receipt struct {
	TransactionID Identifier
	BlockHeight   uint64
	State         uint32 // pending, success, failed
}

func ComputeTransactionRoot(transactions []*Transaction) string {
	bb := [][]byte{}
	for _, tx := range transactions {
		bb = append(bb, tx.ID().ToByte())
	}

	b := bytes.Join(bb, []byte{})
	bHash := sha256.Sum256(b)
	return hex.EncodeToString(bHash[:])
}

type TransactionUsecase interface {
	TransactionByID(ctx context.Context, id string) (*Transaction, error)
	AddPendingTransaction(ctx context.Context, tx *Transaction) error
	PackPendingTransaction(ctx context.Context) error
}
