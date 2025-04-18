package controllers

import (
	"github.com/gin-gonic/gin"
	"litedrive/internal/services/directory"
	"net/http"
)

func CreateDir(c *gin.Context) {
	var dirService directory.DirService

	if err := c.ShouldBind(&dirService); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res := dirService.CreateDir(c)
	c.JSON(http.StatusOK, res)
}

func ListSubDirs(c *gin.Context) {
	var dirService directory.DirService
	res := dirService.ListSubDirs(c)
	c.JSON(http.StatusOK, res)
}

func DeleteDir(c *gin.Context) {
	var dirService directory.DeleteDirService
	if err := c.ShouldBindJSON(&dirService); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	res := dirService.DeleteDir(c)
	c.JSON(http.StatusOK, res)
}

func RenameDir(c *gin.Context) {
	var dirService directory.RenameDirService
	if err := c.ShouldBindJSON(&dirService); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	res := dirService.RenameDir(c)
	c.JSON(http.StatusOK, res)
}

func ListDirFiles(c *gin.Context) {
	res := directory.ListDirFiles(c)
	c.JSON(http.StatusOK, res)
}
