package svc

import (
	"fmt"
	"github.com/boyyang-love/micro-service-wallpaper-rpc/user/helper"
	"github.com/boyyang-love/micro-service-wallpaper-rpc/user/internal/config"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config config.Config
	DB     *gorm.DB
}

func NewServiceContext(c config.Config) *ServiceContext {
	db, err := helper.ConMySQL(c.MySQLConf)
	if err != nil {
		fmt.Printf("mysql connect error: %s", err)
	} else {
		fmt.Println("mysql connect success")
	}
	return &ServiceContext{
		DB:     db,
		Config: c,
	}
}
