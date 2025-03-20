package models

import "gorm.io/gorm"

type UserFile struct {
	gorm.Model
	UserID   uint   `gorm:"not null;index:idx_user"`
	FileSha  string `gorm:"type:char(64);not null;index:idx_sha"`
	FileSize int64  `gorm:"not null;check:file_size >= 0"`
	FileName string `gorm:"type:varchar(255);not null"`
	Status   string `gorm:"type:varchar(20);default:'active';check:status IN ('active', 'deleted', 'locked')"`

	// 外键关联（假设用户表名为 users）
	User User `gorm:"foreignKey:UserID;references:ID"` // 关联用户表

	// 联合唯一索引（用户ID+文件哈希实现去重）
	_ struct{} `gorm:"uniqueIndex:idx_user_file"`
}

// OnUserFileUploadFinished: 更新用户文件表
func (u *UserFile) OnUserFileUploadFinished() error {
	return DB.Create(u).Error
}

// QueryUserFileMetas: 获取指定用户的所有文件列表
func QueryUserFileMetas(userid uint, limit int) ([]UserFile, error) {
	var files []UserFile
	if err := DB.Limit(limit).Where("user_id = ?", userid).Find(&files).Error; err != nil {
		return nil, err
	}
	return files, nil
}
