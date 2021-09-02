package startup

import (
	"strings"

	badger "github.com/dgraph-io/badger/v3"
	"github.com/nite-coder/blackbear/pkg/config"
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
	case "memory":
		opt = badger.DefaultOptions("").WithInMemory(true)
	case "file":
		location, _ := config.String("db.location", "")
		opt = badger.DefaultOptions(location)
	}

	db, err := badger.Open(opt)
	if err != nil {
		return nil, err
	}

	return db, nil
}
