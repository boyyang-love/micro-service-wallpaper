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

type UpdateUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserLogic {
	return &UpdateUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateUserLogic) UpdateUser(in *user.UpdateUserReq) (*user.UpdateUserRes, error) {

	if err := l.svcCtx.
		DB.
		Model(&models.User{}).
		Where("id = ?", in.Id).
		Updates(&models.User{
			Motto:   in.Motto,
			Address: in.Address,
			Tel:     in.Tel,
			Email:   in.Email,
			QQ:      in.QQ,
			Wechat:  in.Wechat,
			GitHub:  in.GitHub,
			Role:    in.Role,
			Avatar:  in.Avatar,
			Cover:   in.Cover,
		}).
		Error; err != nil {
		return nil, errors.New(fmt.Sprint("更新用户信息失败:", err))
	}

	return &user.UpdateUserRes{
		Base: &user.Base{
			Code: 1,
			Msg:  "更新用户信息成功",
		},
	}, nil
}
