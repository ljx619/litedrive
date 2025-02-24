package main

import (
	"github.com/gin-gonic/gin"
	"litedrive/internal/controllers"
	"litedrive/internal/middlewares"
	"litedrive/internal/models"
)

func init() {
	models.ConnectDatabase()
}

func main() {
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

	r.Run(":8080")
}
