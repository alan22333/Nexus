package res

import (
	"Nuxus/pkg/erru"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

func response(ctx *gin.Context, code int, data any, msg string) {
	ctx.JSON(200, Response{
		Code: code,
		Msg:  msg,
		Data: data,
	})
}

func Ok(ctx *gin.Context, data any, msg string) {
	response(ctx, 0, data, msg)
}

func OkWithData(ctx *gin.Context, data any) {
	Ok(ctx, data, "success")
}

func OkWithMsg(ctx *gin.Context, msg string) {
	Ok(ctx, gin.H{}, msg)
}

func Fail(ctx *gin.Context, code int, data any, msg string) {
	response(ctx, code, data, msg)
}

func FailWithAppErr(ctx *gin.Context, err *erru.AppError) {
	response(ctx, err.Code, err.Err, err.Msg)
}

func FailWithMsg(ctx *gin.Context, msg string) {
	Fail(ctx, 1003, nil, msg)
}

// TODO:
func FailWithCode(ctx *gin.Context, code int) {

	Fail(ctx, code, nil, "")
}
