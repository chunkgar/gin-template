package server

import (
	"github.com/chunkgar/gin-template/internal/options"
	"github.com/chunkgar/gin-template/internal/store/mysql"
	"github.com/chunkgar/gokit/log"
	genericoptions "github.com/chunkgar/gokit/options"
	"github.com/chunkgar/gokit/shutdown"
	"github.com/chunkgar/gokit/shutdown/shutdownmanagers/posixsignal"
	"github.com/gin-gonic/gin"
)

var engine *gin.Engine = gin.Default()

type Server struct {
	// 选项
	opts *options.Options

	// 关闭管理器
	gs *shutdown.GracefulShutdown

	jwt *genericoptions.JWT
}

func NewServer(opts *options.Options) *Server {
	gs := shutdown.New()
	gs.AddShutdownManager(posixsignal.NewPosixSignalManager())
	gs.AddShutdownCallback(shutdown.ShutdownFunc(func(s string) error {
		mysqlStore, _ := mysql.GetMySQLFactoryOr(nil)
		if mysqlStore != nil {
			_ = mysqlStore.Close()
		}

		return nil
	}))

	gin.SetMode(opts.Server.Mode)

	return &Server{
		opts: opts,
		gs:   gs,
		jwt:  opts.JWT.NewJwt(),
	}
}

func (s *Server) Prepare() *Server {
	return s
}

func (s *Server) Run() error {
	// start shutdown managers
	if err := s.gs.Start(); err != nil {
		log.Fatalf("start shutdown manager failed: %s", err.Error())
	}

	return engine.Run(s.opts.Server.Address())
}
