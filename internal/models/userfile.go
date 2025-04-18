package models

import (
	"errors"
	"gorm.io/gorm"
)

type UserFile struct {
	gorm.Model
	UserID   uint   `gorm:"not null;uniqueIndex:idx_user_file"`                   // 联合唯一
	FileID   uint   `gorm:"not null;uniqueIndex:idx_user_file"`                   // 联合唯一
	FileName string `gorm:"type:varchar(255);not null;uniqueIndex:idx_user_file"` // 联合唯一

	DirID  uint   `gorm:"default:0;index"` // 所属目录
	Status string `gorm:"type:varchar(20);default:'active';check:status IN ('active', 'deleted', 'locked')"`
	// 外键关联
	User User `gorm:"foreignKey:UserID;references:ID"`
	File File `gorm:"foreignKey:FileID;references:ID"`
	//Dir  UserDir `gorm:"foreignKey:DirID;references:ID"`
}

// OnUserFileUploadFinished: 更新用户文件表
func (u *UserFile) OnUserFileUploadFinished() error {
	return DB.Create(u).Error
}

// QueryUserFileMetas: 获取指定用户的所有文件列表
func QueryUserFileMetas(userid uint, limit int) ([]UserFile, error) {
	var files []UserFile
	if err := DB.Preload("File").Limit(limit).Where("user_id = ?", userid).Find(&files).Error; err != nil {
		return nil, err
	}
	return files, nil
}

// GetUserFileByIDAndUser: 查询用户文件记录
func (u *UserFile) GetUserFileByIDFileNameAndUser() (*UserFile, error) {
	var userFile UserFile
	err := DB.Where("user_id = ? AND file_id = ? AND file_name = ?", u.UserID, u.FileID, u.FileName).First(&userFile).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // 没找到
		}
		return nil, err // 其他查询错误
	}
	return &userFile, nil
}

func (u *UserFile) UpdateFileName() error {
	if u.ID == 0 || u.FileName == "" {
		return errors.New("invalid user file ID or new file name")
	}

	// 使用 GORM 更新 file_name 字段
	err := DB.Model(&UserFile{}).
		Where("id = ? AND user_id = ?", u.ID, u.UserID).
		Update("file_name", u.FileName).Error

	return err
}
