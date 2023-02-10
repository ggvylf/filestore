package client

import (
	"context"
	"encoding/json"

	"github.com/ggvylf/filestore/service/dbproxy/orm"
	dbProto "github.com/ggvylf/filestore/service/dbproxy/proto"

	"github.com/mitchellh/mapstructure"
	"go-micro.dev/v4"
)

// FileMeta : 文件元信息结构
type FileMeta struct {
	FileSha1 string
	FileName string
	FileSize int64
	Location string
	UploadAt string
}

// 定义一个代理客户端
var (
	dbCli dbProto.DBProxyService
)

func init() {
	service := micro.NewService()
	// 初始化， 解析命令行参数等
	service.Init()
	// 初始化一个dbproxy服务的客户端
	dbCli = dbProto.NewDBProxyService("go.micro.service.dbproxy", service.Client())
}

// execAction : 向dbproxy请求执行action
// 执行的所有router参考dbproxy/mapper/mapper.go
func execAction(funcName string, paramJson []byte) (*dbProto.RespExec, error) {
	return dbCli.ExecuteAction(context.TODO(), &dbProto.ReqExec{
		Action: []*dbProto.SingleAction{
			&dbProto.SingleAction{
				Name:   funcName,
				Params: paramJson,
			},
		},
	})
}

// parseBody : 转换rpc返回的结果
func parseBody(resp *dbProto.RespExec) *orm.ExecResult {
	if resp == nil || resp.Data == nil {
		return nil
	}
	resList := []orm.ExecResult{}
	_ = json.Unmarshal(resp.Data, &resList)
	// TODO:
	if len(resList) > 0 {
		return &resList[0]
	}
	return nil
}

// 填充fm结构体
func TableFileToFileMeta(tfile orm.TableFile) FileMeta {
	return FileMeta{
		FileSha1: tfile.FileHash,
		FileName: tfile.FileName.String,
		FileSize: tfile.FileSize.Int64,
		Location: tfile.FileAddr.String,
	}
}

// 转换tbl_user格式
func ToTableUser(src interface{}) orm.TableUser {
	user := orm.TableUser{}
	mapstructure.Decode(src, &user)
	return user
}

// 转换tbl_file格式
func ToTableFile(src interface{}) orm.TableFile {
	file := orm.TableFile{}
	mapstructure.Decode(src, &file)
	return file
}

// 转换tbl_user_file格式
func ToTableUserFile(src interface{}) orm.TableUserFile {
	ufile := orm.TableUserFile{}
	mapstructure.Decode(src, &ufile)
	return ufile
}

// 批量转换tbl_user_file格式
func ToTableUserFiles(src interface{}) []orm.TableUserFile {
	ufile := []orm.TableUserFile{}
	mapstructure.Decode(src, &ufile)
	return ufile
}

// 从tbl_file获取fm信息
func GetFileMeta(filehash string) (*orm.ExecResult, error) {
	uInfo, _ := json.Marshal([]interface{}{filehash})
	res, err := execAction("/file/GetFileMeta", uInfo)
	return parseBody(res), err
}

// 获取fmlist信息
func GetFileMetaList(limitCnt int) (*orm.ExecResult, error) {
	uInfo, _ := json.Marshal([]interface{}{limitCnt})
	res, err := execAction("/file/GetFileMetaList", uInfo)
	return parseBody(res), err
}

// OnFileUploadFinished : 新增/更新文件元信息到tbl_file
func OnFileUploadFinished(fmeta FileMeta) (*orm.ExecResult, error) {
	uInfo, _ := json.Marshal([]interface{}{fmeta.FileSha1, fmeta.FileName, fmeta.FileSize, fmeta.Location})
	res, err := execAction("/file/OnFileUploadFinished", uInfo)
	return parseBody(res), err
}

// 更新tbl_file中文件的location
func UpdateFileLocation(filehash, location string) (*orm.ExecResult, error) {
	uInfo, _ := json.Marshal([]interface{}{filehash, location})
	res, err := execAction("/file/UpdateFileLocation", uInfo)
	return parseBody(res), err
}

// 用户注册 更新tbl_user
func UserSignup(username, encPasswd string) (*orm.ExecResult, error) {
	uInfo, _ := json.Marshal([]interface{}{username, encPasswd})
	res, err := execAction("/user/UserSignup", uInfo)
	return parseBody(res), err
}

// 用户登录 查询tbl_user
func UserSignin(username, encPasswd string) (*orm.ExecResult, error) {
	uInfo, _ := json.Marshal([]interface{}{username, encPasswd})
	res, err := execAction("/user/UserSignin", uInfo)
	return parseBody(res), err
}

// 从tbl_user中获取用户信息
func GetUserInfo(username string) (*orm.ExecResult, error) {
	uInfo, _ := json.Marshal([]interface{}{username})
	res, err := execAction("/user/GetUserInfo", uInfo)
	return parseBody(res), err
}

// 判断用户是否存在tbl_user
func UserExist(username string) (*orm.ExecResult, error) {
	uInfo, _ := json.Marshal([]interface{}{username})
	res, err := execAction("/user/UserExist", uInfo)
	return parseBody(res), err
}

// 更新tbl_user_token的token
func UpdateToken(username, token string) (*orm.ExecResult, error) {
	uInfo, _ := json.Marshal([]interface{}{username, token})
	res, err := execAction("/user/UpdateToken", uInfo)
	return parseBody(res), err
}

// 从tbl_user_file中查询指定用户的文件
func QueryUserFileMeta(username, filehash string) (*orm.ExecResult, error) {
	uInfo, _ := json.Marshal([]interface{}{username, filehash})
	res, err := execAction("/ufile/QueryUserFileMeta", uInfo)
	return parseBody(res), err
}

// 从tbl_user_file中查询用户下的所有文件
func QueryUserFileMetas(username string, limit int) (*orm.ExecResult, error) {
	uInfo, _ := json.Marshal([]interface{}{username, limit})
	res, err := execAction("/ufile/QueryUserFileMetas", uInfo)
	return parseBody(res), err
}

// OnUserFileUploadFinished : 新增/更新fm信息到tbl_user_file
func OnUserFileUploadFinished(username string, fmeta FileMeta) (*orm.ExecResult, error) {
	uInfo, _ := json.Marshal([]interface{}{username, fmeta.FileSha1,
		fmeta.FileName, fmeta.FileSize})
	res, err := execAction("/ufile/OnUserFileUploadFinished", uInfo)
	return parseBody(res), err
}

// 文件重命名 更新tbl_file
func RenameFileName(username, filehash, filename string) (*orm.ExecResult, error) {
	uInfo, _ := json.Marshal([]interface{}{username, filehash, filename})
	res, err := execAction("/ufile/RenameFileName", uInfo)
	return parseBody(res), err
}
