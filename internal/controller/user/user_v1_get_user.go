package user

import (
	"context"
	"hiveify-core/internal/service"

	"hiveify-core/api/user/v1"
)

// GetUser 获取用户信息
func (c *ControllerV1) GetUser(ctx context.Context, _ *v1.GetUserReq) (res *v1.GetUserRes, err error) {
	user, err := service.User().GetUser(ctx)
	res = &v1.GetUserRes{
		User: user,
	}
	return
}
