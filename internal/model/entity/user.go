// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// User is the golang structure for table user.
type User struct {
	Id             int         `json:"id"             ` //
	Name           string      `json:"name"           ` // 用户名
	Password       string      `json:"password"       ` // 密码
	Email          string      `json:"email"          ` // 邮箱
	Phone          string      `json:"phone"          ` // 手机
	EmailValidated bool        `json:"emailValidated" ` // 邮箱是否验证
	PhoneValidated bool        `json:"phoneValidated" ` // 手机是否验证
	CreateTime     *gtime.Time `json:"createTime"     ` // 创建时间
	UpdateTime     *gtime.Time `json:"updateTime"     ` // 上次更新时间
	Profile        string      `json:"profile"        ` // 用户属性
	Status         int         `json:"status"         ` // 账户状态
}
