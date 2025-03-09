package serializer

import (
	"net/http"
)

type Response struct {
	Code  int         `json:"code"`
	Data  interface{} `json:"data,omitempty"`
	Msg   string      `json:"msg"`
	Error string      `json:"error,omitempty"`
}

func ErrorResponse(err error, msg ...string) (response Response) {
	defaultMsg := "操作失败"
	if len(msg) > 0 {
		defaultMsg = msg[0]
	}
	return Response{
		Code:  http.StatusInternalServerError,
		Data:  nil,
		Msg:   defaultMsg,
		Error: err.Error(),
	}
}

func SuccessResponse(data interface{}, msg ...string) Response {
	defaultMsg := "操作成功"
	if len(msg) > 0 {
		defaultMsg = msg[0]
	}
	return Response{
		Code:  http.StatusOK,
		Data:  data,
		Msg:   defaultMsg,
		Error: "",
	}
}
