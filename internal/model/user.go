package model

// UserLoginInput 用户登录输入参数
type UserLoginInput struct {
	Email    string
	Password string
}

// UserRegisterInput 用户注册输入参数
type UserRegisterInput struct {
	Name     string
	Email    string
	Password string
}
