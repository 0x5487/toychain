package node

import (
	"context"
	"fmt"
	"toychain/internal/pkg/startup"
	"toychain/pkg/chain/delivery"
	"toychain/pkg/chain/repository"
	"toychain/pkg/chain/usecase"

	badger "github.com/dgraph-io/badger/v3"
	"github.com/nite-coder/blackbear/pkg/log"
	"github.com/nite-coder/blackbear/pkg/web"
)

var (
	db        *badger.DB
	webServer *web.WebServer
)

func initialize() error {
	ctx := context.Background()

	err := startup.InitConfig()
	if err != nil {
		return err
	}

	err = startup.InitLogger()
	if err != nil {
		return err
	}

	db, err = startup.InitBadger()
	if err != nil {
		return err
	}

	// repository
	accountRepo := repository.NewAccountRepo(db)
	blockRepo := repository.NewBlockRepo(db)
	chainRepo := repository.NewChainRepo(db)

	// usecuase
	accountUsecase := usecase.NewAccountUsecase(db, accountRepo)
	blockUsecase := usecase.NewBlockUsecase(db, blockRepo)
	chainUsecase := usecase.NewChainUsecase(db, accountRepo, blockRepo, chainRepo)
	err = chainUsecase.Initialize(ctx)
	if err != nil {
		return err
	}

	// handlers
	chainHandler := delivery.NewChainHandler(chainUsecase, accountUsecase, blockUsecase)

	// web server
	webServer = web.NewServer()
	webServer.ErrorHandler = func(c *web.Context, err error) {
		c.JSON(400, err.Error())
		fmt.Println(err)
	}

	delivery.RegisterChainRouter(webServer, chainHandler)

	log.Info("main: toychain node is initialized")
	return nil
}
