package explorer

import (
	"github.com/gin-gonic/gin"
	"io"
	"litedrive/internal/models"
	"litedrive/internal/utils"
	"litedrive/pkg/serializer"
	"net/http"
	"os"
	"path/filepath"
)

type FileService struct{}

func (s *FileService) UploadFile(c *gin.Context) serializer.Response {
	// 读取配置文件
	config, err := utils.LoadConfig("./configs/config.yaml")
	if err != nil {
		return serializer.ErrorResponse(err)
	}

	// 获取上下文中的 user_id
	userID, exists := c.Get("user_id")
	if !exists {
		return serializer.Response{
			Code: http.StatusInternalServerError,
			Data: nil,
			Msg:  "User not logged in",
		}
	}
	// 转换为合适类型
	userIDInt := userID.(uint)

	// 获取上传的文件
	file, header, err := c.Request.FormFile("explorer")
	if err != nil {
		return serializer.ErrorResponse(err)
	}
	defer file.Close()

	// 生成存储路径
	filePath := filepath.Join(config.Storage.Root, header.Filename)
	// 创建文件
	outFile, err := os.Create(filePath)
	if err != nil {
		return serializer.ErrorResponse(err)
	}
	defer outFile.Close()

	//写入文件
	fileSize, err := io.Copy(outFile, file)
	if err != nil {
		return serializer.ErrorResponse(err)
	}
	//创建文件记录
	fileRecord := &models.File{
		UserID: userIDInt,
		Name:   header.Filename,
		Size:   fileSize,
		Path:   filePath,
	}

	//调用 Model 层方法存入数据库
	if err := fileRecord.CreateFile(); err != nil {
		return serializer.ErrorResponse(err)
	}
	return serializer.SuccessResponse(fileRecord)
}

func (s *FileService) GetFileInfo(c *gin.Context) serializer.Response {
	fileID := c.Param("fileID")
	file, err := models.GetFileByID(fileID)
	if err != nil {
		return serializer.ErrorResponse(err, "文件记录获取失败")
	}
	return serializer.SuccessResponse(file)
}
func (s *FileService) DownloadFile(c *gin.Context) serializer.Response {
	fileID := c.Param("fileID")
	file, err := models.GetFileByID(fileID)
	if err != nil {
		return serializer.ErrorResponse(err)
	}
	c.FileAttachment(file.Path, file.Name)
	return serializer.SuccessResponse(file)
}
func (s *FileService) DeleteFile(c *gin.Context) serializer.Response {
	fileID := c.Param("fileID")
	file, err := models.GetFileByID(fileID)
	if err != nil {
		return serializer.ErrorResponse(err)
	}
	if err := os.Remove(file.Path); err != nil {
		return serializer.ErrorResponse(err)
	}
	err = models.DeleteFile(fileID)
	if err != nil {
		return serializer.ErrorResponse(err)
	}
	return serializer.SuccessResponse(file, file.Name+"文件成功删除")
}
func (s *FileService) ListFiles(c *gin.Context) serializer.Response {
	allFiles, err := models.GetAllFiles()
	if err != nil {
		return serializer.ErrorResponse(err)
	}
	return serializer.SuccessResponse(allFiles)
}
