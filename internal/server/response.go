package server

import (
	"github.com/chunkgar/gin-template/internal/code"
	"github.com/gin-gonic/gin"
)

type Response struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
	Data  any    `json:"data"`
}

func writeResponse(ctx *gin.Context, c int, data any, msg ...string) {
	var msgStr string
	cc := code.Use(c)

	if len(msg) > 0 {
		msgStr = msg[0]
	} else {
		msgStr = cc.String()
	}

	ctx.JSON(cc.HTTPStatus(), Response{
		Code:  cc.Code(),
		Error: msgStr,
		Data:  data,
	})
}
