package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// RegisterReq 用户注册请求参数
type RegisterReq struct {
	g.Meta   `path:"/user" method:"post" tags:"user" summary:"用户注册"`
	Name     string `json:"name" v:"required#用户名不能为空" dc:"用户名"`
	Email    string `json:"email" v:"required|email#用户名不能为空|邮箱格式不正确" dc:"用户邮箱"`
	Password string `json:"password" v:"required|length:8,20#密码不能为空|密码长度为8到20位" dc:"用户密码"`
}

// RegisterRes 用户注册返回参数
type RegisterRes struct{}
