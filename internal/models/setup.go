package models

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"litedrive/internal/utils"
	"log"
)

var DB *gorm.DB

func InitDatabase() {
	config, err := utils.LoadConfig()

	DB, err = gorm.Open(mysql.Open(config.Database.DSN), &gorm.Config{})

	if err != nil {
		log.Fatal("Error connecting to database")
	} else {
		log.Println("Connected to database successfully")
	}

	DB.AutoMigrate(&User{}, &File{}, &UserFile{}, &UserDir{})

}

// CloseDatabase 关闭数据库连接
func CloseDatabase() {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			log.Printf("Failed to get sqlDB: %v", err)
			return
		}
		_ = sqlDB.Close()
		log.Println("Database connection closed.")
	}
}
