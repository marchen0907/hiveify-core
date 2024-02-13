package database

import (
	"database/sql"
	"github.com/charmbracelet/log"
	psql "hiveify-core/database/postgresql"
)

// User 用户对象
type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

// TableUser 数据库用户对象
type TableUser struct {
	Name  sql.NullString
	Email sql.NullString
	Phone sql.NullString
}

func convertTableUserToUser(tableUser TableUser) User {
	user := User{
		Name:  "",
		Email: "",
		Phone: "",
	}
	if tableUser.Name.Valid {
		user.Name = tableUser.Name.String
	}
	if tableUser.Email.Valid {
		user.Email = tableUser.Email.String
	}
	if tableUser.Phone.Valid {
		user.Phone = tableUser.Phone.String
	}
	return user
}

// Register 注册用户
func Register(name, password, email string) bool {
	sqlStr := "INSERT INTO t_user (name, password, email) VALUES ($1, $2, $3)"
	commandTag, err := psql.DBConnection().Exec(psql.Ctx, sqlStr, name, password, email)
	if err != nil {
		log.Errorf("Failed to insert data, %s", err.Error())
		return false
	}
	return commandTag.Insert()
}

// Login 用户登录
func Login(email, password string) (User, error) {
	return GetUser(email, password)
}

// GetUser 获取用户信息
func GetUser(email, password string) (User, error) {
	var tableUser TableUser
	if len(password) > 0 {
		err := psql.DBConnection().QueryRow(psql.Ctx, "SELECT name, email, phone FROM t_user WHERE email = $1 "+
			"AND password = $2 AND status = 0", email, password).Scan(&tableUser.Name,
			&tableUser.Email, &tableUser.Phone)
		if err != nil {
			log.Infof("Failed to get user with email = %s and password = %s, %s", email, password, err.Error())
			return User{}, err
		}
	} else {
		err := psql.DBConnection().QueryRow(psql.Ctx, "SELECT name, email, phone FROM t_user WHERE email = $1 "+
			"AND status = 0", email).Scan(&tableUser.Name, &tableUser.Email, &tableUser.Phone)
		if err != nil {
			log.Infof("Failed to get user with email = %s, %s", email, err.Error())
			return User{}, err
		}
	}
	return convertTableUserToUser(tableUser), nil
}
