package res

import (
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
func result(code int, data any, msg string) *Response {
	timer := strconv.FormatInt(time.Now().UnixMicro(), 10)
	return &Response{
		Code:    code,
		Data:    data,
		Message: msg,
		Timer:   timer,
	}
}

// SuccessOfMessage 返回成功后的提示消息
func SuccessOfMessage(msg string) *Response {
	return result(success, nil, msg)
}

// SuccessOfData 返回成功后的数据
func SuccessOfData(data any) *Response {
	return result(success, data, "")
}

// FailOfMessage 返回失败后的提示消息
func FailOfMessage(msg string) *Response {
	return result(err, nil, msg)
}

// FailOfCode 返回失败后的状态码
func FailOfCode(code int) *Response {
	var msg string
	switch code {
	case http.StatusUnauthorized:
		msg = "登陆失效"
	case http.StatusForbidden:
		msg = "无权限"
	}
	return result(code, nil, msg)
}

// FailOfData 返回失败后的数据
func FailOfData(data any) *Response {
	return result(err, data, "")
}

// Res 返回结果
func Res(code int, data any, msg string) *Response {
	return result(code, data, msg)
}
