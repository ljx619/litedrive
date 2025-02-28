package router

import (
	"github.com/gin-gonic/gin"
	"litedrive/internal/controllers"
	"litedrive/internal/middlewares"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	api := r.Group("/api/auth")
	{
		api.POST("/register", controllers.Register)
		api.POST("/login", controllers.Login)
	}

	protected := r.Group("/api/admin")
	{
		protected.Use(middlewares.JwtAuthMiddleware())
		protected.GET("/user", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"code": 200,
				"msg":  "ok",
			})
		})
		//protected.GET("/files", controllers.Register())
	}

	return r
}
