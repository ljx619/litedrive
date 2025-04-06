package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"litedrive/internal/middlewares"
	"litedrive/internal/router/controllers"
	"time"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	// 配置 CORS 中间件,默认是放行所有跨域请求
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},                                                // 允许的源
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}, // 允许的方法
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},          // 允许的头部
		AllowCredentials: true,                                                         // 如果需要发送 cookies
		MaxAge:           12 * time.Hour,
	}))

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
		apiFiles.GET("/download/:fileID", controllers.DownloadFile) // 下载文件
		apiFiles.DELETE("/:fileID", controllers.DeleteFile)         // 删除文件
		apiFiles.PUT("/", controllers.RenameFile)                   // 文件重命名
		apiFiles.GET("/list", controllers.ListFiles)                // 获取用户文件列表
		apiFiles.GET("/downloadurl", controllers.DownloadURL)       // 获取下载链接
		apiFiles.POST("/rapidcheck", controllers.RapidCheck)        // 秒传接口
	}

	apiChunk := r.Group("/api/chunk")
	{
		apiChunk.Use(middlewares.JwtAuthMiddleware())
		apiChunk.POST("/initMultUpload", controllers.InitializeMultipartUpload)
		apiChunk.POST("/uploadPart", controllers.UploadPart)
		apiChunk.POST("/completeMultUpload", controllers.CompleteMultipartUpload)
	}

	return r
}
