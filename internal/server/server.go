package server

import (
	"net/http"
	"time"

	"github.com/chunkgar/gin-template/internal/options"
	"github.com/chunkgar/gin-template/internal/pkg/auth"
	"github.com/chunkgar/gin-template/internal/pkg/auth/provider"
	"github.com/chunkgar/gin-template/internal/store/mysql"
	"github.com/chunkgar/gokit/log"
	"github.com/chunkgar/gokit/middleware"
	genericoptions "github.com/chunkgar/gokit/options"
	"github.com/chunkgar/gokit/shutdown"
	"github.com/chunkgar/gokit/shutdown/shutdownmanagers/posixsignal"
	"github.com/gin-gonic/gin"
)

type Server struct {
	// gin 引擎
	engine *gin.Engine

	// 选项
	opts *options.Options

	// 关闭管理器
	gs *shutdown.GracefulShutdown

	// jwt
	jwt      *genericoptions.JWT
	verifier auth.TokenVerifier
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
		engine: gin.New(),
		opts:   opts,
		gs:     gs,
		jwt:    opts.JWT.NewJwt(),
	}
}

func (s *Server) Prepare() *Server {
	log.Infof("server address: %s", s.opts.Server.Address())
	initJWT(s)
	initMiddlewares(s)
	// 初始化路由
	initHealth(s)
	initUser(s)
	initAdmin(s)

	return s
}

func (s *Server) Run() error {
	// start shutdown managers
	if err := s.gs.Start(); err != nil {
		log.Fatalf("start shutdown manager failed: %s", err.Error())
	}

	return s.engine.Run(s.opts.Server.Address())
}

func initJWT(s *Server) {
	// JWT 验证器
	fetcher := auth.NewCachedJWKSFetcher(24 * time.Hour)
	verifier := auth.NewTokenVerifier(fetcher)
	verifier.RegisterProvider(provider.NewAppleProvider(s.opts.AppleAuth.ClientID))

	s.verifier = verifier
}

func initMiddlewares(s *Server) {
	middleware.Apply(s.engine, s.opts.Server.Middlewares)
}

func initHealth(s *Server) {
	s.engine.GET("/api/health", func(c *gin.Context) {
		remoteIP, _ := c.RemoteIP()
		c.JSON(http.StatusOK, gin.H{
			"clientIP": c.ClientIP(),
			"remoteIP": remoteIP,
			"headers":  c.Request.Header,
		})
	})
}
