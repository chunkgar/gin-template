package mysql

import (
	"fmt"
	"sync"

	"github.com/chunkgar/gin-template/internal/store"
	genericoptions "github.com/chunkgar/gokit/options"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type datastore struct {
	db *gorm.DB

	// can include two database instance if needed
	// docker *grom.DB
	// db *gorm.DB
}

func (ds *datastore) Close() error {
	db, err := ds.db.DB()
	if err != nil {
		return errors.Wrap(err, "get gorm db instance failed")
	}

	return db.Close()
}

var (
	mysqlFactory store.Factory
	once         sync.Once
)

// GetMySQLFactoryOr create mysql factory with the given config.
func GetMySQLFactoryOr(opts *genericoptions.MySQLOptions) (store.Factory, error) {
	if opts == nil && mysqlFactory == nil {
		return nil, fmt.Errorf("failed to get mysql store fatory")
	}

	var err error
	var dbIns *gorm.DB
	once.Do(func() {
		dbIns, err = opts.NewClient()

		mysqlFactory = &datastore{dbIns}
	})

	if mysqlFactory == nil || err != nil {
		return nil, fmt.Errorf("failed to get mysql store fatory, mysqlFactory: %+v, error: %w", mysqlFactory, err)
	}

	return mysqlFactory, nil
}

func (ds *datastore) UserAccount() store.UserAccount {
	return &userAccountStore{ds.db}
}

func (ds *datastore) User() store.User {
	return &userStore{ds.db}
}

func (ds *datastore) AdminUser() store.AdminUser {
	return &adminUserStore{ds.db}
}
