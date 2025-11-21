package options

import (
	"encoding/json"

	"github.com/chunkgar/gokit/log"
	genericoptions "github.com/chunkgar/gokit/options"
	cliflag "github.com/marmotedu/component-base/pkg/cli/flag"
)

type Options struct {
	Server    *genericoptions.ServerOptions `json:"server" mapstructure:"server"`
	MySQL     *genericoptions.MySQLOptions  `json:"mysql" mapstructure:"mysql"`
	Redis     *genericoptions.RedisOptions  `json:"redis" mapstructure:"redis"`
	JWT       *genericoptions.JwtOptions    `json:"jwt" mapstructure:"jwt"`
	Log       *log.Options                  `json:"log" mapstructure:"log"`
	AppleAuth *AppleAuthOptions             `json:"apple-auth" mapstructure:"apple-auth"`
}

// NewOptions creates a new Options object with default parameters.
func NewOptions() *Options {
	o := Options{
		Log:       log.NewOptions(),
		Server:    genericoptions.NewServerOptions(),
		MySQL:     genericoptions.NewMySQLOptions(),
		Redis:     genericoptions.NewRedisOptions(),
		JWT:       genericoptions.NewJWTOptions(),
		AppleAuth: NewAppleAuthOptions(),
	}

	return &o
}

// ApplyTo applies the run options to the method receiver and returns self.
// func (o *Options) ApplyTo(c *server.Config) error {
// 	return nil
// }

// Flags returns flags for a specific APIServer by section name.
func (o *Options) Flags() (fss cliflag.NamedFlagSets) {
	o.Log.AddFlags(fss.FlagSet("logs"))
	o.Server.AddFlags(fss.FlagSet("server"))
	o.MySQL.AddFlags(fss.FlagSet("mysql"))
	o.Redis.AddFlags(fss.FlagSet("redis"))
	o.JWT.AddFlags(fss.FlagSet("jwt"))
	o.AppleAuth.AddFlags(fss.FlagSet("apple-auth"))

	return fss
}

func (o *Options) String() string {
	data, _ := json.Marshal(o)

	return string(data)
}

// Complete set default Options.
func (o *Options) Complete() error {
	return nil
	// if o.JwtOptions.Key == "" {
	// 	o.JwtOptions.Key = idutil.NewSecretKey()
	// }

	// return o.SecureServing.Complete()
}

func (o *Options) Validate() []error {
	var errs []error

	errs = append(errs, o.Log.Validate()...)
	errs = append(errs, o.Server.Validate()...)
	errs = append(errs, o.Redis.Validate()...)
	errs = append(errs, o.MySQL.Validate()...)
	errs = append(errs, o.JWT.Validate()...)
	errs = append(errs, o.AppleAuth.Validate()...)

	return errs
}
