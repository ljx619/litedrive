package controllers

import (
	"context"
	"github.com/gin-gonic/gin"
	"litedrive/internal/cache/redis"
	"litedrive/internal/models"
	"litedrive/internal/utils"
	"litedrive/pkg/serializer"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// MultipartUploadInfo: 初始化分块信息
type MultipartUploadInfo struct {
	FileHash   string
	FileSize   int64
	UploadID   string
	ChunkSize  int
	ChunkCount int
}

func InitalMultipartUpload(c *gin.Context) serializer.Response {
	var fileInfo models.File
	//解析用户参数
	if err := c.ShouldBindJSON(&fileInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	//生成分块上传的初始化信息
	upInfo := MultipartUploadInfo{
		FileHash:   fileInfo.Sha,
		FileSize:   fileInfo.Size,
		UploadID:   strconv.FormatInt(time.Now().UnixNano(), 10),
		ChunkSize:  5 * 1024 * 1024, // 5MB
		ChunkCount: int(math.Ceil(float64(fileInfo.Size) / (5 * 1024 * 1024))),
	}

	ctx, cancel := context.WithTimeout(redis.Ctx, 3*time.Second) // 每次调用都新建一个带超时的 ctx
	defer cancel()

	//将初始化信息写入到 redis 缓存
	err := redis.RedisCli.HSet(ctx, "MP_"+upInfo.UploadID,
		"filehash", upInfo.FileHash,
		"filesize", upInfo.FileSize,
		"chunkcount", upInfo.ChunkCount,
	).Err()
	if err != nil {
		log.Fatalf("Redis HSet 失败: %v", err)
	}

	return serializer.SuccessResponse(upInfo)
}

func UploadPart(c *gin.Context) {
	//解析用户请求参数
	type req struct {
		UserID     int    `json:"user_id"`
		UploadID   string `json:"upload_id"`
		ChunkIndex int    `json:"chunk_index"`
	}
	var reqInfo req
	if err := c.ShouldBindJSON(&reqInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	// 加载配置文件
	config, err := utils.LoadConfig()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "配置文件加载失败"})
		return
	}

	// 创建存储分块的目录
	chunkDir := filepath.Join(config.Storage.Root, reqInfo.UploadID)
	if err := os.MkdirAll(chunkDir, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法创建存储目录"})
		return
	}

	// 分块文件路径
	chunkPath := filepath.Join(chunkDir, strconv.Itoa(reqInfo.ChunkIndex))

	// 打开文件句柄，准备写入
	fd, err := os.Create(chunkPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法创建分块文件"})
		return
	}
	defer fd.Close()

	// 获取文件流
	file, _, err := c.Request.FormFile("explorer")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件上传失败"})
		return
	}
	defer file.Close()

	// 读取并写入文件
	buf := make([]byte, 1024*1024) // 1MB 缓冲区
	for {
		n, err := file.Read(buf)
		if err != nil {
			if err.Error() == "EOF" {
				break // 读取完成
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "文件读取错误"})
			return
		}
		if _, err := fd.Write(buf[:n]); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "文件写入失败"})
			return
		}
	}

	//更新redis缓存数据
	// 记录该分块已上传
	ctx, cancel := context.WithTimeout(redis.Ctx, 3*time.Second)
	defer cancel()

	err = redis.RedisCli.HSet(ctx, "MP_"+reqInfo.UploadID,
		"chunk_"+strconv.Itoa(reqInfo.ChunkIndex), 1,
	).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Redis 更新失败"})
		return
	}

	//返回处理结果
	c.JSON(http.StatusOK, gin.H{"message": "分块上传成功"})
}

// CompleteMultipartUpload 完成分块上传并合并文件
//func CompleteMultipartUpload(c *gin.Context) {
//	// 解析用户请求参数
//	type Req struct {
//		UserID   int    `json:"user_id"`
//		UploadID string `json:"upload_id"`
//	}
//	var reqInfo Req
//
//	if err := c.ShouldBindJSON(&reqInfo); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//
//	// 加载配置文件
//	config, err := utils.LoadConfig()
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "配置文件加载失败"})
//		return
//	}
//
//	// 获取分块信息
//	ctx, cancel := context.WithTimeout(redis.Ctx, 3*time.Second)
//	defer cancel()
//
//	// 获取 chunkcount（分块数）和 filehash 等信息
//	chunkCount, err := redis.RedisCli.HGet(ctx, "MP_"+reqInfo.UploadID, "chunkcount").Int()
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法获取分块数"})
//		return
//	}
//
//	fileHash, err := redis.RedisCli.HGet(ctx, "MP_"+reqInfo.UploadID, "filehash").Result()
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法获取文件 hash"})
//		return
//	}
//
//	// 检查所有分块是否上传完成
//	for i := 0; i < chunkCount; i++ {
//		exists, err := redis.RedisCli.HExists(ctx, "MP_"+reqInfo.UploadID, "chunk_"+strconv.Itoa(i)).Result()
//		if err != nil || !exists {
//			c.JSON(http.StatusBadRequest, gin.H{"error": "有分块未上传"})
//			return
//		}
//	}
//
//	// 创建合并后的文件存储路径
//	uploadDir := filepath.Join(config.Storage.Root, reqInfo.UploadID)
//	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法创建目录"})
//		return
//	}
//
//	// 合并所有分块文件
//	mergedFilePath := filepath.Join(uploadDir, fileHash)
//	outFile, err := os.Create(mergedFilePath)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "文件合并失败"})
//		return
//	}
//	defer outFile.Close()
//
//	// 按分块顺序写入文件
//	for i := 0; i < chunkCount; i++ {
//		chunkPath := filepath.Join(uploadDir, strconv.Itoa(i))
//		chunkFile, err := os.Open(chunkPath)
//		if err != nil {
//			c.JSON(http.StatusInternalServerError, gin.H{"error": "读取分块文件失败"})
//			return
//		}
//
//		// 将分块内容写入合并文件
//		if _, err := outFile.ReadFrom(chunkFile); err != nil {
//			c.JSON(http.StatusInternalServerError, gin.H{"error": "合并分块失败"})
//			return
//		}
//		chunkFile.Close()
//
//		// 删除分块文件，节省存储空间
//		if err := os.Remove(chunkPath); err != nil {
//			log.Printf("删除分块文件失败: %v", err)
//		}
//	}
//
//	// 清理 Redis 中的分块记录
//	redis.RedisCli.Del(ctx, "MP_"+reqInfo.UploadID)
//
//	// 返回合并后的文件信息
//	c.JSON(http.StatusOK, gin.H{
//		"message":     "文件上传完成",
//		"file_path":   mergedFilePath,
//		"file_size":   outFile.Stat().Size(),
//		"file_hash":   fileHash,
//		"upload_time": time.Now().Format(time.RFC3339),
//	})
//}
