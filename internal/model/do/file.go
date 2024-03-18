// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// File is the golang structure of table t_file for DAO operations like Where/Data.
type File struct {
	g.Meta     `orm:"table:t_file, do:true"`
	Id         interface{} // 主键，自增长
	Sha512     interface{} // 文件sha512
	Extension  interface{} // 文件扩展名
	Size       interface{} // 文件大小
	Category   interface{} // 文件分类
	CreateTime *gtime.Time // 文件创建时间
	UpdateTime *gtime.Time // 文件更新时间
	Status     interface{} // 是否删除，true为没有删除
}
