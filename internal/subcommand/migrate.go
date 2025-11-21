package subcommand

import (
	"github.com/chunkgar/gin-template/internal/options"
	"github.com/chunkgar/gin-template/internal/store/model"
	"github.com/chunkgar/gokit/app"
	"github.com/chunkgar/gokit/log"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// NOTE: models 数据库模型列表
var models = []interface{}{
	&model.User{},
	&model.UserAccount{},
	&model.UserOAuthToken{},
	&model.Membership{},
	&model.AdminUser{},
	&model.UserDeletion{},
}

func Migrate(opts *options.Options) *app.Command {
	return app.NewCommand(
		"migrate",
		"migrate database",
		app.WithCommandRunFunc(func(args []string) error {
			if err := viper.Unmarshal(opts); err != nil {
				return err
			}

			log.Infof("migrate the database")
			mysqlOptions := opts.MySQL
			log.Infof("Mysql: %+v\n", mysqlOptions)

			db, err := mysqlOptions.NewClient()
			if err != nil {
				log.Errorf("failed to create mysql client: %v", err)
				return err
			}

			for _, model := range models {
				if err := db.AutoMigrate(model); err != nil {
					return errors.Wrapf(err, "failed to migrate model: %v", model)
				}
			}

			log.Info("migrate the database success")

			return nil
		}),
	)
}
