package controllers

import (
	"github.com/gin-gonic/gin"
	"litedrive/internal/services/explorer"
	"litedrive/pkg/serializer"
	"net/http"
)

func UploadFile(c *gin.Context) {
	// 实例化 fileService
	fileService := explorer.FileService{}
	res := fileService.UploadFile(c)
	c.JSON(http.StatusOK, res)
}

func GetFileInfo(c *gin.Context) {
	fileService := explorer.FileService{}
	res := fileService.GetFileInfo(c)
	c.JSON(http.StatusOK, res)
}

func DownloadFile(c *gin.Context) {
	fileService := explorer.FileService{}
	_ = fileService.DownloadFile(c)
	// 多次响应会报错响应内
	//c.JSON(http.StatusOK, res)
}

func DeleteFile(c *gin.Context) {
	fileService := explorer.FileService{}
	res := fileService.DeleteFile(c)
	c.JSON(http.StatusOK, res)
}

func ListFiles(c *gin.Context) {
	fileService := explorer.FileService{}
	res := fileService.ListFiles(c)
	c.JSON(http.StatusOK, res)
}

func DownloadURL(c *gin.Context) {
	fileService := explorer.FileService{}
	res := fileService.DownloadURL(c)
	c.JSON(http.StatusOK, res)
}

func RapidCheck(c *gin.Context) {
	var service explorer.RapidCheckService
	if err := c.ShouldBindJSON(&service); err != nil {
		c.JSON(http.StatusBadRequest, serializer.ErrorResponse(err))
		return
	}
	res := service.RapidCheck(c)

	c.JSON(http.StatusOK, res)
}

func RenameFile(c *gin.Context) {
	var service explorer.RenameFileService
	if err := c.ShouldBindJSON(&service); err != nil {
		c.JSON(http.StatusBadRequest, serializer.ErrorResponse(err))
		return
	}
	res := service.RenameFile(c)
	c.JSON(http.StatusOK, res)
}
