package config

var (

	// public
	Username = "admin"
	Token    = "06805bef4b8ddf7c3ea94d3fc12b8a8c63ef77b5"

	//fastupload
	FasttargetURL = "http://127.0.0.1:28080/file/fastupload"
	Fastfilehash  = "dd89e617b1cbcc2eb396e862b2ae2473be4c9f67"
	Fastfilename  = "a.txt"
	Fastfilesize  = "24474"

	// mupload
	Mpfilehash = "fe1d6ccb2544698b5c567411306e659de0fe922d" //sha1sum filename
	Mpfilesize = "148883574"
	Mpfilename = "/home/ggvylf/Downloads/go1.19.2.linux-amd64.tar.gz"
	MpinitURL  = "http://127.0.0.1:28080/file/mpupload/init"
	MppartURL  = "http://127.0.0.1:28080/file/mpupload/uppart"
	MpcompURL  = "http://127.0.0.1:28080/file/mpupload/complete"
)
