package res

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

// Response 返回的数据结构体
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Timer   string      `json:"timer"`
}

// 返回成功或失败的默认code
const (
	success = http.StatusOK
	err     = http.StatusBadRequest
)

// result 返回结果
func result(c *gin.Context, code int, data any, msg string) {
	timer := strconv.FormatInt(time.Now().UnixMicro(), 10)
	c.JSON(success, Response{
		Code:    code,
		Data:    data,
		Message: msg,
		Timer:   timer,
	})
}

// SuccessOfMessage 返回成功后的提示消息
func SuccessOfMessage(ctx *gin.Context, msg string) {
	result(ctx, success, nil, msg)
}

// SuccessOfData 返回成功后的数据
func SuccessOfData(ctx *gin.Context, data any) {
	result(ctx, success, data, "")
}

// FailOfMessage 返回失败后的提示消息
func FailOfMessage(ctx *gin.Context, msg string) {
	result(ctx, err, nil, msg)
}

func FailOfCode(ctx *gin.Context, code int) {
	var msg string
	switch code {
	case http.StatusUnauthorized:
		msg = "登陆失效"
	case http.StatusForbidden:
		msg = "无权限"
	}
	result(ctx, code, nil, msg)
}

// FailOfData 返回失败后的数据
func FailOfData(ctx *gin.Context, data any) {
	result(ctx, err, data, "")
}

func Res(c *gin.Context, code int, data any, msg string) {
	result(c, code, data, msg)
}
