package main

func main() {
	// 连接到 Ceph
	//s3Client, err := ceph.CephConnect()
	//if err != nil {
	//	log.Fatalln("连接 Ceph 失败:", err)
	//}
	//bucketName := "testbucket1"
	//
	//// 创建存储桶
	//err = ceph.CreateBucket(s3Client, bucketName)
	//if err != nil {
	//	log.Fatalln("创建存储桶失败:", err)
	//}
	//fmt.Println("存储桶创建成功:", bucketName)
	//
	//// 查看存储桶中的对象
	//err = ceph.ListObjects(s3Client, bucketName)
	//if err != nil {
	//	log.Fatalln("查看存储桶中的对象失败:", err)
	//}
	//
	//// 上传一个对象
	//objectKey := "testfile1.txt"
	//objectContent := "这是一个测试文件内容1。"
	//
	//err = ceph.UploadObject(s3Client, bucketName, objectKey, objectContent)
	//if err != nil {
	//	log.Fatalln("上传对象失败:", err)
	//}
	//
	//fmt.Println("对象上传成功:", objectKey)
	//
	//// 再次查看存储桶中的对象
	//err = ceph.ListObjects(s3Client, bucketName)
	//if err != nil {
	//	log.Fatalln("查看存储桶中的对象失败:", err)
	//}
}
