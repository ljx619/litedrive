package controllers

import (
	"github.com/gin-gonic/gin"
	"litedrive/internal/models"
	"litedrive/internal/services/user"
	"net/http"
)

type ReqController struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// TODO 将结构体从controller 解耦到 service中
func Register(c *gin.Context) {
	var req ReqController

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"data": err.Error()})
		return
	}

	u := models.User{
		Username: req.Username,
		Password: req.Password,
	}

	_, err := u.SaveUser()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"data": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "注册成功",
		"data": gin.H{
			"id":       u.ID,
			"username": u.Username,
		},
	})
}

func Login(c *gin.Context) {
	var server user.UserLoginService
	if err := c.ShouldBindJSON(&server); err == nil {
		res := server.Login(c)
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
