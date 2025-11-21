package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/chunkgar/gin-template/internal/options"
	"github.com/chunkgar/gin-template/internal/server"
	"github.com/chunkgar/gin-template/internal/store"
	"github.com/chunkgar/gin-template/internal/store/mysql"
	"github.com/chunkgar/gin-template/internal/subcommand"
	"github.com/chunkgar/gokit/app"
	"github.com/chunkgar/gokit/log"
)

var (
	NAME     = "main"
	BASENAME = "main"
	DESC     = "gin后端"
)

func main() {
	if len(os.Getenv("GOMAXPROCS")) == 0 {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}

	opts := options.NewOptions()

	app.NewApp(NAME, BASENAME,
		app.WithDescription(DESC),
		app.WithOptions(opts),
		app.WithSubCommand(subcommand.Migrate(opts)),
		app.WithSubCommand(subcommand.Admin(opts)),
		app.WithRunFunc(makeRunFunc(opts)),
	).Run()
}

func makeRunFunc(opts *options.Options) app.RunFunc {
	return func(basename string) error {
		fmt.Printf("Log level: %s\n", opts.Log.Level)
		// init
		log.Init(opts.Log)
		defer log.Flush()

		// init mysql
		storeIns, err := mysql.GetMySQLFactoryOr(opts.MySQL)
		if err != nil {
			log.Fatalf("failed to get mysql factory: %v", err)
		}
		store.SetClient(storeIns)

		return server.NewServer(opts).Prepare().Run()
	}
}
