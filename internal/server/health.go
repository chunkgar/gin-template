package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func init() {
	engine.GET("/api/health", func(c *gin.Context) {
		remoteIP, _ := c.RemoteIP()
		c.JSON(http.StatusOK, gin.H{
			"clientIP": c.ClientIP(),
			"remoteIP": remoteIP,
			"headers":  c.Request.Header,
		})
	})
}
