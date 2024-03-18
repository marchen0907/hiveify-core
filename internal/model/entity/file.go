// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// File is the golang structure for table file.
type File struct {
	Id         int         `json:"id"         ` // 主键，自增长
	Sha512     string      `json:"sha512"     ` // 文件sha512
	Extension  string      `json:"extension"  ` // 文件扩展名
	Size       int64       `json:"size"       ` // 文件大小
	Category   string      `json:"category"   ` // 文件分类
	CreateTime *gtime.Time `json:"createTime" ` // 文件创建时间
	UpdateTime *gtime.Time `json:"updateTime" ` // 文件更新时间
	Status     bool        `json:"status"     ` // 是否删除，true为没有删除
}
