package server

import (
	"strconv"

	"github.com/chunkgar/gin-template/internal/code"
	"github.com/chunkgar/gin-template/internal/store"
	"github.com/chunkgar/gokit/log"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func initAdmin(s *Server) {
	// 路由：管理员登录
	s.engine.POST("/api/admin/login", s.postAdminLogin)

	// 路由：管理员认证
	{
		jwtMiddle := newAdminJwtMiddleware(s.jwt)
		r := s.engine.Group("/api/admin/auth")
		r.Use(jwtMiddle)
		r.GET("/profile", s.getAdminProfile)
	}
}

func (s *Server) postAdminLogin(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Errorf("failed to bind json: %v", err)
		writeResponse(c, code.ErrBind, nil, err.Error())
		return
	}

	adminUser, err := store.Client().AdminUser().GetByName(req.Username)
	if err != nil {
		log.Errorf("failed to get admin user by name: %v", err)
		writeResponse(c, code.ErrAccountNotFound, nil, "invalid username or password")
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(adminUser.Password), []byte(req.Password)); err != nil {
		log.Errorf("failed to compare password: %v", err)
		writeResponse(c, code.ErrAccountInvalidPassword, nil)
		return
	}

	token, expiresIn, err := s.jwt.GenerateJWT(strconv.FormatUint(uint64(adminUser.ID), 10), "admin")
	if err != nil {
		log.Errorf("failed to generate jwt: %v", err)
		writeResponse(c, code.ErrInternalServer, nil, err.Error())
		return
	}

	writeResponse(c, code.ErrSuccess, map[string]any{
		"token":     token,
		"expiresIn": expiresIn,
	})
}

func (s *Server) getAdminProfile(c *gin.Context) {
	userID := getUserID(c)

	adminUser, err := store.Client().AdminUser().GetByID(userID)
	if err != nil {
		log.Errorf("failed to get admin user by id: %v", err)
		writeResponse(c, code.ErrAccountNotFound, nil, "admin user not found")
		return
	}

	writeResponse(c, code.ErrSuccess, gin.H{
		"id":       adminUser.ID,
		"username": adminUser.Username,
		"role":     adminUser.Role,
		"isActive": adminUser.IsActive,
	})
}
