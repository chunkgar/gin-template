package server

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/chunkgar/gin-template/internal/code"
	"github.com/chunkgar/gokit/options"
	"github.com/gin-gonic/gin"
)

func getUserID(c *gin.Context) uint {
	userIDStr := c.GetString("userID")
	userIDU64, _ := strconv.ParseUint(userIDStr, 10, 64)
	return uint(userIDU64)
}

func newJwtMiddleware(j *options.JWT) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			writeResponse(c, code.ErrInvalidAuthHeader, "the `Authorization` header was empty")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			writeResponse(c, code.ErrInvalidAuthHeader, "invalid authorization header format")
			c.Abort()
			return
		}

		claims, err := j.ParseJWT(parts[1])
		if claims != nil && claims.Subject == "" {
			err = fmt.Errorf("no subject")
		}
		if err != nil {
			writeResponse(c, code.ErrInvalidAuthHeader, "invalid token: "+err.Error())
			c.Abort()
			return
		}

		c.Set("claims", claims)
		c.Set("jwt", parts[1])
		if claims != nil {
			c.Set("userID", claims.Subject)
		}
		c.Next()
	}
}

func newUserJwtMiddleware(j *options.JWT) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: 额外检查用户是否已删除
	}
}

func newAdminJwtMiddleware(j *options.JWT) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			writeResponse(c, code.ErrInvalidAuthHeader, "the `Authorization` header was empty")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			writeResponse(c, code.ErrInvalidAuthHeader, "invalid authorization header format")
			c.Abort()
			return
		}

		claims, err := j.ParseJWT(parts[1])
		if claims != nil && claims.Subject == "" {
			err = fmt.Errorf("no subject")
		}
		if err != nil {
			writeResponse(c, code.ErrInvalidAuthHeader, "invalid token: "+err.Error())
			c.Abort()
			return
		}
		if claims.Role != "admin" {
			writeResponse(c, code.ErrInvalidAuthHeader, "permission denied")
			c.Abort()
			return
		}

		c.Set("claims", claims)
		c.Set("jwt", parts[1])
		if claims != nil {
			c.Set("userID", claims.Subject)
		}
		c.Next()
	}
}
