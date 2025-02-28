package main

import (
	"litedrive/internal/models"
	"litedrive/internal/router"
	"litedrive/internal/utils"
	"log"
	"strconv"
)

func init() {
	models.ConnectDatabase()
}

func main() {
	config, _ := utils.LoadConfig("./configs/config.yaml")

	api := router.InitRouter()
	if err := api.Run(":" + strconv.Itoa(config.Server.Port)); err != nil {
		log.Fatal(err)
	}

}
