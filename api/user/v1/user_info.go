package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"hiveify-core/internal/model/entity"
)

// GetUserReq 获取用户信息请求
type GetUserReq struct {
	g.Meta `path:"/user" method:"get" tags:"user" summary:"获取用户信息"`
}

// GetUserRes 获取用户信息响应
type GetUserRes struct {
	*entity.User
}
