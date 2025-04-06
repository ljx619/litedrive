package models

import (
	"errors"
	"gorm.io/gorm"
)

// File 文件元信息结构
type File struct {
	gorm.Model
	Sha  string `json:"sha" gorm:"type:char(64);unique;not null"`
	Size int64  `json:"size" gorm:"not null;check:size >= 0"`
	Path string `json:"path" gorm:"type:varchar(255);not null"`
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

// GetFileBySha 根据 SHA 获取文件
func GetFileBySha(sha string) (*File, error) {
	var file File
	if err := DB.Where("sha = ?", sha).First(&file).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // 没找到，返回 nil
		}
		return nil, err
	}
	return &file, nil
}

// UpdateFilePathBySha 根据 SHA 更新文件存储路径
func UpdateFilePathBySha(sha string, newPath string) error {
	if sha == "" || newPath == "" {
		return errors.New("sha 和 newPath 不能为空")
	}

	result := DB.Model(&File{}).Where("sha = ?", sha).Update("path", newPath)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("未找到匹配的文件，更新失败")
	}

	return nil
}
