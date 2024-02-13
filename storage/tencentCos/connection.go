package tencentCos

import (
	"context"
	"github.com/charmbracelet/log"
	"github.com/tencentyun/cos-go-sdk-v5"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"time"
)

var cosClient *cos.Client
var ctx = context.Background()

// Client 创建CosClient对象
func client() *cos.Client {
	if cosClient == nil {
		// 创建 CosClient 对象
		u, _ := url.Parse("https://cos.wujunyi.cn")
		baseURL := &cos.BaseURL{BucketURL: u}
		// 1.永久密钥
		cosClient = cos.NewClient(baseURL, &http.Client{
			Transport: &cos.AuthorizationTransport{
				SecretID:  os.Getenv("COS_SECRET_ID"),
				SecretKey: os.Getenv("COS_SECRET_KEY"),
			},
		})
	}
	return cosClient
}

// GetBucketList 获取存储桶列表
func GetBucketList() *cos.ServiceGetResult {
	bucketList, _, err := client().Service.Get(ctx, nil)
	if err != nil {
		log.Errorf("Get bucket list failed, err: %s", err.Error())
		return nil
	}
	return bucketList
}

// PutFile 上传文件到cos
func PutFile(path string, fileHeader *multipart.FileHeader) bool {
	file, err := fileHeader.Open()
	if err != nil {
		log.Warnf("Failed to open file, err: %s", err.Error())
		return false
	}
	_, err = client().Object.Put(ctx, path, file, nil)
	if err != nil {
		return false
	}
	return true
}

// GetUploadURL 获取预签名上传链接
func GetUploadURL(path string) (string, error) {
	uploadURL, err := client().Object.GetPresignedURL2(ctx, http.MethodPut, path, time.Hour*12, nil)
	if err != nil {
		log.Warnf("Failed to get presigned URL, err: %s", err.Error())
		return "", err
	}
	return uploadURL.String(), nil
}

// CheckFileExist 检测文件是否存在于存储桶中
func CheckFileExist(path string) bool {
	ok, err := client().Object.IsExist(ctx, path)
	if err != nil {
		return false
	}
	return ok
}

// GetDownloadURL 获取预签名下载链接
func GetDownloadURL(path string) (string, error) {
	downloadURL, err := client().Object.GetPresignedURL2(ctx, http.MethodGet, path, time.Hour*12, nil)
	if err != nil {
		log.Warnf("Get presigned url failed, err: %s", err.Error())
		return "", err
	}
	return downloadURL.String(), nil
}
