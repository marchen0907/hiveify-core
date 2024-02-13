package handle

import (
	"context"
	"errors"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"hiveify-core/cache/redis"
	"hiveify-core/database"
	"hiveify-core/response"
	"hiveify-core/util"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"
)

// ChunkInfo 文件分块上传信息
type ChunkInfo struct {
	Sha512      string `json:"sha512"` // 原始文件的sha512
	Size        int64  `json:"size"`   // 原始文件大小
	ChunkID     string `json:"chunkID"`
	ChunkSize   int64  `json:"chunkSize"`
	ChunkNum    int    `json:"chunkNum"`
	UploadChunk []int  `json:"uploadChunk"`
}

var ctx = context.Background()

// InitChunkHandle 初始化分块上传
func InitChunkHandle(c *gin.Context) {
	// 参数解析
	user, _ := util.GetUserByRequest(c)
	sha512 := c.Request.FormValue("sha512")
	size, err := strconv.ParseInt(c.Request.FormValue("size"), 10, 64)
	if err != nil {
		log.Warnf("参数非法：%s", err.Error())
		response.Fail(c, http.StatusBadRequest, "参数非法！")
		return
	}

	// 初始化分块上传信息
	chunkInfo := ChunkInfo{
		Sha512:      sha512,
		Size:        size,
		ChunkID:     user.Email + fmt.Sprintf("%x", time.Now().UnixNano()), // email + 时间戳
		ChunkSize:   10 * 1024 * 1024,                                      // 10MB
		ChunkNum:    int(math.Ceil(float64(size) / 10 * 1024 * 1024)),
		UploadChunk: nil,
	}

	// 分块信息写入缓存
	jsonData, err := sonic.Marshal(chunkInfo)
	if err != nil {
		log.Errorf("Failed to marshal chunk info: %s", err.Error())
		return
	}
	err = redis.Pool().Set(ctx, "chunk_"+chunkInfo.ChunkID, jsonData, 0).Err()
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "初始化分块上传失败！")
		return
	}
	response.Success(c, "初始化分块成功！", chunkInfo)
}

// UploadChunkHandle 上传分块
func UploadChunkHandle(c *gin.Context) {
	// 参数解析
	chunkID := c.Request.FormValue("chunkID")
	chunkIndex, err := strconv.Atoi(c.Request.FormValue("chunkIndex"))
	if err != nil {
		log.Warnf("参数非法：%s", err.Error())
		response.Fail(c, http.StatusBadRequest, "参数非法！")
		return
	}
	chunkSha512 := c.Request.FormValue("chunkSha512")

	// 检查分块信息
	jsonData, err := redis.Pool().Get(ctx, chunkID).Result()
	if err != nil {
		log.Warnf("Failed to get chunk info with chunkID %s, %s", chunkID, err.Error())
		response.Fail(c, http.StatusBadRequest, "分块信息校验失败！")
		return
	}
	var chunkInfo ChunkInfo
	err = sonic.Unmarshal([]byte(jsonData), &chunkInfo)
	if err != nil {
		log.Warnf("Failed to unmarshal chunk info with chunkID %s, %s", chunkID, err.Error())
		response.Fail(c, http.StatusInternalServerError, "分块文件上传失败！")
		return
	}

	// 获取分块文件并校验
	fileHeader, err := c.FormFile("file")
	if err != nil {
		log.Infof("Failed to get file, %s", err.Error())
		return
	}
	if ok := util.CheckFile(fileHeader, chunkSha512); !ok {
		response.Fail(c, http.StatusBadRequest, "分块文件校验失败！")
		return
	}

	// 分块文件存储
	if ok := util.SaveToLocal(fileHeader, strconv.Itoa(chunkIndex),
		"/tmp/"+chunkID+"/chunk_"+strconv.Itoa(chunkIndex)); ok {
		log.Errorf("Failed to chunk file to local, %s", err.Error())
		response.Fail(c, http.StatusInternalServerError, "分块文件存储失败！")
		return
	}

	// 更新分块信息
	chunkInfo.UploadChunk = append(chunkInfo.UploadChunk, chunkIndex)

	updateDate, err := sonic.Marshal(chunkInfo)
	if err != nil {
		log.Errorf("Failed to marshal chunk info: %s", err.Error())
		response.Fail(c, http.StatusInternalServerError, "分块文件上传失败！")
		return
	}
	err = redis.Pool().Set(ctx, "chunk_"+chunkInfo.ChunkID, updateDate, 0).Err()
	if err != nil {
		log.Errorf("Failed to update chunk info, %s", err.Error())
		response.Fail(c, http.StatusInternalServerError, "分块文件上传失败！")
		return
	}
	response.Success(c, "分块文件上传成功！", chunkInfo)
}

