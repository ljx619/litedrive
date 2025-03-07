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
		return serializer.Response{
			Code:  http.StatusInternalServerError,
			Data:  nil,
			Msg:   "",
			Error: err.Error(),
		}
	}

	// 获取上下文中的 user_id
	userID, exists := c.Get("user_id")
	if !exists {
		return serializer.Response{
			Code:  http.StatusInternalServerError,
			Data:  nil,
			Msg:   "User not logged in",
			Error: err.Error(),
		}
	}
	// 转换为合适类型
	userIDInt := userID.(uint)

	// 获取上传的文件
	file, header, err := c.Request.FormFile("explorer")
	if err != nil {
		return serializer.Response{
			Code:  http.StatusInternalServerError,
			Data:  nil,
			Msg:   "",
			Error: err.Error(),
		}
	}
	defer file.Close()

	// 生成存储路径
	filePath := filepath.Join(config.Storage.Root, header.Filename)
	// 创建文件
	outFile, err := os.Create(filePath)
	if err != nil {
		return serializer.Response{
			Code:  http.StatusInternalServerError,
			Data:  nil,
			Msg:   "",
			Error: err.Error(),
		}
	}
	defer outFile.Close()

	//写入文件
	fileSize, err := io.Copy(outFile, file)
	if err != nil {
		return serializer.Response{
			Code:  http.StatusInternalServerError,
			Data:  nil,
			Msg:   "",
			Error: err.Error(),
		}
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
		return serializer.Response{
			Code:  http.StatusInternalServerError,
			Data:  nil,
			Msg:   "model 调用错误",
			Error: err.Error(),
		}
	}
	return serializer.Response{
		Code: http.StatusOK,
		Data: fileRecord,
		Msg:  "操作成功",
	}
}
