package main

import (
	"encoding/json"
	"github.com/joho/godotenv"
	"litedrive/internal/cache/rabbitmq"
	"litedrive/internal/firesystem/cos"
	"litedrive/internal/models"
	"log"
	"os"
)

func ProcessTransfer(msg []byte) bool {
	//解析msg
	pubData := rabbitmq.TransferData{}
	err := json.Unmarshal(msg, &pubData)
	if err != nil {
		log.Println(err)
		return false
	}

	//根据临时存储文件路径，创建文件句柄
	file, err := os.Open(pubData.CurLocation)
	if err != nil {
		log.Println(err)
		return false
	}
	defer file.Close()

	//通过文件句柄将文件内容读出来并且上传到 cos
	err = cos.UploadFile(pubData.FileHash, file)
	if err != nil {
		log.Println(err)
		return false
	}

	//更新文件的存储路径到文件表
	err = models.UpdateFilePathBySha(pubData.FileHash, pubData.DestLocation)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

func main() {
	log.Println("开始监听转移消息队列...")

	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	models.InitDatabase()
	cos.InitCosClient()

	// ✅ **先手动初始化 RabbitMQ**
	if !rabbitmq.InitChannel() {
		log.Fatal("RabbitMQ 初始化失败")
	}

	rabbitmq.StartConsume(
		rabbitmq.TransCOSQueueName,
		"transfer_cos",
		ProcessTransfer)
}
