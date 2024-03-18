package middleware

import (
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
	"hiveify-core/utility"
	"net/http"
)

// Auth 用户认证中间件
func Auth(req *ghttp.Request) {
	URI := req.URL.Path
	if URI == "/token" || (URI == "/user" && req.Method == "POST") {
		req.Middleware.Next()
	} else {
		ctx := gctx.New()
		token := req.Cookie.Get("Authorization").String()
		if len(token) == 0 || !utility.ValidateToken(ctx, token) {
			req.Response.WriteJson(struct {
				Code    int         `json:"code"`
				Message string      `json:"message"`
				Data    interface{} `json:"data"`
			}{
				Code:    http.StatusUnauthorized,
				Message: "未登录或登录已过期，请重新登录！",
				Data:    nil,
			})
		} else {
			req.Middleware.Next()
		}
	}

}
