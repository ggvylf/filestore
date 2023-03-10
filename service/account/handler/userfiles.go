package handler

import (
	"context"
	"encoding/json"
	"log"

	"github.com/ggvylf/filestore/common"
	proto "github.com/ggvylf/filestore/service/account/proto"

	dbcli "github.com/ggvylf/filestore/service/dbproxy/client"
)

// UserFiles : 获取用户文件列表
func (u *User) UserFilesList(ctx context.Context, req *proto.ReqUserFile, res *proto.RespUserFile) error {
	dbResp, err := dbcli.QueryUserFileMetas(req.Username, int(req.Limit))
	if err != nil || !dbResp.Suc {
		log.Println(err.Error())
		res.Code = common.StatusServerError
		return err
	}

	// 格式化数据
	userFiles := dbcli.ToTableUserFiles(dbResp.Data)
	data, err := json.Marshal(userFiles)
	if err != nil {
		log.Println(err.Error())
		res.Code = common.StatusServerError
		return nil
	}

	res.FileData = data
	return nil
}

// UserFiles : 用户文件重命名
func (u *User) UserFileRename(ctx context.Context, req *proto.ReqUserFileRename, res *proto.RespUserFileRename) error {
	dbResp, err := dbcli.RenameFileName(req.Username, req.Filehash, req.NewFileName)
	if err != nil || !dbResp.Suc {
		log.Println(err.Error())
		res.Code = common.StatusServerError
		return err
	}

	userFiles := dbcli.ToTableUserFiles(dbResp.Data)
	data, err := json.Marshal(userFiles)
	if err != nil {
		log.Println(err.Error())
		res.Code = common.StatusServerError
		return nil
	}

	res.FileData = data
	return nil
}
