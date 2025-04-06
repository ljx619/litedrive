package user

import (
	"github.com/gin-gonic/gin"
	"litedrive/internal/models"
	"litedrive/pkg/serializer"
)

type UserService struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (service *UserService) Login(c *gin.Context) string {
	token, _ := models.LoginCheck(service.Username, service.Password)
	return token
}

func (service *UserService) Register(c *gin.Context) serializer.Response {
	u := models.User{
		Username: service.Username,
		Password: service.Password,
	}
	_, err := u.SaveUser()
	if err != nil {
		return serializer.ErrorResponse(err)
	}
	return serializer.SuccessResponse(gin.H{
		"id":       u.ID,
		"username": u.Username,
	}, "注册成功")

}
