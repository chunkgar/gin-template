package server

import (
	"strconv"

	"github.com/chunkgar/gin-template/internal/code"
	"github.com/chunkgar/gin-template/internal/server/vo"
	"github.com/chunkgar/gin-template/internal/store"
	"github.com/chunkgar/gin-template/internal/store/model"
	"github.com/chunkgar/gokit/log"
	"github.com/gin-gonic/gin"
)

func initUser(s *Server) {
	jwtMiddleware := newJwtMiddleware(s.jwt)

	// 路由：账号绑定与解绑
	{
		r := s.engine.Group("/api/account")
		r.Use(jwtMiddleware)

		r.POST("/bind/idtoken", s.postBindIDToken)
		r.POST("/unbind", s.postUnbind)
	}

	// 路由：用户
	{
		r := s.engine.Group("/api/user")

		// 登录
		r.POST("/login/anon", s.postLoginAnon)
		r.POST("/login/idtoken", s.postLoginIDToken)

		// 信息
		r.GET("/profile", jwtMiddleware, s.getProfile)

		// 删除
		r.POST("/delete", jwtMiddleware, s.postDelete)
	}
}

func (s *Server) postBindIDToken(c *gin.Context) {
	userID := getUserID(c)
	var req vo.BindIDTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Errorw("bindIDToken: ShouldBindJSON", "err", err)
		writeResponse(c, code.ErrInvalidParams, nil)
		return
	}

	claims, provider, err := s.verifier.Verify(c, req.IDToken, req.Nonce)
	if err != nil {
		log.Errorf("login idtoken verify failed: %s", err.Error())
		writeResponse(c, code.ErrSignatureInvalid, nil)
		return
	}

	userAccount, errCode, err := store.Client().UserAccount().Bind(userID, claims.GetSubject(), model.AccountType(provider.GetName()))
	if err != nil {
		log.Errorw("bindIDToken: Bind", "err", err, "userID", userID, "accountID", claims.GetSubject(), "accountType", provider.GetName())
		writeResponse(c, errCode, nil)
		return
	}

	token, expiresIn, err := s.jwt.GenerateJWT(strconv.FormatUint(uint64(userAccount.UserID), 10), "user")
	if err != nil {
		log.Errorw("bindIDToken: GenerateToken", "err", err, "userID", userAccount.UserID, "username", userAccount.User.Nickname)
		writeResponse(c, code.ErrInternalServer, nil)
		return
	}

	writeResponse(c, code.ErrSuccess, vo.LoginResponse{
		Token:     token,
		ExpiresIn: expiresIn.UnixMicro(),
		UserID:    userAccount.UserID,
		Nickname:  userAccount.User.Nickname,
		AvatarURL: userAccount.User.AvatarURL,
	})
}

func (s *Server) postUnbind(c *gin.Context) {
	var req vo.UnbindRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Errorw("unbind: ShouldBindJSON", "err", err)
		writeResponse(c, code.ErrInvalidParams, nil)
		return
	}

	userID := getUserID(c)
	if err := store.Client().UserAccount().Unbind(userID, model.AccountType(req.Type)); err != nil {
		log.Errorw("unbind: Unbind", "err", err, "userID", userID, "accountType", req.Type)
		writeResponse(c, code.ErrInternalServer, nil)
		return
	}

	writeResponse(c, code.ErrSuccess, nil)
}

func (s *Server) postLoginAnon(c *gin.Context) {
	// TODO: Cache

	var req vo.LoginAnonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Errorw("login anon json failed: %s", err.Error(), "deviceID", req.DeviceID, "type", req.Type)
		writeResponse(c, code.ErrInvalidParams, nil)
		return
	}

	userAccount, _, err := store.Client().UserAccount().FindOrCreate(req.DeviceID, model.AccountType(req.Type))
	if err != nil {
		log.Errorw("login anon find or create failed: %s", err.Error(), "deviceID", req.DeviceID, "type", req.Type)
		writeResponse(c, code.ErrInternalServer, nil)
		return
	}

	token, expiresIn, err := s.jwt.GenerateJWT(strconv.FormatUint(uint64(userAccount.UserID), 10), "user")
	if err != nil {
		log.Errorw("login anon generate jwt failed: %s", err.Error(), "deviceID", req.DeviceID, "type", req.Type)
		writeResponse(c, code.ErrInternalServer, nil)
		return
	}

	writeResponse(c, code.ErrSuccess, vo.LoginResponse{
		UserID:    userAccount.UserID,
		Token:     token,
		ExpiresIn: expiresIn.UnixMicro(),
		Nickname:  userAccount.User.Nickname,
		AvatarURL: userAccount.User.AvatarURL,
	})
}

func (s *Server) postLoginIDToken(c *gin.Context) {
	// TODO: Cache

	var req vo.LoginIDTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Errorw("login idtoken json failed: %s", err.Error(), "idtoken", req.IDToken)
		writeResponse(c, code.ErrInvalidParams, nil)
		return
	}

	claims, provider, err := s.verifier.Verify(c, req.IDToken, req.Nonce)
	if err != nil {
		log.Errorw("login idtoken verify failed: %s", err.Error(), "idtoken", req.IDToken)
		writeResponse(c, code.ErrSignatureInvalid, nil)
		return
	}

	userAccount, isNew, err := store.Client().UserAccount().FindOrCreate(claims.GetSubject(), model.AccountType(provider.GetName()))
	if err != nil {
		log.Errorw("login idtoken find or create failed: %s", err.Error(), "subject", claims.GetSubject(), "provider", provider.GetName())
		writeResponse(c, code.ErrInternalServer, nil)
		return
	}

	if isNew {
		err = store.Client().UserAccount().UpdateMeta(userAccount.ID, map[string]any{
			"email":             claims.GetEmail(),
			"is_email_verified": claims.GetEmailVerified(),
		})
		if err != nil {
			log.Errorw("login idtoken update meta failed: %s", err.Error(), "subject", claims.GetSubject(), "provider", provider.GetName())
		}
	}

	token, expiresIn, err := s.jwt.GenerateJWT(strconv.FormatUint(uint64(userAccount.UserID), 10), "user")
	if err != nil {
		log.Errorf("login idtoken generate jwt failed: %s", err.Error())
		writeResponse(c, code.ErrInternalServer, nil)
		return
	}

	writeResponse(c, code.ErrSuccess, vo.LoginResponse{
		UserID:    userAccount.UserID,
		Token:     token,
		ExpiresIn: expiresIn.UnixMicro(),
		Nickname:  userAccount.User.Nickname,
		AvatarURL: userAccount.User.AvatarURL,
	})
}

func (s *Server) getProfile(c *gin.Context) {
	userID := getUserID(c)

	user, err := store.Client().User().GetByID(userID)
	if err != nil {
		log.Errorw("getProfile: GetByID", "err", err, "userID", userID)
		writeResponse(c, code.ErrInternalServer, nil)
		return
	}

	writeResponse(c, code.ErrSuccess, user)
}

func (s *Server) postDelete(c *gin.Context) {
	userID := getUserID(c)

	if err := store.Client().User().RequestDeletion(userID); err != nil {
		log.Errorw("postDelete: RequestDeletion", "err", err, "userID", userID)

		if err.Error() == "duplicated request" {
			writeResponse(c, code.ErrDuplicatedRequest, nil)
			return
		}

		writeResponse(c, code.ErrInternalServer, nil)
		return
	}

	writeResponse(c, code.ErrSuccess, nil)
}
