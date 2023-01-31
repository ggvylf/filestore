package mq

import (
	"github.com/ggvylf/filestore/common"
)

// 定义消息结构
type TransferData struct {
	FileHash      string
	CurLocation   string
	DestLocation  string
	DestStoreType common.StoreType
	FileSize      int64
	Bucket        string
}
