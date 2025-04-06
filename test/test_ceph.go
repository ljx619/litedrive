package main

//func main() {
//ceph.InitCephClient()
//// 直接获取对象数据
//data, err := ceph.DownloadObject("testbucket1", "/ceph/73bf2b49f53267e2cd319404bf8b24bad8e3a4203bed90d845d657d07a8c0d9c/")
//if err != nil {
//	log.Fatalln("获取对象失败:", err)
//}
//
//// 创建目标文件
//tmpFile, err := os.Create("./storage/tmp/")
//if err != nil {
//	log.Fatal("创建文件失败:", err)
//}
//defer tmpFile.Close()
//
//// 将 data 写入文件
//filesize, err := tmpFile.Write(data)
//if err != nil {
//	log.Fatal("写入文件失败:", err)
//}
//
//// 计算 data 的 SHA-256 哈希值
//hash := sha256.Sum256(data)
//fileSha := hex.EncodeToString(hash[:])
//
//fmt.Println("文件大小：", filesize, "哈希值：", fileSha)
//}
