package handler

import (
	"context"

	cfg "github.com/ggvylf/filestore/service/download/config"
	dlProto "github.com/ggvylf/filestore/service/download/proto"
)

// Dwonload :download结构体
type Download struct{}

// DownloadEntry : 获取下载入口
func (u *Download) DownloadEntry(
	ctx context.Context,
	req *dlProto.ReqDownloadEntry,
	res *dlProto.RespDownloadEntry) error {

	res.Entry = cfg.DownloadEntry
	return nil
}
