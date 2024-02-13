package interceptor

import (
	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	"hiveify-core/response"
	"hiveify-core/util"
	"net/http"
)

// HTTPInterceptor token验证的中间件
func HTTPInterceptor() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("Authorization")
		if err != nil {
			log.Warnf("interceptor token获取失败：%s", err.Error())
			c.Abort()
			response.Fail(c, http.StatusUnauthorized, "未登录，请先登录！")
			return
		}
		if !util.ValidateToken(token) {
			log.Warnf("token验证失败：%s", err.Error())
			c.Abort()
			response.Fail(c, http.StatusUnauthorized, "登陆已过期，请重新登录！")
			return
		}
		c.Next()
	}
}
