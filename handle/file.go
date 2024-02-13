package handle

import (
	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	"hiveify-core/database"
	"hiveify-core/response"
	"hiveify-core/storage/tencentCos"
	"hiveify-core/storage/tencentSts"
	"hiveify-core/util"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// UploadHandle 文件上传，文件不在t_file中
func UploadHandle(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		log.Infof("Failed to get file, %s", err.Error())
		response.Fail(c, http.StatusInternalServerError, "上传失败！")
		return
	}

	fileMeta := database.File{
		Name:       fileHeader.Filename,
		CreateTime: time.Now(),
	}

	//  计算文件哈希值
	fileInfo, err := util.FileSha512(fileHeader)
	if err != nil {
		log.Warnf("Failed to get file sha512 %s", err)
		response.Fail(c, http.StatusInternalServerError, "上传失败！")
		return
	}
	fileMeta.Sha512 = fileInfo.Sha512
	fileMeta.Size = fileInfo.Size
	fileMeta.Path = fileInfo.Sha512 + filepath.Ext(fileMeta.Name)

	// 先保存文件到存储
	if ok := tencentCos.PutFile(fileMeta.Path, fileHeader); !ok {
		response.Fail(c, http.StatusInternalServerError, "上传失败！")
		return
	}

	//  保存文件信息到t_file
	if dbResult := database.CreateFile(fileMeta.Sha512, fileMeta.Name, fileMeta.Path, fileMeta.Size); !dbResult {
		response.Fail(c, http.StatusInternalServerError, "上传失败！")
		return
	}

	// 将文件信息保存到t_user_file
	user, _ := util.GetUserByRequest(c)
	if userFileDB := database.CreateUserFile(user.Email, fileMeta.Sha512, fileMeta.Name); !userFileDB {
		response.Fail(c, http.StatusInternalServerError, "上传失败！")
		return
	}
	response.Success(c, "文件上传成功！",
		database.File{
			Sha512:     fileMeta.Sha512,
			Name:       fileMeta.Name,
			Size:       fileMeta.Size,
			Path:       "https://cos.wujunyi.cn/" + fileMeta.Path,
			CreateTime: fileMeta.CreateTime,
		})
}

// CompleteCOSHandle 通知文件已经上传到cos
func CompleteCOSHandle(c *gin.Context) {
	sha512 := c.Request.FormValue("sha512")
	name := c.Request.FormValue("name")
	path := c.Request.FormValue("path")
	size, err := strconv.ParseInt(c.Request.FormValue("size"), 10, 64)
	if err != nil {
		log.Warnf("参数非法：%s", err.Error())
		response.Fail(c, http.StatusBadRequest, "参数非法！")
		return
	}

	user, _ := util.GetUserByRequest(c)

	if !tencentCos.CheckFileExist(path) ||
		!database.CreateFile(sha512, name, path, size) ||
		!database.CreateUserFile(user.Email, sha512, name) {
		response.Fail(c, http.StatusBadRequest, "文件上传失败！")
		return
	}

	response.Success(c, "文件上传成功！", database.File{
		Sha512:     sha512,
		Name:       name,
		Size:       size,
		Path:       path,
		CreateTime: time.Now(),
	})
}

// FastUploadHandle 秒传接口，文件已经存在于t_file表中
func FastUploadHandle(c *gin.Context) {
	sha512 := c.Request.FormValue("sha512")
	fileName := c.Request.FormValue("fileName")
	user, _ := util.GetUserByRequest(c)

	file, err := database.GetFile(sha512)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "文件不存在！")
		return
	}
	file.Name = fileName

	if ok := database.CreateUserFile(user.Email, file.Sha512, file.Name); !ok {
		response.Fail(c, http.StatusInternalServerError, "文件上传失败！")
		return
	}

	response.Success(c, "文件上传成功！", file)
}

// GetFileHandle 获取文件信息
func GetFileHandle(c *gin.Context) {
	user, _ := util.GetUserByRequest(c)
	userFiles, err := database.GetAllUserFile(user.Email)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "获取文件信息失败！")
		return
	}
	response.Success(c, "获取文件信息成功！", userFiles)
}

// DownloadFileHandle 下载文件
func DownloadFileHandle(c *gin.Context) {
	sha512 := c.Request.FormValue("sha512")
	user, _ := util.GetUserByRequest(c)
	_, err := database.GetUserFile(user.Email, sha512)
	if err != nil {
		response.Fail(c, http.StatusNotFound, "文件不存在！")
		return
	}
	file, err := database.GetFile(sha512)
	if err != nil {
		response.Fail(c, http.StatusNotFound, "文件不存在！")
		return
	}
	file.Path, err = tencentCos.GetDownloadURL(file.Path)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "获取文件下载链接失败！")
		return
	}
	response.Success(c, "查询成功", file)
}

// UpdateFileHandle 文件重命名
func UpdateFileHandle(c *gin.Context) {
	sha512 := c.Request.FormValue("sha512")
	newName := c.Request.FormValue("newName")
	user, _ := util.GetUserByRequest(c)
	file, err := database.GetUserFile(user.Email, sha512)
	if err != nil {
		response.Fail(c, http.StatusNotFound, "文件不存在！")
		return
	}
	file.FileName = newName
	dbResult := database.UpdateUserFile(user.Email, file.Sha512, file.FileName)
	file, err = database.GetUserFile(user.Email, sha512)
	if err != nil && !dbResult {
		response.Fail(c, http.StatusInternalServerError, "修改失败！")
		return
	}
	response.Success(c, "修改成功", file)

}

// DeleteFileHandle 删除文件
func DeleteFileHandle(c *gin.Context) {
	sha512 := c.Request.FormValue("sha512")
	file, err := database.GetFile(sha512)
	if err != nil {
		response.Fail(c, http.StatusNotFound, "文件不存在！")
		return
	}
	err = os.Remove(file.Path)
	removeFile := database.DeleteFile(sha512)
	if err != nil && !removeFile {
		response.Fail(c, http.StatusInternalServerError, "删除失败")
		return
	}
	response.Success(c, "删除成功", file.Name)
}

// GetCredentialHandle 获取cos的临时凭证
func GetCredentialHandle(c *gin.Context) {
	credential, err := tencentSts.GetCredential()
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "获取临时密钥失败")
		return
	}
	response.Success(c, "获取临时密钥成功", credential)
}
