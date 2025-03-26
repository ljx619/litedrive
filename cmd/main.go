package main

import (
	"github.com/joho/godotenv"
	"litedrive/internal/cache/redis"
	"litedrive/internal/firesystem/ceph"
	"litedrive/internal/models"
	"litedrive/internal/router"
	"litedrive/internal/utils"
	"log"
	"strconv"
)

func init() {
	models.InitDatabase()
	redis.InitRedis()
	ceph.InitCephClient()
}

func main() {
	defer models.CloseDatabase()
	defer redis.CloseRedis()
	//加载环境变量
	_ = godotenv.Load()
	//加载配置文件
	config, _ := utils.LoadConfig()
	//注册路由
	api := router.InitRouter()
	if err := api.Run(":" + strconv.Itoa(config.Server.Port)); err != nil {
		log.Fatal(err)
	}

}
