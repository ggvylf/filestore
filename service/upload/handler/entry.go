package handler

import (
	"context"

	cfg "github.com/ggvylf/filestore/service/upload/config"
	upProto "github.com/ggvylf/filestore/service/upload/proto"
)

// Upload : upload结构体
type Upload struct{}

// UploadEntry : 获取上传入口
func (u *Upload) UploadEntry(
	ctx context.Context,
	req *upProto.ReqEntry,
	res *upProto.RespEntry) error {

	res.Entry = cfg.UploadEntry
	return nil
}
