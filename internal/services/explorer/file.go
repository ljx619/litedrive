package explorer

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"io"
	//"litedrive/internal/firesystem/ceph"
	"litedrive/internal/firesystem/cos"
	"litedrive/internal/models"
	"litedrive/internal/utils"
	"litedrive/pkg/serializer"
	"os"
	"path/filepath"
	"strconv"
)

type FileService struct{}

type RapidCheckService struct {
	FileName string `json:"fileName"`
	FileHash string `json:"fileHash"`
}

type RenameFileService struct {
	ID          uint   `json:"id"`          // UserFile.ID
	NewFileName string `json:"newFileName"` // 新的文件名
}

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

	//重置文件指针
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return serializer.ErrorResponse(err)
	}

	// TODO 经测试 此桶未作 sha 唯一验证 但是文件表中做了验证 所以数据会不一致
	// 同时写入到 ceph 存储
	//cephFilePath := "/ceph/" + fileSha + "/" + header.Filename
	//err = ceph.UploadObject("testbucket1", cephFilePath, file)
	//if err != nil {
	//	return serializer.ErrorResponse(err)
	//}
	//err = ceph.ListObjects("testbucket1")
	//if err != nil {
	//	return serializer.ErrorResponse(err)
	//}

	//TODO COS逻辑实现
	//cosFilePath := "/cos/" + fileSha + "/" + header.Filename
	// 暂时注释测试 rabbitmq
	//err = cos.UploadFile(filePath, file)
	//if err != nil {
	//	return serializer.ErrorResponse(err)
	//}
	//objects, err := cos.ListObjects("litedrive-filestore-1320309154")
	//if err != nil {
	//	return serializer.ErrorResponse(err)
	//}
	//fmt.Println(objects)

	////rabbitmq 部分
	//data := rabbitmq.TransferData{
	//	FileHash:      fileSha,
	//	CurLocation:   filePath,
	//	DestLocation:  cosFilePath,
	//	DestStoreType: common.StoreOSS,
	//}
	//pubData, _ := json.Marshal(data)
	//suc := rabbitmq.Publish(rabbitmq.TransExchangeName, rabbitmq.TransOSSRoutingKey, pubData)
	//if !suc {
	//	// 加入重拾发送消息逻辑
	//	return serializer.ErrorResponse(err)
	//}

	//创建文件记录
	fileRecord := &models.File{
		Sha:  fileSha,
		Size: fileSize,
		Path: filePath,
	}

	//调用 Model 层方法存入文件表
	if err := fileRecord.CreateFile(); err != nil {
		return serializer.ErrorResponse(err)
	}

	//绑定用户文件
	userFileRecord := &models.UserFile{
		UserID:   userIDInt,
		FileID:   fileRecord.ID,
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
	userID, exists := c.Get("user_id")
	if !exists {
		return serializer.ErrorResponse(errors.New("用户未登录"))
	}
	userIDInt := userID.(uint)

	fileIDStr := c.Param("fileID")
	fileIDUint, err := strconv.ParseUint(fileIDStr, 10, 64)
	if err != nil {
		return serializer.ErrorResponse(errors.New("无效的文件ID"))
	}

	// 权限校验
	var userFile models.UserFile
	if err := models.DB.Where("user_id = ? AND file_id = ?", userIDInt, uint(fileIDUint)).First(&userFile).Error; err != nil {
		return serializer.ErrorResponse(errors.New("文件不存在或无权限"))
	}

	// 获取文件信息
	var file models.File
	if err := models.DB.First(&file, "id = ?", userFile.FileID).Error; err != nil {
		return serializer.ErrorResponse(err)
	}

	// 检查文件是否存在
	if _, err := os.Stat(file.Path); os.IsNotExist(err) {
		return serializer.ErrorResponse(errors.New("文件已被删除或丢失"))
	}

	// 返回文件
	c.FileAttachment(file.Path, userFile.FileName)
	return serializer.SuccessResponse(file)
}

