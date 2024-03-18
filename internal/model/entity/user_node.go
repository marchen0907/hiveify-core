// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// UserNode is the golang structure for table user_node.
type UserNode struct {
	Id         int         `json:"id"         ` // 主键
	Email      string      `json:"email"      ` // 邮箱与用户表关联
	Type       string      `json:"type"       ` // 节点类型
	Sha512     string      `json:"sha512"     ` // sha512
	ParentNode string      `json:"parentNode" ` // 父节点
	Name       string      `json:"name"       ` // 用户定义的文件（夹）名
	Sync       bool        `json:"sync"       ` // 是否自动同步
	CreateTime *gtime.Time `json:"createTime" ` // 创建时间
	UpdateTime *gtime.Time `json:"updateTime" ` // 更新时间
	Status     bool        `json:"status"     ` // 状态
}
