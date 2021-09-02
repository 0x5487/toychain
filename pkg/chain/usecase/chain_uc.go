package usecase

import (
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"errors"
	"sync"
	"time"
	"toychain/pkg/domain"

	badger "github.com/dgraph-io/badger/v3"
	"github.com/nite-coder/blackbear/pkg/config"
	"github.com/nite-coder/blackbear/pkg/log"
)

type chainUsecase struct {
	db                   *badger.DB
	accountRepo          domain.AccountRepository
	blockRepo            domain.BlockRepository
	chainRepo            domain.ChainRepository
	transactionPool      []*domain.Transaction
	transactionThreshold uint32
	mu                   sync.Mutex
	validator            string
}

func NewChainUsecase(db *badger.DB, accountRepo domain.AccountRepository, blockRepo domain.BlockRepository, chainRepo domain.ChainRepository) domain.ChainUsecase {
	validator, _ := config.String("app.validator", "")

	return &chainUsecase{
		db:                   db,
		accountRepo:          accountRepo,
		blockRepo:            blockRepo,
		chainRepo:            chainRepo,
		transactionPool:      make([]*domain.Transaction, 0),
		transactionThreshold: 1,
		validator:            validator,
	}
}

func (uc *chainUsecase) Initialize(ctx context.Context) error {
	lastBlockHeader, err := uc.chainRepo.LastBlockHeader(ctx)

	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		return err
	}

	if lastBlockHeader != nil {
		return nil
	}

	// create new chain
	tx := domain.NewTransaction("", domain.GodAddress, 1000000)

	block := domain.Block{
		Header: &domain.BlockHeader{
			Height:    1,
			Timestamp: uint64(time.Now().Unix()),
			Validator: uc.validator,
		},
	}
	block.AddTransaction(tx)

	accountKey := []byte(domain.GodAddress)

	return uc.db.Update(func(txn *badger.Txn) error {
		account, err := uc.accountRepo.AccountTX(ctx, accountKey, txn)
		if err != nil && !errors.Is(err, domain.ErrNotFound) {
			return err
		}

		if account != nil {
			return errors.New("GOD Account has already exist")
		}

		if account == nil {
			account = &domain.Account{
				Address: domain.GodAddress,
				Balance: tx.Payload.Value,
			}

			err = uc.accountRepo.StoreTX(ctx, account, txn)
			if err != nil {
				return err
			}
		}

		err = uc.blockRepo.StoreTX(ctx, &block, txn)
		if err != nil {
			return err
		}

		err = uc.chainRepo.StoreLastBlockHeaderTX(ctx, block.Header, txn)
		if err != nil {
			return err
		}

		return nil
	})
}

func (uc *chainUsecase) LastBlockHeader(ctx context.Context) (*domain.BlockHeader, error) {
	return uc.chainRepo.LastBlockHeader(ctx)
}

func (uc *chainUsecase) AddPendingTransaction(ctx context.Context, tx *domain.Transaction) (string, error) {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	tx.Timestamp = uint64(time.Now().Unix())
	tx.Type = 1
	uc.transactionPool = append(uc.transactionPool, tx)
	txID := tx.ID().String()

	if len(uc.transactionPool) >= int(uc.transactionThreshold) {
		block, err := uc.packTransactionToBlock(ctx, uc.transactionPool)
		uc.transactionPool = []*domain.Transaction{}
		if err != nil {
			return "", err
		}

		if block != nil && len(block.Transactions) > 0 {
			// TODO: broadcast block to other nodes.
		}

		return txID, nil
	}

	// TODO: broadcast transaction to other nodes.
	return txID, nil
}

func (uc *chainUsecase) packTransactionToBlock(ctx context.Context, pendingTransactions []*domain.Transaction) (*domain.Block, error) {
	logger := log.FromContext(ctx)

	block := domain.Block{}

	// parse transactions
	err := uc.db.Update(func(txn *badger.Txn) error {
		lastHeader, err := uc.chainRepo.LastBlockHeader(ctx)
		block.Header = &domain.BlockHeader{
			Height:    lastHeader.Height + 1,
			Timestamp: uint64(time.Now().Unix()),
			PrevHash:  lastHeader.ID().String(),
			Validator: uc.validator,
		}

		logger.Debugf("pending transaction count: %d", len(uc.transactionPool))

		for i := 0; i < len(pendingTransactions); i++ {
			tx := pendingTransactions[i]

			// valid transaction
			pubKey, err := domain.AddressToPubKey(tx.From)
			if err != nil {
				return err
			}

			msg, err := tx.Payload.Serialize()
			if err != nil {
				return err
			}

			bPayload, err := hex.DecodeString(tx.Signature)
			if err != nil {
				return err
			}

			if !ed25519.Verify(pubKey, msg, bPayload) {
				logger.Debugf("%v", tx.Payload)
				return domain.ErrInvalidSignature
			}

			// change fromAccount balance
			accountKey := domain.AccountKey(tx.From)
			fromAccount, err := uc.accountRepo.AccountTX(ctx, accountKey, txn)
			if err != nil {
				return err
			}

			if fromAccount.Nonce+1 != tx.Payload.Nonce {
				logger.Debugf("nonce should be %d, but is %d", fromAccount.Nonce+1, tx.Payload.Nonce)
				continue
			}

			if fromAccount.Balance < tx.Payload.Value {
				logger.Debugf("not enough balance")
				continue
			}

			fromAccount.Balance -= tx.Payload.Value
			fromAccount.Nonce += 1
			err = uc.accountRepo.StoreTX(ctx, fromAccount, txn)
			if err != nil {
				return err
			}

			// change toAccount balance
			accountKey = domain.AccountKey(tx.Payload.To)
			toAccount, err := uc.accountRepo.AccountTX(ctx, accountKey, txn)
			if err != nil {
				if errors.Is(err, domain.ErrNotFound) {
					toAccount = &domain.Account{
						Address: tx.Payload.To,
						Balance: 0,
					}
				} else {
					return err
				}
			}

			toAccount.Balance += tx.Payload.Value
			err = uc.accountRepo.StoreTX(ctx, toAccount, txn)
			if err != nil {
				return err
			}

			block.AddTransaction(tx)
		}

		if len(block.Transactions) == 0 {
			logger.Debug("no transaction in the block")
			return nil
		}

		err = uc.blockRepo.StoreTX(ctx, &block, txn)
		if err != nil {
			return err
		}

		logger.Debugf("last blcok number: %d, ", block.Header.Height)
		err = uc.chainRepo.StoreLastBlockHeaderTX(ctx, block.Header, txn)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &block, nil
}
