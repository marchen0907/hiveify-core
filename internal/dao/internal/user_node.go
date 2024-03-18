// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// UserNodeDao is the data access object for table t_user_node.
type UserNodeDao struct {
	table   string          // table is the underlying table name of the DAO.
	group   string          // group is the database configuration group name of current DAO.
	columns UserNodeColumns // columns contains all the column names of Table for convenient usage.
}

// UserNodeColumns defines and stores column names for table t_user_node.
type UserNodeColumns struct {
	Id         string // 主键
	Email      string // 邮箱与用户表关联
	Type       string // 节点类型
	Sha512     string // sha512
	ParentNode string // 父节点
	Name       string // 用户定义的文件（夹）名
	Sync       string // 是否自动同步
	CreateTime string // 创建时间
	UpdateTime string // 更新时间
	Status     string // 状态
}

// userNodeColumns holds the columns for table t_user_node.
var userNodeColumns = UserNodeColumns{
	Id:         "id",
	Email:      "email",
	Type:       "type",
	Sha512:     "sha512",
	ParentNode: "parent_node",
	Name:       "name",
	Sync:       "sync",
	CreateTime: "create_time",
	UpdateTime: "update_time",
	Status:     "status",
}

// NewUserNodeDao creates and returns a new DAO object for table data access.
func NewUserNodeDao() *UserNodeDao {
	return &UserNodeDao{
		group:   "default",
		table:   "t_user_node",
		columns: userNodeColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *UserNodeDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *UserNodeDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *UserNodeDao) Columns() UserNodeColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *UserNodeDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *UserNodeDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *UserNodeDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
