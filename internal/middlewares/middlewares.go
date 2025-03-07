package middlewares

import (
	"github.com/gin-gonic/gin"
	"litedrive/internal/utils"
	"net/http"
)

func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := utils.TokenValid(c); err != nil {
			c.String(http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}

		// 从 token 中提取 user_id
		userID, err := utils.ExtractTokenID(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to extract user_id"})
			c.Abort()
			return
		}

		// 将 user_id 存入上下文
		c.Set("user_id", userID)

		// 继续执行请求
		c.Next()
	}
}
