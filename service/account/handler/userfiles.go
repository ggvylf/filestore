package handler

import (
	"context"

	proto "github.com/ggvylf/filestore/service/account/proto"
)

func (u *User) UserFiles(ctx context.Context, req *proto.ReqUserFiles, resp *proto.RespUserFiles) error {
	return nil
}

func (u *User) UserFileRename(ctx context.Context, req *proto.ReqUserFileRename, resp *proto.RespUserFileRename) error {
	return nil
}
