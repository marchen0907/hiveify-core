package user

import (
	"context"
	"hiveify-core/internal/model"
	"hiveify-core/internal/service"

	"hiveify-core/api/user/v1"
)

// Login 用户登录
func (c *ControllerV1) Login(ctx context.Context, req *v1.LoginReq) (res *v1.LoginRes, err error) {
	token, err := service.User().Login(ctx, model.UserLoginInput{
		Email:    req.Email,
		Password: req.Password,
	})
	res = &v1.LoginRes{
		Token: token,
	}
	return
}
