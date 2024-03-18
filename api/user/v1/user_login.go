package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// LoginReq 用户登录请求参数
type LoginReq struct {
	g.Meta   `path:"/token" method:"get" tags:"user" summary:"用户登录"`
	Email    string `name:"email" in:"query"  v:"required|email#邮箱不能为空|请输入邮箱" dc:"用户邮箱"`
	Password string `name:"password" in:"query" v:"required|length:8,20#密码不能为空|密码长度为8到20位" dc:"用户密码"`
}

// LoginRes 用户登录返回参数
type LoginRes struct {
	Token string `json:"token"`
}
