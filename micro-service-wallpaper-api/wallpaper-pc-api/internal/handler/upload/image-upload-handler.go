package upload

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/boyyang-love/micro-service-wallpaper-api/helper"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/logic/upload"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/svc"
	"github.com/boyyang-love/micro-service-wallpaper-api/internal/types"
	"github.com/boyyang-love/micro-service-wallpaper-rpc/upload/models"
	"github.com/boyyang-love/micro-service-wallpaper-rpc/upload/uploadclient"
	"github.com/zeromicro/go-zero/rest/httpx"
	"gorm.io/gorm"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"
)

func ImageUploadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ImageUploadReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		// userid
		userid := r.Context().Value("Id")

		file, fileHeader, err := r.FormFile("file")

		img, imgType, err := image.Decode(file)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		originBuffer := new(bytes.Buffer)
		switch imgType {
		case "png":
			if err = png.Encode(originBuffer, img); err != nil {
				httpx.ErrorCtx(r.Context(), w, err)
				return
			}
		case "jpeg", "jpg":
			if err = jpeg.Encode(originBuffer, img, nil); err != nil {
				httpx.ErrorCtx(r.Context(), w, err)
				return
			}
		}
		path := fmt.Sprintf("%s/%s/%s/%s", req.RootDir, req.Dir, userid, "compressed")
		oriPath := fmt.Sprintf("%s/%s/%s/%s", req.RootDir, req.Dir, userid, "original")
		hash := helper.ImageFileHashByBytes(originBuffer.Bytes())

		isHave, info := IsHave(svcCtx.DB, hash, path)
		if isHave {
			resp := &types.ImageUploadRes{
				Base: types.Base{
					Code: 1,
					Msg:  "图片上传成功",
				},
				Data: types.ImageUploadResdata{
					FileName:   info.FileName,
					Path:       info.FilePath,
					OriginPath: info.OriginFilePath,
				},
			}

			httpx.OkJsonCtx(r.Context(), w, resp)
			return
		}

		l := upload.NewImageUploadLogic(r.Context(), svcCtx)
		resp, err := l.ImageUpload(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {

			fileUpload, err := svcCtx.UploadService.ImageUpload(
				r.Context(),
				&uploadclient.ImageUploadReq{
					File:       originBuffer.Bytes(),
					FileName:   fileHeader.Filename,
					Path:       path,
					OriPath:    oriPath,
					Quality:    uint32(req.Quality),
					BucketName: req.BucketName,
				},
			)
			if err != nil {
				httpx.ErrorCtx(r.Context(), w, err)
				return
			}

			err := InsertToDB(svcCtx.DB, &models.Upload{
				Hash:           fileUpload.Data.ETag,
				FileName:       fileHeader.Filename,
				OriginFileSize: int64(fileUpload.Data.OriSize),
				FileSize:       int64(fileUpload.Data.Size),
				OriginType:     imgType,
				FileType:       ".webp",
				OriginFilePath: fileUpload.Data.OriPath,
				FilePath:       fileUpload.Data.Path,
				UserId:         userid,
				Type:           req.Type,
				Status:         false,
				W:              img.Bounds().Dx(),
				H:              img.Bounds().Dy(),
			})
			if err != nil {
				httpx.ErrorCtx(r.Context(), w, err)
				return
			}

			resp = &types.ImageUploadRes{
				Base: types.Base{
					Code: int(fileUpload.Base.Code),
					Msg:  fileUpload.Base.Msg,
				},
				Data: types.ImageUploadResdata{
					FileName:   req.FileName,
					Path:       fileUpload.Data.Path,
					OriginPath: fileUpload.Data.OriPath,
				},
			}

			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

func IsHave(db *gorm.DB, imageHash string, minioPath string) (isHave bool, info *models.Upload) {
	if err := db.
		Model(&models.Upload{}).
		Select("hash", "file_name", "file_path", "origin_file_path").
		Where("hash = ? and file_path = ?", imageHash, minioPath).
		First(&info).
		Error; errors.As(err, &gorm.ErrRecordNotFound) {
		return false, info
	}

	return true, info
}

func InsertToDB(db *gorm.DB, info *models.Upload) (err error) {
	if err := db.
		Model(&models.Upload{}).
		Create(&info).
		Error; err != nil {
		return err
	}
}
