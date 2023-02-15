package config

const (
	// UploadServiceHost : 上传服务监听的地址
	UploadServiceHost = "0.0.0.0:8888"

	// UploadLBHost: 上传服务LB地址
	// 这里是对外的地址，跟具体服务监听的ip端口无关
	UploadLBHost = "http://127.0.0.1:28080"

	// DownloadLBHost: 下载服务LB地址
	DownloadLBHost = "http://127.0.0.1:38080"

	// TracerAgentHost: tracing agent地址
	TracerAgentHost = "127.0.0.1:6831"
)
