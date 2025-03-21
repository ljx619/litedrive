package controllers

import (
	"github.com/gin-gonic/gin"
	"litedrive/internal/services/explorer"
	"net/http"
)

// InitializeMultipartUpload 初始化分块上传
func InitializeMultipartUpload(c *gin.Context) {
	res := explorer.InitalMultipartUpload(c)
	c.JSON(http.StatusOK, res)
}

// UploadPart 处理分块上传
func UploadPart(c *gin.Context) {
	res := explorer.UploadPart(c)
	c.JSON(http.StatusOK, res)
}

// CompleteMultipartUpload 处理分块合并
func CompleteMultipartUpload(c *gin.Context) {
	res := explorer.CompleteMultipartUpload(c)
	c.JSON(http.StatusOK, res)
}
