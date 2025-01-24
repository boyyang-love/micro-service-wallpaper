package logic

import (
	"context"
	"errors"
	"github.com/boyyang-love/micro-service-wallpaper-rpc/upload/helper"
	"github.com/boyyang-love/micro-service-wallpaper-rpc/upload/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-rpc/upload/models"
	"github.com/boyyang-love/micro-service-wallpaper-rpc/upload/pb/upload"
	"gorm.io/gorm"

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
	compressedImage, err := helper.Image2Webp(in.File, uint(in.Quality))
	if err != nil {
		return nil, err
	}

	urls, err := helper.MinioFileUpload(&helper.MinioFileUploadParams{
		Ctx:         l.ctx,
		MinioClient: l.svcCtx.MinioClient,
		Buf:         compressedImage.Buf,
		OriBuf:      compressedImage.OriBuf,
		FileHash:    compressedImage.Hash,
		Folder:      in.Folder,
		OriFolder:   in.OriFolder,
		FileName:    in.FileName,
		FileType:    in.FileType,
		BucketName:  in.BucketName,
	})
	if err != nil {
		return nil, err
	}

	if err = l.InsertToMySql(&models.Upload{
		Hash:           compressedImage.Hash,
		FileName:       in.FileName,
		OriginFileSize: compressedImage.OriSize,
		FileSize:       compressedImage.Size,
		OriginType:     in.FileType,
		FileType:       "webp",
		OriginFilePath: urls.OriFilePath,
		FilePath:       urls.FilePath,
		UserId:         in.UserId,
		Type:           in.Type,
		Status:         false,
		W:              int(in.W),
		H:              int(in.H),
	}); err != nil {
		return nil, err
	}

	return &upload.FileUploadRes{
		FileName:    in.FileName,
		FilePath:    urls.FilePath,
		OriFilePath: urls.OriFilePath,
	}, nil
}

func (l *FileUploadLogic) InsertToMySql(uploadParams *models.Upload) error {
	var uploadModel models.Upload
	if err := l.svcCtx.
		DB.
		Model(&models.Upload{}).
		Where("hash = ? and type = ?", uploadParams.Hash, uploadParams.Type).
		First(&uploadModel).
		Error; errors.As(err, &gorm.ErrRecordNotFound) {
		if err = l.svcCtx.
			DB.
			Model(&models.Upload{}).
			Create(uploadParams).
			Error; err != nil {
			return err
		}
	}

	return nil
}
