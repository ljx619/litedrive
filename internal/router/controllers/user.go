package controllers

import (
	"github.com/gin-gonic/gin"
	"litedrive/internal/services/user"
	"litedrive/pkg/serializer"
	"net/http"
)

func Register(c *gin.Context) {
	var service user.UserService

	if err := c.ShouldBindJSON(&service); err != nil {
		c.JSON(http.StatusBadRequest, serializer.ErrorResponse(err))
		return
	}

	res := service.Register(c)

	c.JSON(http.StatusOK, res)
}

func Login(c *gin.Context) {
	var service user.UserService
	if err := c.ShouldBindJSON(&service); err == nil {
		res := service.Login(c)

		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": "登录成功",
			"data":    gin.H{"token": res},
		})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{
			"data": "用户名或密码错误",
		})
	}

}
