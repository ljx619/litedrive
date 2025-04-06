package main

import (
	"github.com/joho/godotenv"
	"litedrive/internal/cache/redis"
	"litedrive/internal/firesystem/ceph"
	"litedrive/internal/firesystem/cos"
	"litedrive/internal/models"
	"litedrive/internal/router"
	"litedrive/internal/utils"
	"log"
	"strconv"
)

func init() {
	godotenv.Load()
	models.InitDatabase()
	redis.InitRedis()
	ceph.InitCephClient()
	cos.InitCosClient()
}

func main() {
	defer models.CloseDatabase()
	defer redis.CloseRedis()
	//加载配置文件
	config, _ := utils.LoadConfig()
	//注册路由
	api := router.InitRouter()
	if err := api.Run(":" + strconv.Itoa(config.Server.Port)); err != nil {
		log.Fatal(err)
	}

}
