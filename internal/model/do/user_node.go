// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// UserNode is the golang structure of table t_user_node for DAO operations like Where/Data.
type UserNode struct {
	g.Meta     `orm:"table:t_user_node, do:true"`
	Id         interface{} // 主键
	Email      interface{} // 邮箱与用户表关联
	Type       interface{} // 节点类型
	Sha512     interface{} // sha512
	ParentNode interface{} // 父节点
	Name       interface{} // 用户定义的文件（夹）名
	Sync       interface{} // 是否自动同步
	CreateTime *gtime.Time // 创建时间
	UpdateTime *gtime.Time // 更新时间
	Status     interface{} // 状态
}
