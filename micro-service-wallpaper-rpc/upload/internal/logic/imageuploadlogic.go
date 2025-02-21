package logic

import (
	"context"
	"fmt"
	"github.com/boyyang-love/micro-service-wallpaper-rpc/upload/helper"
	"github.com/minio/minio-go/v7"

	"github.com/boyyang-love/micro-service-wallpaper-rpc/upload/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-rpc/upload/pb/upload"

	"github.com/zeromicro/go-zero/core/logx"
)

type ImageUploadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewImageUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ImageUploadLogic {
	return &ImageUploadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ImageUploadLogic) ImageUpload(in *upload.ImageUploadReq) (*upload.ImageUploadRes, error) {

	comp, err := helper.Image2Webp(&in.File, uint(in.Quality))
	if err != nil {
		return nil, err
	}

	imagePath := fmt.Sprintf("%s/%s.webp", in.Path, helper.FileNameNoExt(in.FileName))
	oriImagePath := fmt.Sprintf("%s/%s", in.OriPath, in.FileName)
	uploadInfo, err := l.svcCtx.
		MinioClient.
		PutObject(
			l.ctx,
			in.BucketName,
			imagePath,
			comp.Buf,
			comp.Size,
			minio.PutObjectOptions{},
		)
	if err != nil {
		return nil, err
	}

	oriUploadInfo, err := l.svcCtx.
		MinioClient.
		PutObject(
			l.ctx,
			in.BucketName,
			oriImagePath,
			comp.OriBuf,
			comp.OriSize,
			minio.PutObjectOptions{},
		)
	if err != nil {
		return nil, err
	}

	return &upload.ImageUploadRes{
		Base: &upload.Base{
			Code: 1,
			Msg:  "图片上传成功",
		},
		Data: &upload.ImageUploadResData{
			Path:    imagePath,
			OriPath: oriImagePath,
			ETag:    uploadInfo.ETag,
			OriETag: oriUploadInfo.ETag,
			Size:    uint64(uploadInfo.Size),
			OriSize: uint64(oriUploadInfo.Size),
		},
	}, nil

}
