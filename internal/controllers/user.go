package controllers

import (
	"github.com/gin-gonic/gin"
	"litedrive/internal/models"
	"net/http"
)

type ReqController struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

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
	var req ReqController
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"data": err.Error()})
		return
	}

	u := models.User{
		Username: req.Username,
		Password: req.Password,
	}

	token, err := models.LoginCheck(u.Username, u.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"data": "用户名或密码错误",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "登录成功",
		"data":    gin.H{"token": token},
	})

}
