// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
	"hiveify-core/internal/model"
	"hiveify-core/internal/model/entity"
)

type (
	IUser interface {
		IsEmailAvailable(ctx context.Context, email string) (bool, error)
		Register(ctx context.Context, in model.UserRegisterInput) (err error)
		Login(ctx context.Context, in model.UserLoginInput) (string, error)
		GetUser(ctx context.Context) (*entity.User, error)
	}
)

var (
	localUser IUser
)

func User() IUser {
	if localUser == nil {
		panic("implement not found for interface IUser, forgot register?")
	}
	return localUser
}

func RegisterUser(i IUser) {
	localUser = i
}
