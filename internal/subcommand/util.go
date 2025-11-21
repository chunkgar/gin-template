package subcommand

import (
	"github.com/chunkgar/gin-template/internal/options"
	"github.com/spf13/viper"
)

func unmarshalOptions(opts *options.Options) error {
	if err := viper.Unmarshal(opts); err != nil {
		return err
	}
	return nil
}
