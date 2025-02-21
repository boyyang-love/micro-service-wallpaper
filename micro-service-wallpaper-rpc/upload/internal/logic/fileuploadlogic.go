package logic

import (
	"bytes"
	"context"
	"fmt"
	"github.com/boyyang-love/micro-service-wallpaper-rpc/upload/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-rpc/upload/pb/upload"
	"github.com/minio/minio-go/v7"
	"github.com/zeromicro/go-zero/core/logx"
)

type FileUploadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFileUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileUploadLogic {
	return &FileUploadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FileUploadLogic) FileUpload(in *upload.FileUploadReq) (*upload.FileUploadRes, error) {

	filePath := fmt.Sprintf("%s%s", in.Path, in.FileName)
	buf := bytes.NewReader(in.File)
	uploadInfo, err := l.svcCtx.
		MinioClient.
		PutObject(
			l.ctx,
			in.BucketName,
			filePath,
			buf,
			buf.Size(),
			minio.PutObjectOptions{},
		)
	if err != nil {
		return nil, err
	}

	return &upload.FileUploadRes{
		Base: &upload.Base{
			Code: 1,
			Msg:  "文件上传成功",
		},
		Data: &upload.FileUploadResData{
			FileName: in.FileName,
			Path:     filePath,
			ETag:     uploadInfo.ETag,
			Size:     uint64(uploadInfo.Size),
		},
	}, nil
}
