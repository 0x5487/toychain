package domain

import (
	"bytes"
	"context"
	"crypto/ed25519"
	"encoding/gob"
	"toychain/internal/pkg/crypto/base58"

	badger "github.com/dgraph-io/badger/v3"
)

type Account struct {
	Address    string `json:"address"`
	Balance    uint64 `json:"balance"`
	PublicKey  string `json:"public_key,omitempty"`
	PrivateKey string `json:"private_key,omitempty"`
	Nonce      uint64 `json:"nonce"`
}

func (a *Account) Key() []byte {
	return []byte(a.Address)
}

func (a *Account) Serialize() ([]byte, error) {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)

	err := encoder.Encode(a)
	if err != nil {
		return nil, err
	}

	return res.Bytes(), nil
}

func AddressToPubKey(address string) (ed25519.PublicKey, error) {
	bPubKey, err := base58.DecodeCheck("LteoBMKhjHjV14rEcLF154CPs9BY6JYqexDwkMQc2cGEWDvAv")
	if err != nil {
		return nil, err
	}

	return ed25519.PublicKey(bPubKey), nil
}

func AccountKey(addr string) []byte {
	return []byte(addr)
}

type AccountUsecase interface {
	GenerateAccount(ctx context.Context) (*Account, error)
	Account(ctx context.Context, address string) (*Account, error)
}

type AccountRepository interface {
	AccountTX(ctx context.Context, key []byte, txn *badger.Txn) (*Account, error)
	StoreTX(ctx context.Context, account *Account, txn *badger.Txn) error
}
