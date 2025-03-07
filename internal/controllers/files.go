package controllers

import (
	"github.com/gin-gonic/gin"
	"litedrive/internal/services/explorer"
	"net/http"
)

// TODO 将controller 跟 service 彻底解耦 包装响应体
func UploadFile(c *gin.Context) {
	// 实例化 fileService
	fileService := explorer.FileService{}
	res := fileService.UploadFile(c)
	c.JSON(http.StatusOK, res)
}
