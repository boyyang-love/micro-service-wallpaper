package logic

import (
	"context"
	"errors"
	"fmt"
	"github.com/boyyang-love/micro-service-wallpaper-rpc/user/models"

	"github.com/boyyang-love/micro-service-wallpaper-rpc/user/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-rpc/user/pb/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoLogic {
	return &UserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserInfoLogic) UserInfo(in *user.UserInfoReq) (*user.UserInfoRes, error) {
	var userInfo models.User
	if err := l.svcCtx.
		DB.
		Model(&models.User{}).
		Select("Id", "Username", "Account", "Motto", "Address", "Tel", "Email", "QQ", "Wechat", "GitHub", "Role", "Avatar", "Cover").
		Where("id = ?", in.Id).
		First(&userInfo).
		Error; err != nil {
		return nil, errors.New(fmt.Sprint("数据库查询失败", err))
	}

	return &user.UserInfoRes{
		Base: &user.Base{
			Code: 1,
			Msg:  "查询成功",
		},
		Data: &user.UserInfoResData{
			Id:       uint64(userInfo.Id),
			Username: userInfo.Username,
			Account:  userInfo.Account,
			Motto:    userInfo.Motto,
			Address:  userInfo.Address,
			Tel:      userInfo.Tel,
			Email:    userInfo.Email,
			QQ:       userInfo.QQ,
			Wechat:   userInfo.Wechat,
			GitHub:   userInfo.GitHub,
			Role:     userInfo.Role,
			Avatar:   userInfo.Avatar,
			Cover:    userInfo.Cover,
		},
	}, nil
}
