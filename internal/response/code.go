package response

import (
	"net/http"
)

var (
	// CodeOK 请求成功
	CodeOK = localCode{http.StatusOK, "OK", nil}
	// CodeNotAuthorized 未授权
	CodeNotAuthorized = localCode{http.StatusUnauthorized, "Not Authorized", nil}
	// CodeServerError 服务器内部错误
	CodeServerError = localCode{http.StatusInternalServerError, "Internal Error", nil}
	// CodeBadRequest 错误请求
	CodeBadRequest = localCode{http.StatusBadRequest, "Bad Request", nil}
)
