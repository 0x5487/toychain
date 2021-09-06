package startup

import (
	"strings"

	badger "github.com/dgraph-io/badger/v3"
	"github.com/nite-coder/blackbear/pkg/config"
	"github.com/nite-coder/blackbear/pkg/log"
)

type Database struct {
	Name             string
	ConnectionString string `mapstructure:"connection_string"`
	Type             string
	Migration        bool `mapstructure:"migration"`
}

func InitBadger() (*badger.DB, error) {
	dbType, _ := config.String("db.type", "memory")

	var opt badger.Options
	switch strings.ToLower(dbType) {
	default:
	case "memory":
		log.Info("badge in memory mode")
		opt = badger.DefaultOptions("").WithInMemory(true)
	case "file":
		log.Info("badge in file mode")
		location, _ := config.String("db.location", "")
		opt = badger.DefaultOptions(location)
	}

	db, err := badger.Open(opt)
	if err != nil {
		return nil, err
	}

	return db, nil
}
