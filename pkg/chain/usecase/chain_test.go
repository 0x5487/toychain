package usecase

import (
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"log"
	"testing"
	"toychain/pkg/chain/repository"
	"toychain/pkg/domain"

	badger "github.com/dgraph-io/badger/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateBlockChain(t *testing.T) {
	ctx := context.Background()

	opt := badger.DefaultOptions("").WithInMemory(true)
	db, err := badger.Open(opt)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	accountRepo := repository.NewAccountRepo(db)
	blockRepo := repository.NewBlockRepo(db)
	chainRepo := repository.NewChainRepo(db)

	usecase := NewChainUsecase(db, accountRepo, blockRepo, chainRepo)

	err = usecase.Initialize(ctx)
	require.NoError(t, err)

	accountKey := []byte(domain.GodAddress)

	_ = db.View(func(txn *badger.Txn) error {
		account, err := accountRepo.AccountTX(ctx, accountKey, txn)
		require.NoError(t, err)

		assert.Equal(t, domain.GodAddress, account.Address)
		assert.Equal(t, uint64(1000000), account.Balance)

		block, err := blockRepo.BlockByHeightTX(ctx, 1, txn)
		require.NoError(t, err)

		assert.Equal(t, 1, len(block.Transactions))

		tx := block.Transactions[0]
		assert.Equal(t, "", tx.From)
		assert.Equal(t, domain.GodAddress, tx.Payload.To)
		assert.Equal(t, uint64(1000000), tx.Payload.Value)

		return nil
	})
}

func TestTransfer(t *testing.T) {
	ctx := context.Background()

	opt := badger.DefaultOptions("").WithInMemory(true)
	db, err := badger.Open(opt)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	accountRepo := repository.NewAccountRepo(db)
	blockRepo := repository.NewBlockRepo(db)
	chainRepo := repository.NewChainRepo(db)

	usecase := NewChainUsecase(db, accountRepo, blockRepo, chainRepo)

	err = usecase.Initialize(ctx)
	require.NoError(t, err)

	tx := domain.Transaction{
		Type: 1,
		From: "LteoBMKhjHjV14rEcLF154CPs9BY6JYqexDwkMQc2cGEWDvAv",
		Payload: domain.Payload{
			To:         "2vii4rSRYLgUP56q8jYUP6wQpNevqBDveLTrk1qB86gPEeHp3L",
			Value:      1000,
			GasLimited: 10,
			GasPrice:   1,
			Nonce:      1,
		},
	}

	bytePrivKey, err := hex.DecodeString("e34f9a593a9332b97395a60319fbf7186ca7f59168a5d73cdf203bfcdc02b0e12d29f680b4136828cd53abe747e6b023033d414d260fb64307d4ff8f3bf746e3")
	require.NoError(t, err)

	bPayload, err := tx.Payload.Serialize()
	require.NoError(t, err)

	signedTX := ed25519.Sign(bytePrivKey, bPayload)
	tx.Signature = hex.EncodeToString(signedTX)

	txID, err := usecase.AddPendingTransaction(ctx, &tx)
	t.Log(txID)
	require.NoError(t, err)

	_ = db.View(func(txn *badger.Txn) error {
		accountKey := domain.AccountKey("LteoBMKhjHjV14rEcLF154CPs9BY6JYqexDwkMQc2cGEWDvAv")
		fromAccount, err := accountRepo.AccountTX(ctx, accountKey, txn)
		require.NoError(t, err)
		assert.Equal(t, uint64(999000), fromAccount.Balance)

		accountKey = domain.AccountKey("2vii4rSRYLgUP56q8jYUP6wQpNevqBDveLTrk1qB86gPEeHp3L")
		toAccount, err := accountRepo.AccountTX(ctx, accountKey, txn)
		require.NoError(t, err)
		assert.Equal(t, uint64(1000), toAccount.Balance)

		return nil
	})
}
