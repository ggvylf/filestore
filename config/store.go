package config

import "github.com/ggvylf/filestore/common"

const (
	// TempLocalRootDir : 本地临时存储地址的路径
	TempLocalRootDir = "/tmp"
	// TempPartRootDir : 分块文件在本地临时存储地址的路径
	TempPartRootDir = "/tmp/fileserver_part"
	// CephRootDir : Ceph的存储路径prefix
	CephRootDir = "/ceph"
	// OSSRootDir : OSS的存储路径prefix
	OSSRootDir = "/minio"
	// CurrentStoreType : 设置当前文件的存储类型
	CurrentStoreType = common.StoreOSS
)
