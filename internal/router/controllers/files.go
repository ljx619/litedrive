package controllers

import (
	"github.com/gin-gonic/gin"
	"litedrive/internal/services/explorer"
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
	res := fileService.DownloadFile(c)
	c.JSON(http.StatusOK, res)
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