// UploadStatusHandle 获取上传状态
func UploadStatusHandle(c *gin.Context) {
	//  参数解析
	chunkID := c.Request.FormValue("chunkID")
	// 获取分块信息
	jsonData, err := redis.Pool().Get(ctx, "chunk_"+chunkID).Result()
	if err != nil {
		log.Warnf("Failed to get chunk info, %s", err.Error())
		response.Fail(c, http.StatusInternalServerError, "获取上传状态失败！")
		return
	}
	var chunkInfo ChunkInfo
	err = sonic.Unmarshal([]byte(jsonData), &chunkInfo)
	if err != nil {
		log.Warnf("Failed to unmarshal chunk info, %s", err.Error())
		response.Fail(c, http.StatusInternalServerError, "获取上传状态失败！")
		return
	}
	response.Success(c, "获取上传状态成功！", chunkInfo)
}

// CancelUploadHandle 取消上传
func CancelUploadHandle(c *gin.Context) {
	// 参数解析

	chunkID := c.Request.FormValue("chunkID")
	sha512 := c.Request.FormValue("sha512")
	user, _ := util.GetUserByRequest(c)

	// 删除redis中的分块信息
	err := redis.Pool().Del(ctx, "chunk_"+chunkID).Err()
	if err != nil {
		log.Warnf("Failed to delete chunk info: %s", err.Error())
		response.Fail(c, http.StatusInternalServerError, "取消上传失败！")
		return
	}

	//  删除本地存储的分块文件
	err = os.RemoveAll("/tmp/" + chunkID)
	if err != nil {
		log.Warnf("Failed to delete chunk info: %s", err.Error())
		response.Fail(c, http.StatusInternalServerError, "取消上传失败！")
		return
	}

	// 查询数据库中是否存在对应的文件信息
	// 存在则删除
	deleteFile := true
	file, err := database.GetFile(sha512)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			log.Warnf("获取文件信息失败：%s", err.Error())
			response.Fail(c, http.StatusInternalServerError, "取消上传失败！")
			return
		}
		// 文件信息不存在，无需删除
		deleteFile = false
	}

	if deleteFile {
		if ok := database.DeleteFile(file.Sha512); !ok {
			log.Warnf("删除文件信息失败：%s", err.Error())
			response.Fail(c, http.StatusInternalServerError, "取消上传失败！")
			return
		}
		// 成功删除文件信息
	}

	deleteUserFile := true
	userFile, err := database.GetUserFile(user.Email, sha512)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			log.Warnf("获取用户文件信息失败：%s", err.Error())
			response.Fail(c, http.StatusInternalServerError, "取消上传失败！")
			return
		}
		// 用户文件信息不存在，无需删除
		deleteUserFile = false
	}
	if deleteUserFile {
		// 用户文件信息存在，删除它
		if ok := database.DeleteUserFile(userFile.Email, userFile.Sha512); ok {
			response.Success(c, "取消上传成功！", chunkID)
		}
	}
}

// CompleteUploadHandle 合并分块文件
func CompleteUploadHandle(c *gin.Context) {
	// 参数解析
	chunkID := c.Request.FormValue("chunkID")
	sha512 := c.Request.FormValue("sha512")
	name := c.Request.FormValue("name")
	size, err := strconv.ParseInt(c.Request.FormValue("size"), 10, 64)
	if err != nil {
		log.Warnf("参数非法：%s", err.Error())
		response.Fail(c, http.StatusBadRequest, "参数非法！")
		return
	}
	user, _ := util.GetUserByRequest(c)

	// 分块上传情况校验
	jsonData, err := redis.Pool().Get(ctx, "chunk_"+chunkID).Result()
	if err != nil {
		log.Errorf("Failed to get chunk info from redis: %s", err.Error())
		response.Fail(c, http.StatusInternalServerError, "分块文件上传失败！")
		return
	}
	var chunkInfo ChunkInfo
	err = sonic.Unmarshal([]byte(jsonData), &chunkInfo)
	if err != nil {
		log.Errorf("Failed to unmarshal chunk info: %s", err.Error())
		response.Fail(c, http.StatusInternalServerError, "分块文件上传失败！")
		return
	}
	if sha512 != chunkInfo.Sha512 || chunkInfo.ChunkNum != len(chunkInfo.UploadChunk) || chunkInfo.Size != size {
		log.Errorf("Failed to verify chunk info: %s", err.Error())
		response.Fail(c, http.StatusInternalServerError, "分块文件上传失败！")
		return
	}

	// 合并分块文件
	if ok := util.MergeChunk("/tmp/"+name, "/tmp/"+chunkID); !ok {
		response.Fail(c, http.StatusInternalServerError, "分块文件上传失败！")
		return
	}

	// 文件信息写入数据库
	if ok := database.CreateFile(chunkInfo.Sha512, name, "/tmp/"+name, chunkInfo.Size); !ok {
		response.Fail(c, http.StatusInternalServerError, "分块文件上传失败！")
		return
	}
	if ok := database.CreateUserFile(user.Email, chunkInfo.Sha512, name); !ok {
		response.Fail(c, http.StatusInternalServerError, "分块文件上传失败！")
		return
	}

	// 删除redis中的分块信息
	err = redis.Pool().Del(ctx, "chunk_"+chunkID).Err()
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "分块文件上传失败！")
		return
	}
	response.Success(c, "分块文件上传成功！", database.File{
		Sha512:     sha512,
		Name:       name,
		Size:       chunkInfo.Size,
		Path:       "/tmp/" + name,
		CreateTime: time.Now(),
	})
}
