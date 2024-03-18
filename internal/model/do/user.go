// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// User is the golang structure of table t_user for DAO operations like Where/Data.
type User struct {
	g.Meta         `orm:"table:t_user, do:true"`
	Id             interface{} //
	Name           interface{} // 用户名
	Password       interface{} // 密码
	Email          interface{} // 邮箱
	Phone          interface{} // 手机
	EmailValidated interface{} // 邮箱是否验证
	PhoneValidated interface{} // 手机是否验证
	CreateTime     *gtime.Time // 创建时间
	UpdateTime     *gtime.Time // 上次更新时间
	Profile        interface{} // 用户属性
	Status         interface{} // 账户状态
}