func (s *FileService) DeleteFile(c *gin.Context) serializer.Response {
	// 获取当前用户 ID
	userID, exists := c.Get("user_id")
	if !exists {
		return serializer.ErrorResponse(errors.New("用户未登录"))
	}
	userIDInt := userID.(uint)

	// 获取 fileID
	fileID := c.Param("fileID")

	// 查询 UserFile 确保用户有权限删除该文件
	var userFile models.UserFile
	if err := models.DB.Where("user_id = ? AND file_id = ?", userIDInt, fileID).First(&userFile).Error; err != nil {
		return serializer.ErrorResponse(errors.New("文件不存在或无权限"))
	}

	// 删除 UserFile 记录
	if err := models.DB.Unscoped().Delete(&userFile).Error; err != nil {
		return serializer.ErrorResponse(err)
	}

	// 检查该文件是否还有其他用户在使用
	var count int64
	if err := models.DB.Model(&models.UserFile{}).Where("file_id = ?", userFile.FileID).Count(&count).Error; err != nil {
		return serializer.ErrorResponse(err)
	}

	// 如果没有其他用户使用该文件，则物理删除文件
	if count == 0 {
		var file models.File
		if err := models.DB.First(&file, userFile.FileID).Error; err != nil {
			return serializer.ErrorResponse(err)
		}

		// 删除文件物理文件
		if err := os.Remove(file.Path); err != nil {
			return serializer.ErrorResponse(err)
		}

		// 删除 File 记录
		if err := models.DB.Unscoped().Delete(&file).Error; err != nil {
			return serializer.ErrorResponse(err)
		}
	}

	return serializer.SuccessResponse(nil, "文件删除成功")
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

// RapidCheck:秒传逻辑接口
func (s *RapidCheckService) RapidCheck(c *gin.Context) serializer.Response {
	// 获取上下文中的 user_id
	userID, exists := c.Get("user_id")
	if !exists {
		return serializer.ErrorResponse(errors.New("user not logged in"))
	}
	userIDInt := userID.(uint)

	// 从前端请求获取文件 SHA256 和 文件名
	fileSha := s.FileHash
	fileName := s.FileName
	if fileSha == "" {
		return serializer.ErrorResponse(errors.New("file SHA256 is required"))
	}
	if fileName == "" {
		return serializer.ErrorResponse(errors.New("file name is required"))
	}

	// 查找数据库，看该哈希值是否存在
	existingFile, err := models.GetFileBySha(fileSha)
	if err != nil || existingFile == nil {
		// 文件不存在，返回秒传失败
		return serializer.ErrorResponse(errors.New("秒传失败，文件不存在"))
	}

	// 绑定用户文件
	userFileRecord := &models.UserFile{
		UserID:   userIDInt,
		FileID:   existingFile.ID,
		FileName: fileName,
	}

	if err := userFileRecord.OnUserFileUploadFinished(); err != nil {
		return serializer.ErrorResponse(err)
	}

	return serializer.SuccessResponse(existingFile)
}

// TODO COS URL下载 关于下载时候的文件名 默认情况下 应该是上传时候指定的 key 如果想保持原来的文件名，可以在上传的时候指定一个元信息
func (s *FileService) DownloadURL(c *gin.Context) serializer.Response {
	filesha := c.PostForm("filesha")
	//从文件表中查找记录
	file, _ := models.GetFileBySha(filesha)

	// TODO 判断文件存储在Cos还是Ceph中

	//生成下载链接
	signedURL, _ := cos.DownloadURL(file.Path)
	return serializer.SuccessResponse(signedURL)
}

// RenameFile: 文件重命名
func (s *RenameFileService) RenameFile(c *gin.Context) serializer.Response {
	// 获取 user_id（可用于权限校验）
	userID, exists := c.Get("user_id")
	if !exists {
		return serializer.ErrorResponse(errors.New("user not logged in"))
	}
	userIDInt := userID.(uint)

	if s.ID == 0 || s.NewFileName == "" {
		return serializer.ErrorResponse(errors.New("缺少必要参数"))
	}

	// 查找 UserFile
	var userFile models.UserFile
	if err := models.DB.First(&userFile, "id = ? AND user_id = ?", s.ID, userIDInt).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return serializer.ErrorResponse(errors.New("文件记录不存在或无权限"))
		}
		return serializer.ErrorResponse(err)
	}

	// 更新文件名
	userFile.FileName = s.NewFileName
	if err := models.DB.Save(&userFile).Error; err != nil {
		return serializer.ErrorResponse(err)
	}

	return serializer.SuccessResponse("文件重命名成功")
}
