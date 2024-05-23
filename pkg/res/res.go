package res

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code int    `json:"code"`
	Data any    `json:"data"`
	Msg  string `json:"msg"`
}

const (
	Success = 0
	Failed  = 7
)

// 请求访问成功的响应方法
func OK(data any, msg string, ctx *gin.Context) {
	ctx.JSON(http.StatusOK, Response{Code: Success, Data: data, Msg: msg})
}

func OKWithMsg(msg string, ctx *gin.Context) {
	OK(nil, msg, ctx)
}

func OKWithData(data any, ctx *gin.Context) {
	OK(data, "请求成功", ctx)
}

// 请求访问失败的响应方法
func Fail(code int, data any, msg string, ctx *gin.Context) {    // 错误码不止一种
	ctx.JSON(http.StatusOK, Response{Code: code, Data: data, Msg: msg})
}

func FailWithMsg(msg string, ctx *gin.Context) {
	Fail(Failed, nil, msg, ctx)
}

func FailWithData(data any, ctx *gin.Context) {
	Fail(Failed, data, "系统错误", ctx)
}