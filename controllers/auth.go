package controllers

import (
	"github.com/gin-gonic/gin"
	"litedrive/models"
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
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "register success",
		"data":    req,
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
		c.JSON(http.StatusBadRequest, gin.H{
			"data": "username or password is incorrect.",
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})

}
