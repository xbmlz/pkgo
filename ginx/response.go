package ginx

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	SuccessCode = 0
	ErrorCode   = 1
)

type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

func ResponseOk(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Response{
		Code: SuccessCode,
		Msg:  "success",
		Data: data,
	})
}

func ResponseError(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, Response{
		Code: ErrorCode,
		Msg:  msg,
		Data: nil,
	})
}

func ResponseCustom(c *gin.Context, code int, msg string, data any) {
	c.JSON(http.StatusOK, Response{
		Code: code,
		Msg:  msg,
		Data: data,
	})
}
