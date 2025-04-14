package configs

import "litedrive/pkg/common"

const (
	// TempLocalRootDir : 本地临时存储地址的路径
	TempLocalRootDir = "/data/fileserver/"
	// TempPartRootDir : 分块文件在本地临时存储地址的路径
	TempPartRootDir = "/data/fileserver_part/"
	// CephRootDir : Ceph的存储路径prefix
	CephRootDir = "/ceph"
	// OSSRootDir : OSS的存储路径prefix
	OSSRootDir = "cos/"
	// CurrentStoreType : 设置当前文件的存储类型
	CurrentStoreType = common.StoreCOS
)
