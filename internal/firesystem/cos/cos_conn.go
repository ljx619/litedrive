package cos

import (
	"context"
	"github.com/tencentyun/cos-go-sdk-v5"
	"litedrive/internal/utils"
	"log"
	"net/http"
	"net/url"
	"os"
)

// 接入腾讯云COS对象存储

var CosClient *cos.Client

func InitCosClient() {
	// 加载配置文件
	config, _ := utils.LoadConfig()
	endpoint := config.Cos.Endpoint
	secretId := config.Cos.SecretID
	secretKey := config.Cos.SecretKey

	u, err := url.Parse(endpoint)
	if err != nil {
		log.Fatalf("解析 COS URL 失败: %v", err)
	}

	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  secretId,
			SecretKey: secretKey,
		},
	})

	// 连接测试
	_, _, err = c.Service.Get(context.Background())
	if err != nil {
		log.Fatalf("腾讯云 COS 连接失败: %v", err)
	}

	// 给全局变量赋值,方便service全局调用
	CosClient = c
	log.Println("腾讯云 COS 连接测试成功")
}

// 创建存储桶 暂未配置
func CreateBucket() error {
	_, err := CosClient.Bucket.Put(context.Background(), nil)
	if err != nil {
		log.Printf("创建存储桶失败: %v", err)
		return err
	}
	log.Println("存储桶创建成功")
	return nil
}

// 查询存储桶列表
func ListBuckets() ([]cos.Bucket, error) {
	res, _, err := CosClient.Service.Get(context.Background())
	if err != nil {
		log.Printf("查询存储桶列表失败: %v", err)
		return nil, err
	}
	return res.Buckets, nil
}

// 上传文件
func UploadFile(objectKey, filePath string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = CosClient.Object.Put(context.Background(), objectKey, f, nil)
	if err != nil {
		log.Printf("上传文件失败: %v", err)
		return err
	}
	log.Println("文件上传成功:", objectKey)
	return nil
}

// 删除文件
func DeleteFile(objectKey string) error {
	_, err := CosClient.Object.Delete(context.Background(), objectKey)
	if err != nil {
		log.Printf("删除文件失败: %v", err)
		return err
	}
	log.Println("文件删除成功:", objectKey)
	return nil
}

// 查询对象列表
func ListObjects(prefix string) ([]cos.Object, error) {
	opt := &cos.BucketGetOptions{
		Prefix:  prefix,
		MaxKeys: 100,
	}
	res, _, err := CosClient.Bucket.Get(context.Background(), opt)
	if err != nil {
		log.Printf("查询对象列表失败: %v", err)
		return nil, err
	}
	return res.Contents, nil
}

// 下载对象
func DownloadFile(objectKey, localPath string) error {
	_, err := CosClient.Object.GetToFile(context.Background(), objectKey, localPath, nil)
	if err != nil {
		log.Printf("下载对象失败: %v", err)
		return err
	}
	log.Println("文件下载成功:", objectKey)
	return nil
}
