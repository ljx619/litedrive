package explorer

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/gin-gonic/gin"
	"io"
	"litedrive/internal/models"
	"litedrive/internal/utils"
	"litedrive/pkg/serializer"
	"os"
	"path/filepath"
	"strconv"
)

type FileService struct{}

func (s *FileService) UploadFile(c *gin.Context) serializer.Response {
	// 读取配置文件
	config, err := utils.LoadConfig()
	if err != nil {
		return serializer.ErrorResponse(err)
	}

	// 获取上下文中的 user_id
	userID, exists := c.Get("user_id")
	if !exists {
		return serializer.ErrorResponse(errors.New("user not logged in"))
	}
	// 转换为合适类型
	userIDInt := userID.(uint)

	// 获取上传的文件
	file, header, err := c.Request.FormFile("explorer")
	if err != nil {
		return serializer.ErrorResponse(err)
	}
	defer file.Close()

	//计算文件 SHA-256 哈希值
	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return serializer.ErrorResponse(err)
	}
	fileSha := hex.EncodeToString(hash.Sum(nil))

	//重置文件指针
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return serializer.ErrorResponse(err)
	}

	//秒传逻辑
	existingFile, err := models.GetFileBySha(fileSha)
	if err == nil && existingFile != nil {
		//说明数据库此文件已存在 直接绑定用户
		userFileRecord := &models.UserFile{
			UserID:   userIDInt,
			FileSha:  fileSha,
			FileSize: header.Size,
			FileName: header.Filename,
		}

		if err := userFileRecord.OnUserFileUploadFinished(); err != nil {
			return serializer.ErrorResponse(err)
		}

		return serializer.SuccessResponse(existingFile)
	}

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
		Sha:    fileSha,
		UserID: userIDInt,
		Name:   header.Filename,
		Size:   fileSize,
		Path:   filePath,
	}

	//调用 Model 层方法存入文件表
	if err := fileRecord.CreateFile(); err != nil {
		return serializer.ErrorResponse(err)
	}

	//绑定用户文件
	userFileRecord := &models.UserFile{
		UserID:   userIDInt,
		FileSha:  fileSha,
		FileSize: fileSize,
		FileName: header.Filename,
	}

	if err := userFileRecord.OnUserFileUploadFinished(); err != nil {
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
	// 获取上下文中的 user_id
	userID, exists := c.Get("user_id")
	if !exists {
		return serializer.ErrorResponse(errors.New("user id not found"))
	}
	// 转换为合适类型
	userIDInt := userID.(uint)

	limit, _ := strconv.Atoi(c.Query("limit"))

	allFiles, err := models.QueryUserFileMetas(userIDInt, limit)
	if err != nil {
		return serializer.ErrorResponse(err)
	}
	return serializer.SuccessResponse(allFiles)
}

//func TryFastUpload(c *gin.Context) serializer.Response {
//	//解析请求参数
//	// 获取上下文中的 user_id
//	userID, exists := c.Get("user_id")
//	if !exists {
//		return serializer.ErrorResponse(errors.New("user not logged in"))
//	}
//	// 转换为合适类型
//	userIDInt := userID.(uint)
//
//	// 获取上传的文件
//	file, header, err := c.Request.FormFile("explorer")
//	if err != nil {
//		return serializer.ErrorResponse(err)
//	}
//	defer file.Close()
//
//	//从文件表查找相同hash文件
//	//查不到记录则返回秒传失败
//	//上传过则将文件信息写入用户文件表，返回成功
//	return true
//}
