package delivery

import (
	"crypto/ed25519"
	"encoding/hex"
	"toychain/pkg/domain"

	"github.com/nite-coder/blackbear/pkg/log"
	"github.com/nite-coder/blackbear/pkg/web"
)

type ChainHandler struct {
	chainUsecase   domain.ChainUsecase
	accountUsecase domain.AccountUsecase
	blockUsecase   domain.BlockUseucase
}

func NewChainHandler(chainUsecase domain.ChainUsecase, accountUsecase domain.AccountUsecase, blockUsecase domain.BlockUseucase) *ChainHandler {
	return &ChainHandler{
		chainUsecase:   chainUsecase,
		accountUsecase: accountUsecase,
		blockUsecase:   blockUsecase,
	}
}

func (h *ChainHandler) GetAccountEndpoint(c *web.Context) error {
	ctx := c.Request.Context()
	address := c.Param("address")

	amount, err := h.accountUsecase.Account(ctx, address)
	if err != nil {
		return err
	}

	return c.JSON(200, amount)
}

func (h *ChainHandler) GetLastBlockHeaderEndpoint(c *web.Context) error {
	ctx := c.Request.Context()

	header, err := h.chainUsecase.LastBlockHeader(ctx)
	if err != nil {
		return err
	}

	return c.JSON(200, header)
}

func (h *ChainHandler) SendTransactionEndpoint(c *web.Context) error {
	ctx := c.Request.Context()
	logger := log.FromContext(ctx)
	logger.Debug("=== SendTransactionEndpoint ===")

	tx := domain.Transaction{}
	err := c.BindJSON(&tx)
	if err != nil {
		return err
	}

	logger.Debugf("%v", tx)
	txID, err := h.chainUsecase.AddPendingTransaction(ctx, &tx)
	if err != nil {
		return err
	}

	return c.JSON(200, txID)
}

type EncryptPayload struct {
	PrivateKey string         `json:"private_key"`
	Payload    domain.Payload `json:"payload"`
}

func (h *ChainHandler) EncryptPayloadEndpoint(c *web.Context) error {
	ctx := c.Request.Context()
	logger := log.FromContext(ctx)
	logger.Debug("=== EncryptPayloadEndpoint ===")

	payload := EncryptPayload{}
	err := c.BindJSON(&payload)
	if err != nil {
		return err
	}

	bPrivKey, err := hex.DecodeString(payload.PrivateKey)

	b, err := payload.Payload.Serialize()
	if err != nil {
		return err
	}
	signedPayload := ed25519.Sign(bPrivKey, b)
	result := hex.EncodeToString(signedPayload)

	return c.String(200, result)
}

func (h *ChainHandler) GetBlockHeightEndpoint(c *web.Context) error {
	ctx := c.Request.Context()
	logger := log.FromContext(ctx)
	logger.Debug("=== GetBlockHeightEndpoint ===")

	height, _ := c.ParamInt("height")

	block, err := h.blockUsecase.BlockByHeight(ctx, uint64(height))
	if err != nil {
		return err
	}

	return c.JSON(200, block)
}

func (h *ChainHandler) CreateAccountEndpoint(c *web.Context) error {
	ctx := c.Request.Context()
	logger := log.FromContext(ctx)
	logger.Debug("=== CreateAccountEndpoint ===")

	account, err := h.accountUsecase.GenerateAccount(ctx)
	if err != nil {
		return err
	}

	return c.JSON(200, account)
}
