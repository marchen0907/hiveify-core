package database

import (
	"database/sql"
	"errors"
	"github.com/charmbracelet/log"
	"github.com/jackc/pgx/v5"
	psql "hiveify-core/database/postgresql"
	"time"
)

// TableFile 文件表结构体
type TableFile struct {
	Sha512     string
	Name       sql.NullString
	Size       sql.NullInt64
	Path       sql.NullString
	CreateTime sql.NullTime
}

// File 文件上传信息
type File struct {
	Sha512     string    `json:"sha512"`
	Name       string    `json:"name"`
	Size       int64     `json:"size"`
	Path       string    `json:"path"`
	CreateTime time.Time `json:"createTime"`
}

// convertTableFileToFile 将TableFile转换为File
func convertTableFileToFile(tableFile TableFile) File {
	file := File{
		Sha512:     tableFile.Sha512,
		Name:       "",
		Size:       0,
		Path:       "",
		CreateTime: time.Time{},
	}

	// 检查是否为 NULL，如果不是，则赋值给对应的字段
	if tableFile.Name.Valid {
		file.Name = tableFile.Name.String
	}

	if tableFile.Size.Valid {
		file.Size = tableFile.Size.Int64
	}

	if tableFile.Path.Valid {
		file.Path = tableFile.Path.String
	}

	if tableFile.CreateTime.Valid {
		file.CreateTime = tableFile.CreateTime.Time
	}

	return file
}

// CreateFile 在t_file中写入文件信息
func CreateFile(sha512, name, path string, size int64) bool {
	conn, err := psql.DBConnection().Acquire(psql.Ctx)
	if err != nil {
		log.Errorf("Failed to get a connection from the pool, %s", err.Error())
		return false
	}
	defer conn.Release()
	sqlStr := "INSERT INTO t_file(sha512, name, size, path) VALUES ($1,$2,$3,$4)"
	commandTag, err := conn.Exec(psql.Ctx, sqlStr, sha512, name, size, path)
	if err != nil {
		log.Errorf("Failed to insert data, %s", err.Error())
		return false
	}
	return commandTag.Insert()
}

// GetFile 从数据库中获取文件信息
func GetFile(sha512 string) (File, error) {
	sqlStr := "SELECT sha512,name,size,path,create_time FROM t_file WHERE sha512 = $1 AND status IS true"
	var tableFile TableFile
	err := psql.DBConnection().QueryRow(psql.Ctx, sqlStr, sha512).Scan(&tableFile.Sha512, &tableFile.Name,
		&tableFile.Size, &tableFile.Path, &tableFile.CreateTime)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// 没有找到匹配的行
			log.Warnf("No file found with sha512 %s", sha512)
			return File{}, err
		}
		log.Warnf("Failed to retrieve file with sha512 %s: %s", sha512, err.Error())
		return File{}, err
	}
	return convertTableFileToFile(tableFile), nil
}

// DeleteFile 从数据库中删除文件信息
func DeleteFile(sha512 string) bool {
	sqlStr := "UPDATE t_file SET status = false WHERE sha512 = $1 AND status IS true"
	commandTag, err := psql.DBConnection().Exec(psql.Ctx, sqlStr, sha512)
	if err != nil {
		log.Errorf("Failed to delete file, %s", err.Error())
		return false
	}
	return commandTag.Update()
}
