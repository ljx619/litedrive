package user

import (
	"github.com/gin-gonic/gin"
	"litedrive/internal/models"
)

type UserLoginService struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (service *UserLoginService) Login(c *gin.Context) string {
	token, _ := models.LoginCheck(service.Username, service.Password)
	return token
}
