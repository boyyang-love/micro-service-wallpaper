package svc

import (
	"fmt"
	"github.com/boyyang-love/micro-service-wallpaper-api/helper"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/config"
	"github.com/boyyang-love/micro-service-wallpaper-rpc/upload/uploadclient"
	"github.com/zeromicro/go-zero/zrpc"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config        config.Config
	UploadService uploadclient.Upload
	DB            *gorm.DB
}

func NewServiceContext(c config.Config) *ServiceContext {
	uploadClient := zrpc.MustNewClient(c.UploadRpc)

	db, err := helper.ConMySQL(c.MySQLConf)
	if err != nil {
		fmt.Printf("数据库连接失败(%s)\n", err.Error())
	} else {
		fmt.Println("数据库连接成功")
	}
	return &ServiceContext{
		Config:        c,
		UploadService: uploadclient.NewUpload(uploadClient),
		DB:            db,
	}
}
