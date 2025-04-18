package models

import "gorm.io/gorm"

type UserDir struct {
	gorm.Model
	UserID   uint   `gorm:"not null;index"`                                                                    // 用户 ID，外键
	ParentID uint   `gorm:"default:0;index"`                                                                   // 父目录 ID，根目录为 0
	Name     string `gorm:"type:varchar(255);not null"`                                                        // 目录名
	Status   string `gorm:"type:varchar(20);default:'active';check:status IN ('active', 'deleted', 'locked')"` // 目录状态

	// 外键关联
	User User `gorm:"foreignKey:UserID;references:ID"`
}
