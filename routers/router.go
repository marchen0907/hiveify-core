package routers

import (
	"github.com/gin-gonic/gin"
	"hiveify-core/handle"
	"hiveify-core/interceptor"
)

// Router 路由
func Router() *gin.Engine {
	router := gin.Default()
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	// 设置上传文件的最大大小为4GB
	router.MaxMultipartMemory = 4 * 1024 * 1024 * 1024

	// 不需要校验token的接口
	router.GET("/user/register", handle.RegisterHandler)
	router.GET("/user/login", handle.LoginHandler)

	// 拦截器
	router.Use(interceptor.HTTPInterceptor())

	// 需要校验token的接口
	router.GET("/user/getUserInfo", handle.GetUserInfoHandler)

	router.POST("/file/upload", handle.UploadHandle)
	router.GET("/file/fastUpload", handle.FastUploadHandle)
	router.GET("/file/getFile", handle.GetFileHandle)
	router.GET("/file/downloadFile", handle.DownloadFileHandle)
	router.GET("/file/updateFile", handle.UpdateFileHandle)
	router.GET("/file/deleteFile", handle.DeleteFileHandle)

	router.GET("/file/chunk/initChunk", handle.InitChunkHandle)
	router.POST("/file/chunk/upload", handle.UploadChunkHandle)
	router.GET("/file/chunk/status", handle.UploadStatusHandle)
	router.GET("/file/chunk/cancel", handle.CancelUploadHandle)
	router.GET("/file/chunk/complete", handle.CompleteUploadHandle)

	router.GET("/file/getCredential", handle.GetCredentialHandle)
	router.GET("/file/completeCOSUpload", handle.CompleteCOSHandle)

	return router
}
