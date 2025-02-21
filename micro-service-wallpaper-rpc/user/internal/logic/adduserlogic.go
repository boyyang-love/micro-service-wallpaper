package logic

import (
	"context"
	"errors"
	"fmt"
	"github.com/boyyang-love/micro-service-wallpaper-rpc/user/models"
	"gorm.io/gorm"

	"github.com/boyyang-love/micro-service-wallpaper-rpc/user/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-rpc/user/pb/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddUserLogic {
	return &AddUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddUserLogic) AddUser(in *user.AddUserReq) (*user.AddUserRes, error) {
	userInfo := models.User{
		Username: in.Username,
		Account:  in.Account,
		Password: in.Password,
		Role:     in.Role,
	}

	if err := l.svcCtx.
		DB.
		Model(&models.User{}).
		Select("id", "account", "username").
		Where("account = ? or username = ?", in.Account, in.Username).
		First(&userInfo).
		Error; errors.As(err, &gorm.ErrRecordNotFound) {
		if err = l.svcCtx.
			DB.
			Model(&models.User{}).
			Create(&userInfo).
			Error; err != nil {
			return nil, errors.New(fmt.Sprintf("create user failed: %v", err))
		}

		return &user.AddUserRes{
			Base: &user.Base{
				Code: 1,
				Msg:  "用户新增成功",
			},
			Data: &user.AddUserResData{
				Id:       uint64(userInfo.Id),
				Username: userInfo.Username,
				Account:  userInfo.Account,
			},
		}, nil
	} else {

		return nil, errors.New("该用户已经存在")
	}
}
