package svc

import (
	"fmt"
	"github.com/boyyang-love/micro-service-wallpaper-rpc/upload/helper"
	"github.com/boyyang-love/micro-service-wallpaper-rpc/upload/internal/config"

	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config      config.Config
	DB          *gorm.DB
	MinioClient *minio.Client
}

func NewServiceContext(c config.Config) *ServiceContext {

	db, err := helper.ConMySQL(c.MySQLConf)
	if err != nil {
		fmt.Printf("数据库连接失败(%s)\n", err.Error())
	} else {
		fmt.Println("数据库连接成功")
	}

	minioClient, err := helper.Minio(
		c.MinioClientConf.Endpoint,
		c.MinioClientConf.AccessKey,
		c.MinioClientConf.SecretKey,
		c.MinioClientConf.Secure,
	)
	if err != nil {
		fmt.Printf("对象存储连接失败(%s)\n", err.Error())
	} else {
		fmt.Println("对象存储连接成功")
	}

	return &ServiceContext{
		Config:      c,
		DB:          db,
		MinioClient: minioClient,
	}
}
