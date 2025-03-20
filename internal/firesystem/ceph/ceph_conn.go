package ceph

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"golang.org/x/net/context"
	"litedrive/internal/utils"
	"log"
)

var CephClient *s3.Client

// InitCephClient 初始化 Ceph S3 客户端
func InitCephClient() {
	cephConfig, _ := utils.LoadConfig()
	endpoint := cephConfig.Ceph.Endpoint
	accessKey := cephConfig.Ceph.AccessKey
	secretKey := cephConfig.Ceph.SecretKey
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(""),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{
					URL:           endpoint,
					SigningRegion: region,
				}, nil
			},
		)),
	)
	if err != nil {
		log.Fatalf("加载 Ceph 配置失败: %v", err)
	}

	CephClient = s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true // Ceph 需要路径访问模式
	})
	log.Println("Ceph 客户端初始化成功")
}

// CheckBucket 检查存储桶是否存在
func CheckBucket(bucketName string) error {
	_, err := CephClient.HeadBucket(context.TODO(), &s3.HeadBucketInput{
		Bucket: &bucketName,
	})
	if err != nil {
		var notFoundErr *types.NotFound
		if errors.As(err, &notFoundErr) {
			return fmt.Errorf("存储桶 %s 不存在", bucketName)
		}
		return fmt.Errorf("检查存储桶失败: %w", err)
	}
	return nil
}

// CreateBucket 创建存储桶
func CreateBucket(bucketName string) error {
	_, err := CephClient.CreateBucket(context.TODO(), &s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		return fmt.Errorf("创建存储桶失败: %w", err)
	}
	return nil
}

// DeleteBucket 删除存储桶
func DeleteBucket(bucketName string) error {
	_, err := CephClient.DeleteBucket(context.TODO(), &s3.DeleteBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		return fmt.Errorf("删除存储桶失败: %w", err)
	}
	return nil
}

// ListObjects 列出存储桶中的对象
func ListObjects(bucketName string) error {
	resp, err := CephClient.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		return fmt.Errorf("列出对象失败: %w", err)
	}

	if len(resp.Contents) == 0 {
		fmt.Println("存储桶为空")
	} else {
		fmt.Println("存储桶中的对象:")
		for _, object := range resp.Contents {
			fmt.Printf("- %s\n", *object.Key)
		}
	}
	return nil
}

// UploadObject 上传对象
func UploadObject(bucketName, objectKey, content string) error {
	_, err := CephClient.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
		Body:   bytes.NewReader([]byte(content)),
	})
	if err != nil {
		return fmt.Errorf("上传对象失败: %w", err)
	}
	return nil
}

// DownloadObject 下载对象
func DownloadObject(bucketName, objectKey string) ([]byte, error) {
	resp, err := CephClient.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		return nil, fmt.Errorf("下载对象失败: %w", err)
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取对象数据失败: %w", err)
	}
	return buf.Bytes(), nil
}

// DeleteObject 删除对象
func DeleteObject(bucketName, objectKey string) error {
	_, err := CephClient.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		return fmt.Errorf("删除对象失败: %w", err)
	}
	return nil
}
