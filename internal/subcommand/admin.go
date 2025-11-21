package subcommand

import (
	"fmt"

	"github.com/chunkgar/gin-template/internal/options"
	"github.com/chunkgar/gin-template/internal/store/model"
	"github.com/chunkgar/gokit/app"
	cliflag "github.com/marmotedu/component-base/pkg/cli/flag"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

func newCreate(opts *options.Options) *app.Command {

	return app.NewCommand(
		"create",
		"create [username] [password]",
		app.WithCommandRunFunc(func(args []string) error {
			if err := unmarshalOptions(opts); err != nil {
				return err
			}

			// 检查args长度
			if len(args) != 2 {
				return fmt.Errorf("create command requires username and password")
			}

			mysqlOptions := opts.MySQL
			db, err := mysqlOptions.NewClient()
			if err != nil {
				return errors.Wrapf(err, "failed to create mysql client")
			}

			username := args[0]
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(args[1]), bcrypt.DefaultCost)
			if err != nil {
				return errors.Wrapf(err, "failed to hash password")
			}

			admin := &model.AdminUser{
				Username: username,
				Password: string(hashedPassword),
				Role:     "admin",
				IsActive: true,
			}

			if err := db.Create(admin).Error; err != nil {
				return errors.Wrapf(err, "failed to create admin account")
			}

			fmt.Printf("Admin account created successfully: %s\n", username)

			return nil
		}),
	)
}

type SubOptions struct {
	Msg   string
	Times int
}

func (o *SubOptions) Flags() (fss cliflag.NamedFlagSets) {
	fs := fss.FlagSet("subcommand")
	fs.StringVar(&o.Msg, "msg", o.Msg, "message to print")
	fs.IntVar(&o.Times, "times", o.Times, "number of times to print")
	return fss
}

func (o *SubOptions) Validate() []error { return nil }

func Admin(opts *options.Options) *app.Command {
	subOpts := &SubOptions{Msg: "Hello, World from subcommand!", Times: 1}

	cmd := app.NewCommand(
		"admin",
		"admin commands",
		app.WithCommandOptions(subOpts),
	)

	cmd.AddCommand(newCreate(opts))

	return cmd
}
