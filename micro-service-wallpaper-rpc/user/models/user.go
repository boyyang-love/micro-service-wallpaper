package models

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"strings"
)

type User struct {
	Base
	// 字段
	Username string `json:"username" form:"username"`
	Account  string `json:"account" form:"account"`
	Password string `json:"password" form:"password"`
	Motto    string `json:"motto" form:"motto"`
	Address  string `json:"address" form:"address"`
	Tel      string `json:"tel" form:"tel"`
	Email    string `json:"email" form:"email"`
	QQ       string `json:"qq" form:"qq"`
	Wechat   string `json:"wechat" form:"wechat"`
	GitHub   string `json:"git_hub" form:"git_hub"`
	Role     string `json:"role" form:"role" gorm:"default:user"`
	Avatar   string `json:"avatar" form:"avatar"`
	Cover    string `json:"cover" form:"cover"`
}

func (u *User) TableName() string {
	return "user"
}

func (u *User) BeforeCreate(*gorm.DB) (err error) {
	uid := uuid.NewV1()
	u.Uid = strings.Replace(uid.String(), "-", "", -1)
	return nil
}
