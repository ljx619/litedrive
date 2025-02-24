package models

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"litedrive/internal/utils"
	"log"
)

var DB *gorm.DB

func ConnectDatabase() {
	config, err := utils.LoadConfig("./configs/config.yaml")

	DB, err = gorm.Open(mysql.Open(config.Database.DSN), &gorm.Config{})

	if err != nil {
		log.Fatal("Error connecting to database")
	} else {
		fmt.Printf("Connect to database success")
	}

	DB.AutoMigrate(&User{}, &File{})

}
