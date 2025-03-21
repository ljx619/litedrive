package router

import (
	"github.com/gin-gonic/gin"
	"litedrive/internal/middlewares"
	"litedrive/internal/router/controllers"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	api := r.Group("/api/auth")
	{
		api.POST("/register", controllers.Register)
		api.POST("/login", controllers.Login)
	}

	apiFiles := r.Group("/api/files")
	{
		apiFiles.Use(middlewares.JwtAuthMiddleware())
		apiFiles.POST("/upload", controllers.UploadFile)            // 上传文件
		apiFiles.GET("/:fileID", controllers.GetFileInfo)           // 获取文件信息
		apiFiles.GET("/:fileID/download", controllers.DownloadFile) // 下载文件
		apiFiles.DELETE("/:fileID", controllers.DeleteFile)         // 删除文件
		apiFiles.GET("/list", controllers.ListFiles)                // 获取文件列表
		//apiFiles.PUT("/:fileID/move", controllers.MoveFile)
	}

	apiChunk := r.Group("/api/chunk")
	{
		apiChunk.Use(middlewares.JwtAuthMiddleware())
		apiChunk.POST("/initMultUpload", controllers.InitializeMultipartUpload)
		apiChunk.POST("/uploadPart", controllers.UploadPart)
		apiChunk.POST("/completemultupload", controllers.CompleteMultipartUpload)
	}

	return r
}
