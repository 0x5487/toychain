package delivery

import "github.com/nite-coder/blackbear/pkg/web"

func RegisterChainRouter(webServer *web.WebServer, h *ChainHandler) {

	webServer.Get("/accounts/:address", h.GetAccountEndpoint)
	webServer.Post("/accounts", h.CreateAccountEndpoint)

	webServer.Post("/transactions/encrypt", h.EncryptPayloadEndpoint)
	webServer.Post("/transactions", h.SendTransactionEndpoint)

	webServer.Get("/blocks/last-block-header", h.GetLastBlockHeaderEndpoint)
	webServer.Get("/blocks/:height", h.GetBlockHeightEndpoint)

}
