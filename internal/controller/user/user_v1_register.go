package user

import (
	"context"
	"hiveify-core/internal/model"
	"hiveify-core/internal/service"

	"hiveify-core/api/user/v1"
)

// Register 用户注册
func (c *ControllerV1) Register(ctx context.Context, req *v1.RegisterReq) (res *v1.RegisterRes, err error) {
	err = service.User().Register(ctx, model.UserRegisterInput{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	})
	return
}
