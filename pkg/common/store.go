package common

import "strings"

// 存储类型(表示文件存到哪里)
type StoreType int

const (
	_ StoreType = iota
	// StoreLocal : 节点本地
	StoreLocal
	// StoreCeph : Ceph集群
	StoreCeph
	// StoreCOS : 腾讯COS
	StoreCOS
	// StoreMix : 混合(Ceph及OSS)
	StoreMix
	// StoreAll : 所有类型的存储都存一份数据
	StoreAll
)

// ParseStoreType 将字符串转换为 StoreType 枚举
func ParseStoreType(t string) StoreType {
	switch strings.ToLower(t) {
	case "local":
		return StoreLocal
	case "ceph":
		return StoreCeph
	case "cos":
		return StoreCOS
	case "mix":
		return StoreMix
	case "all":
		return StoreAll
	default:
		return StoreLocal
	}
}
