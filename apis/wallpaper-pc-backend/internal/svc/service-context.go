package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"upload/uploadclient"
	"wallpaper-pc-backend/internal/config"
)

type ServiceContext struct {
	Config        config.Config
	UploadService uploadclient.Upload
}

func NewServiceContext(c config.Config) *ServiceContext {
	uploadClient := zrpc.MustNewClient(c.UploadRpc)
	return &ServiceContext{
		Config:        c,
		UploadService: uploadclient.NewUpload(uploadClient),
	}
}
