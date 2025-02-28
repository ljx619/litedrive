package models

import (
	"gorm.io/gorm"
	"time"
)

// File 表示文件系统中的文件或目录
type File struct {
	gorm.Model
	Name        string     `json:"name" gorm:"not null"`
	Path        string     `json:"path" gorm:"not null"`
	Size        int64      `json:"size" gorm:"default:0"`
	Type        string     `json:"type" gorm:"not null"` // file 或 directory
	ContentType string     `json:"content_type"`
	UserID      uint       `json:"user_id" gorm:"not null"`
	ParentID    *uint      `json:"parent_id"`
	IsDeleted   bool       `json:"is_deleted" gorm:"default:false"`
	ShareCode   string     `json:"share_code"`
	ShareExpire *time.Time `json:"share_expire"`
}

// FileService 文件服务接口
type FileService interface {
	Upload(userID uint, filename string, parentID *uint, fileData []byte, contentType string) (*File, error)
	Download(userID uint, fileID uint) (*File, []byte, error)
	List(userID uint, parentID *uint) ([]File, error)
	Delete(userID uint, fileID uint) error
	CreateDirectory(userID uint, name string, parentID *uint) (*File, error)
	Move(userID uint, fileID uint, newParentID *uint) error
	Rename(userID uint, fileID uint, newName string) error
	Share(userID uint, fileID uint, expireDays int) (string, error)
	GetByShareCode(shareCode string) (*File, []byte, error)
}
