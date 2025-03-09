package models

import (
	"gorm.io/gorm"
)

// File 表示文件系统中的文件或目录
type File struct {
	gorm.Model
	UserID uint   `json:"user_id"`
	Name   string `json:"name"`
	Size   int64  `json:"size"`
	Path   string `json:"path"`
}

// CreateFile 创建文件记录
func (f *File) CreateFile() error {
	return DB.Create(f).Error
}

// DeleteFile 删除文件记录
func DeleteFile(fileID string) error {
	return DB.Delete(&File{}, "id = ?", fileID).Error
}

// GetFileByID 根据 ID 获取文件信息
func GetFileByID(fileID string) (*File, error) {
	var file File
	if err := DB.First(&file, "id = ?", fileID).Error; err != nil {
		return nil, err
	}
	return &file, nil
}

// GetAllFiles 获取所有文件信息
func GetAllFiles() ([]File, error) {
	var files []File
	if err := DB.Find(&files).Error; err != nil {
		return nil, err
	}
	return files, nil
}
