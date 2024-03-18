package user

import (
	"context"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"hiveify-core/internal/consts"
	"hiveify-core/internal/dao"
	"hiveify-core/internal/model"
	"hiveify-core/internal/model/do"
	"hiveify-core/internal/model/entity"
	"hiveify-core/internal/response"
	"hiveify-core/internal/service"
	"hiveify-core/utility"
	"net/http"
)

type (
	sUser struct{}
)

// init 注册用户服务
func init() {
	service.RegisterUser(New())
}

// New 实例化用户服务
func New() service.IUser {
	return &sUser{}
}

// IsEmailAvailable 检查此邮箱是否可用
func (s *sUser) IsEmailAvailable(ctx context.Context, email string) (bool, error) {
	count, err := dao.User.Ctx(ctx).Where(do.User{
		Email: email,
	}).Count()
	if err != nil {
		return false, err
	}
	return count == 0, nil
}

// Register 用户注册
func (s *sUser) Register(ctx context.Context, in model.UserRegisterInput) (err error) {
	ip := g.RequestFromCtx(ctx).GetClientIp()
	count, err := g.Redis().Get(ctx, ip)
	if err != nil {
		_, err2 := g.Redis().Do(ctx, "SET", in.Email+"_token", 0, "EX", 24*60*60)
		if err2 != nil {
			return gerror.NewCodef(response.CodeServerError, "服务器出错，注册失败!")
		}
	} else {
		if count.Int() > 5 {
			return gerror.NewCodef(gcode.New(http.StatusTooManyRequests, "", nil), "请求过于频繁，请稍后再试!")
		}
		_, err2 := g.Redis().Do(ctx, "INCR", ip)
		if err2 != nil {
			return gerror.NewCodef(response.CodeServerError, "服务器出错，注册失败!")
		}
	}
	emailAvailable, err := s.IsEmailAvailable(ctx, in.Email)
	if err != nil {
		return gerror.NewCodef(response.CodeServerError, "服务器出错，注册失败!")
	}
	if !emailAvailable {
		return gerror.NewCodef(response.CodeBadRequest, "邮箱已被注册!")
	}
	_, err = dao.User.Ctx(ctx).Data(do.User{
		Name:     in.Name,
		Email:    in.Email,
		Password: utility.StringSha512(in.Password + consts.PasswordSalt),
	}).Insert()
	if err != nil {
		return gerror.NewCodef(response.CodeServerError, "服务器出错，注册失败!")
	}
	return gerror.NewCodef(response.CodeOK, "注册成功!")
}

// Login 用户登录
func (s *sUser) Login(ctx context.Context, in model.UserLoginInput) (string, error) {
	var user *entity.User
	err := dao.User.Ctx(ctx).Fields("name, email, phone").Where(do.User{
		Email:    in.Email,
		Password: utility.StringSha512(in.Password + consts.PasswordSalt),
		Status:   0,
	}).Scan(&user)
	if err != nil {
		return "", gerror.NewCodef(response.CodeServerError, "登陆失败，服务器出错!")
	}
	if user == nil {
		return "", gerror.NewCodef(response.CodeNotAuthorized, "登陆失败，用户名或密码错误!")
	}

	token, err := utility.GetToken(user)
	if err != nil {
		return "", err
	}

	_, err = g.Redis().Do(ctx, "SET", in.Email+"_token", token, "EX", 24*60*60)
	if err != nil {
		g.Log().Errorf(ctx, "Error setting token: %s", err.Error())
		return "", gerror.NewCodef(response.CodeServerError, "登陆失败，服务器出错!")
	}
	return token, gerror.NewCodef(response.CodeOK, "登陆成功!")
}

// GetUser 获取用户信息
func (s *sUser) GetUser(ctx context.Context) (*entity.User, error) {
	token := g.RequestFromCtx(ctx).Cookie.Get("Authorization").String()
	user, err := utility.ParseToken(ctx, token)
	if err != nil {
		return nil, err
	}
	return user, gerror.NewCodef(response.CodeOK, "获取用户信息成功!")
}
