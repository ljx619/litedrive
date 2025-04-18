package directory

import (
	"errors"
	"github.com/gin-gonic/gin"
	"litedrive/internal/models"
	"litedrive/pkg/serializer"
	"strconv"
)

type DirService struct {
	Name     string `json:"name" binding:"required"`
	ParentID uint   `json:"parentId"`
}

type DeleteDirService struct {
	DirID uint `json:"dir_id" binding:"required"`
}

type RenameDirService struct {
	DirID   uint   `json:"dir_id" binding:"required"` // 要重命名的目录ID
	NewName string `json:"name" binding:"required"`   // 新的目录名称
}

type ListDirFilesResponse struct {
	Dirs  []models.UserDir  `json:"dirs"`
	Files []models.UserFile `json:"files"`
}

func (service *DirService) CreateDir(c *gin.Context) serializer.Response {
	// 获取上下文中的 user_id
	userID, exists := c.Get("user_id")
	if !exists {
		return serializer.ErrorResponse(errors.New("user not logged in"))
	}
	// 转换为合适类型
	userIDInt := userID.(uint)

	// 目录不能重名
	var count int64
	models.DB.Model(&models.UserDir{}).
		Where("user_id = ? AND parent_id = ? AND name = ? AND status = 'active'", userID, service.ParentID, service.Name).
		Count(&count)
	if count > 0 {
		return serializer.ErrorResponse(errors.New("目录已存在"))
	}

	dir := models.UserDir{
		UserID:   userIDInt,
		ParentID: service.ParentID,
		Name:     service.Name,
		Status:   "active",
	}

	if err := models.DB.Create(&dir).Error; err != nil {
		return serializer.ErrorResponse(err)
	}

	return serializer.SuccessResponse(dir, "目录创建成功")
}

func (service *DirService) ListSubDirs(c *gin.Context) serializer.Response {
	// 获取上下文中的 user_id
	userID, exists := c.Get("user_id")
	if !exists {
		return serializer.ErrorResponse(errors.New("user not logged in"))
	}
	parentID, _ := strconv.ParseUint(c.Query("parent_id"), 10, 64)

	// 转换为合适类型
	userIDInt := userID.(uint)

	var dirs []models.UserDir
	if err := models.DB.Where("user_id = ? AND parent_id = ? AND status = 'active'", userIDInt, parentID).
		Find(&dirs).Error; err != nil {
		return serializer.ErrorResponse(err, "查询子目录失败")
	}
	return serializer.SuccessResponse(dirs)
}

func (service *DeleteDirService) DeleteDir(c *gin.Context) serializer.Response {
	// 获取上下文中的 user_id
	userID, exists := c.Get("user_id")
	if !exists {
		return serializer.ErrorResponse(errors.New("user not logged in"))
	}

	// 转换为合适类型
	userIDInt := userID.(uint)

	var dir models.UserDir
	if err := models.DB.
		Where("id = ? AND user_id = ?", service.DirID, userID).
		First(&dir).Error; err != nil {
		return serializer.ErrorResponse(err, "目录不存在")
	}

	var subDirCount int64
	models.DB.Model(&models.UserDir{}).
		Where("parent_id = ? AND user_id = ?", service.DirID, userID).
		Count(&subDirCount)
	if subDirCount > 0 {
		return serializer.ErrorResponse(errors.New(""), "该目录下存在子目录，无法删除")
	}

	var fileCount int64
	models.DB.Model(&models.UserFile{}).
		Where("dir_id = ? AND user_id = ?", service.DirID, userID).
		Count(&fileCount)
	if fileCount > 0 {
		return serializer.ErrorResponse(errors.New("该目录下存在文件，无法删除"))
	}

	if err := models.DB.
		Unscoped().
		Where("id = ? AND user_id = ?", service.DirID, userIDInt).
		Delete(&models.UserDir{}).Error; err != nil {
		return serializer.ErrorResponse(err, "删除目录失败")
	}

	return serializer.SuccessResponse("删除目录成功")

}

func (service *RenameDirService) RenameDir(c *gin.Context) serializer.Response {
	// 获取上下文中的 user_id
	userID, exists := c.Get("user_id")
	if !exists {
		return serializer.ErrorResponse(errors.New("user not logged in"))
	}

	// 转换为合适类型
	userIDInt := userID.(uint)

	// 检查是否有重名目录在当前目录下
	var dir models.UserDir
	if err := models.DB.
		Where("id = ? AND user_id = ?", service.DirID, userIDInt).
		First(&dir).Error; err != nil {
		return serializer.ErrorResponse(err, "目录不存在")
	}

	var count int64
	if err := models.DB.
		Model(&models.UserDir{}).
		Where("user_id = ? AND parent_id = ? AND name = ? AND id != ?", userID, dir.ParentID, service.NewName, service.DirID).
		Count(&count).Error; err != nil {
		return serializer.ErrorResponse(err, "重命名失败")
	}
	if count > 0 {
		return serializer.ErrorResponse(errors.New("同一目录下已存在同名目录"))
	}

	// 更新目录名
	if err := models.DB.
		Model(&models.UserDir{}).
		Where("id = ? AND user_id = ?", service.DirID, userID).
		Update("name", service.NewName).Error; err != nil {
		return serializer.ErrorResponse(err, "目录重命名失败")
	}

	return serializer.SuccessResponse("", "重命名成功")
}

func ListDirFiles(c *gin.Context) serializer.Response {
	var dirs []models.UserDir
	var files []models.UserFile

	// 获取上下文中的 user_id
	userID, exists := c.Get("user_id")
	if !exists {
		return serializer.ErrorResponse(errors.New("user not logged in"))
	}
	// 转换为合适类型
	userIDInt := userID.(uint)
	dirIDStr := c.Query("dirId")
	dirID, err := strconv.ParseUint(dirIDStr, 10, 64)
	if err != nil {
		return serializer.ErrorResponse(err, "无效的目录 ID")
	}

	// 查询目录
	if err := models.DB.Where("user_id = ? AND parent_id = ?", userIDInt, dirID).Find(&dirs).Error; err != nil {
		return serializer.ErrorResponse(err, "获取子目录失败")
	}

	// 查询文件
	if err := models.DB.Where("user_id = ? AND dir_id = ?", userIDInt, dirID).Preload("File").Find(&files).Error; err != nil {
		return serializer.ErrorResponse(err, "获取文件失败")
	}

	return serializer.SuccessResponse(ListDirFilesResponse{
		Dirs:  dirs,
		Files: files,
	})
}
