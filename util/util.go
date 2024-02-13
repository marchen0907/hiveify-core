package util

import (
	"crypto/sha512"
	"encoding/hex"
	"github.com/charmbracelet/log"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"regexp"
	"sort"
)

// FileInfo 计算sha512返回的文件信息
type FileInfo struct {
	Sha512 string
	Size   int64
}

// FileSha512 计算文件sha512
func FileSha512(fileHeader *multipart.FileHeader) (FileInfo, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return FileInfo{}, err
	}
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			log.Warnf("Failed to close file: %s", err.Error())
		}
	}(file)

	hash := sha512.New()

	if _, err := io.Copy(hash, file); err != nil {
		return FileInfo{}, err
	}

	hashInBytes := hash.Sum(nil)
	hashString := hex.EncodeToString(hashInBytes)

	// 返回字符串形式的SHA-512校验和
	return FileInfo{
		Sha512: hashString,
		Size:   fileHeader.Size,
	}, nil
}

// StringSha512 计算字符串的sha512
func StringSha512(str string) string {
	hash := sha512.New()
	hash.Write([]byte(str))
	return hex.EncodeToString(hash.Sum(nil))
}

// CheckEmail 检查email格式
func CheckEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// CheckFile 检查文件是否被篡改
func CheckFile(fileHeader *multipart.FileHeader, sha512 string) bool {
	fileInfo, err := FileSha512(fileHeader)
	if err != nil {
		log.Warnf("Failed to calculate file SHA-512: %s", err.Error())
		return false
	}
	return fileInfo.Sha512 == sha512
}

// SaveToLocal 保存文件到本地
func SaveToLocal(fileHeader *multipart.FileHeader, fileName, path string) bool {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		log.Warnf("Failed to create directory %s, %s", path, err.Error())
		return false
	}
	newFile, createError := os.Create(path)
	if createError != nil {
		log.Warnf("Failed to create %s, %s", fileName, createError.Error())
		return false
	}
	defer func(newFile *os.File) {
		err := newFile.Close()
		if err != nil {
			log.Warnf("Failed to close %s, %s", newFile.Name(), err.Error())
			return
		}
	}(newFile)
	file, err := fileHeader.Open()
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			log.Warnf("Failed to close %s, %s", fileHeader.Filename, err.Error())
			return
		}
	}(file)
	_, err = io.Copy(newFile, file)
	if err != nil {
		log.Warnf("Failed to save %s, %s", newFile.Name(), err.Error())
		return false
	}
	return true
}

// getChunkFiles 获取所有分块文件
func getChunkFiles(directoryPath string) ([]string, error) {
	var files []string

	// 获取目录中的文件列表
	err := filepath.Walk(directoryPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// 只处理普通文件
		if info.Mode().IsRegular() {
			files = append(files, info.Name())
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}

// MergeChunk 合并分块文件
func MergeChunk(outputPath, directoryPath string) bool {
	files, err := getChunkFiles(directoryPath)
	if err != nil {
		return false
	}

	// 按照文件名（分块编号）排序
	sort.Strings(files)

	// 创建合并后的文件
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return false
	}
	defer func() {
		err := outputFile.Close()
		if err != nil {
			log.Warnf("Failed to close %s, %s", outputFile.Name(), err.Error())
		}
	}()

	// 逐个合并分块
	for _, file := range files {
		chunkPath := filepath.Join(directoryPath, file)

		// 打开分块文件
		chunkFile, err := os.Open(chunkPath)
		if err != nil {
			return false
		}

		// 将分块内容写入合并后的文件
		_, err = io.Copy(outputFile, chunkFile)
		if err != nil {
			err := chunkFile.Close()
			if err != nil {
				return false
			}
			return false
		}

		// 关闭分块文件
		err = chunkFile.Close()
		if err != nil {
			log.Warnf("Failed to close %s, %s", chunkFile.Name(), err.Error())
		}
	}
	err = os.RemoveAll(directoryPath)
	if err != nil {
		return false
	}
	return true
}
