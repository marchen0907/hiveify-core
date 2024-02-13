package database

import (
	"errors"
	"github.com/charmbracelet/log"
	"github.com/jackc/pgx/v5"
	psql "hiveify-core/database/postgresql"
	"time"
)

// UserFile 用户文件表结构体
type UserFile struct {
	Email      string    `json:"email"`
	Sha512     string    `json:"sha512"`
	FileName   string    `json:"fileName"`
	CreateTime time.Time `json:"createTime"`
}

// CreateUserFile 创建用户文件
func CreateUserFile(email, sha512, fileName string) bool {
	sqlStr := "INSERT INTO t_user_file (email, sha512, file_name) VALUES ($1, $2, $3)"
	commandTag, err := psql.DBConnection().Exec(psql.Ctx, sqlStr,
		email, sha512, fileName)
	if err != nil {
		log.Warnf("Failed to create user file with eamil = %s AND sha512 = %s AND fileName = %s, %s",
			email, sha512, fileName, err.Error())
		return false
	}
	return commandTag.Insert()
}

// GetUserFile 获取用户文件
func GetUserFile(email, sha512 string) (UserFile, error) {
	// 实现获取用户文件的逻辑
	var userFile UserFile
	sqlStr := "SELECT email, sha512, file_name, create_time FROM t_user_file WHERE email = $1 AND sha512 = $2"
	err := psql.DBConnection().QueryRow(psql.Ctx, sqlStr,
		email, sha512).Scan(&userFile.Email, &userFile.Sha512, &userFile.FileName, &userFile.CreateTime)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// 没有找到匹配的行
			log.Warnf("No file found with sha512 %s", sha512)
			return UserFile{}, err
		}
		log.Warnf("Failed to retrieve user_file with sha512 %s: %s", sha512, err.Error())
		return UserFile{}, err
	}
	return userFile, nil
}

// GetAllUserFile 获取所有用户文件
func GetAllUserFile(email string) ([]File, error) {
	sqlStr := "SELECT t_user_file.sha512, t_user_file.file_name, t_file.size, t_file.path, t_user_file.create_time " +
		"from t_user_file, t_file WHERE t_user_file.sha512 = t_file.sha512 AND t_user_file.email = $1 " +
		"AND t_user_file.status IS true AND t_file.status IS true;"

	var allUserFile []File
	rows, err := psql.DBConnection().Query(psql.Ctx, sqlStr, email)
	if err != nil {
		log.Errorf("Failed to query t_user_file and t_file: %s", err.Error())
		return []File{}, err
	}
	for rows.Next() {
		var file File
		err := rows.Scan(&file.Sha512, &file.Name, &file.Size, &file.Path, &file.CreateTime)
		if err != nil {
			log.Warnf("Failed to scan row: %s", err.Error())
			return []File{}, err
		}

		// 将每一行的数据添加到切片中
		allUserFile = append(allUserFile, file)
	}

	// 检查是否有错误导致迭代结束
	if err := rows.Err(); err != nil {
		log.Warnf("Error while scanning rows: %s", err.Error())
		return []File{}, nil
	}
	return allUserFile, nil
}

// UpdateUserFile 更新用户文件
func UpdateUserFile(email, sha512, newFileName string) bool {
	commandTag, err := psql.DBConnection().Exec(psql.Ctx, "UPDATE t_user_file SET file_name = $1 WHERE email = $2 "+
		"AND sha512 = $3 AND status IS true", newFileName, email, sha512)
	if err != nil {
		log.Warnf("Failed to update t_user_file: %s", err.Error())
		return false
	}
	return commandTag.Update()
}

// DeleteUserFile 删除用户文件
func DeleteUserFile(email, sha512 string) bool {
	commandTag, err := psql.DBConnection().Exec(psql.Ctx, "UPDATE t_user_file SET status = false WHERE email = $1 "+
		"AND sha512 = $2 AND status IS true", email, sha512)
	if err != nil {
		log.Warnf("Failed to delete t_user_file: %s", err.Error())
		return false
	}
	return commandTag.Update()
}
