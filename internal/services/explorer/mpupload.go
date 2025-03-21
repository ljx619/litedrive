package explorer

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"io"
	"litedrive/internal/cache/redis"
	"litedrive/internal/models"
	"litedrive/internal/utils"
	"litedrive/pkg/serializer"
	"log"
	"math"
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
	if err := c.ShouldBindJSON(&fileInfo); err != nil {
		return serializer.ErrorResponse(errors.New("参数解析错误"))
	}

	upInfo := MultipartUploadInfo{
		FileHash:   fileInfo.Sha,
		FileSize:   fileInfo.Size,
		UploadID:   strconv.FormatInt(time.Now().UnixNano(), 10),
		ChunkSize:  5 * 1024 * 1024, // 5MB
		ChunkCount: int(math.Ceil(float64(fileInfo.Size) / (5 * 1024 * 1024))),
	}

	ctx, cancel := context.WithTimeout(redis.Ctx, 3*time.Second)
	defer cancel()

	err := redis.RedisCli.HSet(ctx, "MP_"+upInfo.UploadID,
		"filehash", upInfo.FileHash,
		"filesize", upInfo.FileSize,
		"chunkcount", upInfo.ChunkCount,
	).Err()
	if err != nil {
		log.Printf("Redis HSet 失败: %v", err)
		return serializer.ErrorResponse(errors.New("Redis 写入失败"))
	}

	return serializer.SuccessResponse(upInfo)
}

func UploadPart(c *gin.Context) serializer.Response {
	type req struct {
		UserID     int    `json:"user_id"`
		UploadID   string `json:"upload_id"`
		ChunkIndex int    `json:"chunk_index"`
	}
	var reqInfo req
	if err := c.ShouldBindJSON(&reqInfo); err != nil {
		return serializer.ErrorResponse(errors.New("参数解析错误"))
	}

	config, err := utils.LoadConfig()
	if err != nil {
		return serializer.ErrorResponse(errors.New("配置文件加载失败"))
	}

	chunkDir := filepath.Join(config.Storage.Root, reqInfo.UploadID)
	if err := os.MkdirAll(chunkDir, os.ModePerm); err != nil {
		return serializer.ErrorResponse(errors.New("无法创建存储目录"))
	}

	chunkPath := filepath.Join(chunkDir, strconv.Itoa(reqInfo.ChunkIndex))
	fd, err := os.Create(chunkPath)
	if err != nil {
		return serializer.ErrorResponse(errors.New("无法创建分块文件"))
	}
	defer fd.Close()

	file, _, err := c.Request.FormFile("explorer")
	if err != nil {
		return serializer.ErrorResponse(errors.New("文件上传失败"))
	}
	defer file.Close()

	buf := make([]byte, 1024*1024)
	for {
		n, err := file.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return serializer.ErrorResponse(errors.New("文件读取错误"))
		}
		if _, err := fd.Write(buf[:n]); err != nil {
			return serializer.ErrorResponse(errors.New("文件写入失败"))
		}
	}

	ctx, cancel := context.WithTimeout(redis.Ctx, 3*time.Second)
	defer cancel()

	err = redis.RedisCli.HSet(ctx, "MP_"+reqInfo.UploadID,
		"chunk_"+strconv.Itoa(reqInfo.ChunkIndex), 1,
	).Err()
	if err != nil {
		return serializer.ErrorResponse(errors.New("Redis 更新失败"))
	}

	return serializer.SuccessResponse("分块上传成功")
}

func CompleteMultipartUpload(c *gin.Context) serializer.Response {
	type Req struct {
		UserID   int    `json:"user_id"`
		UploadID string `json:"upload_id"`
	}
	var reqInfo Req

	if err := c.ShouldBindJSON(&reqInfo); err != nil {
		return serializer.ErrorResponse(errors.New("参数解析错误"))
	}

	config, err := utils.LoadConfig()
	if err != nil {
		return serializer.ErrorResponse(errors.New("配置文件加载失败"))
	}

	ctx, cancel := context.WithTimeout(redis.Ctx, 3*time.Second)
	defer cancel()

	chunkCountStr, err := redis.RedisCli.HGet(ctx, "MP_"+reqInfo.UploadID, "chunkcount").Result()
	if err != nil {
		return serializer.ErrorResponse(errors.New("无法获取分块数"))
	}
	chunkCount, _ := strconv.Atoi(chunkCountStr)

	fileHash, err := redis.RedisCli.HGet(ctx, "MP_"+reqInfo.UploadID, "filehash").Result()
	if err != nil {
		return serializer.ErrorResponse(errors.New("无法获取文件 hash"))
	}

	mergedFilePath := filepath.Join(config.Storage.Root, fileHash)
	outFile, err := os.Create(mergedFilePath)
	if err != nil {
		return serializer.ErrorResponse(errors.New("创建合并文件失败"))
	}
	defer outFile.Close()

	for i := 0; i < chunkCount; i++ {
		chunkPath := filepath.Join(config.Storage.Root, reqInfo.UploadID, strconv.Itoa(i))
		chunkFile, err := os.Open(chunkPath)
		if err != nil {
			return serializer.ErrorResponse(errors.New("读取分块失败"))
		}
		_, err = io.Copy(outFile, chunkFile)
		chunkFile.Close()
		if err != nil {
			return serializer.ErrorResponse(errors.New("合并分块失败"))
		}
		os.Remove(chunkPath)
	}

	fileInfo, err := outFile.Stat()
	if err != nil {
		return serializer.ErrorResponse(errors.New("获取合并文件信息失败"))
	}

	redis.RedisCli.Del(ctx, "MP_"+reqInfo.UploadID)

	return serializer.SuccessResponse(map[string]interface{}{
		"message":   "文件上传完成",
		"file_path": mergedFilePath,
		"file_size": fileInfo.Size(),
		"file_hash": fileHash,
	})
}
