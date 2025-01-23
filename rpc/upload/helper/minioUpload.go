package helper

import (
	"bytes"
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
)

type MinioFileUploadParams struct {
	Ctx         context.Context
	MinioClient *minio.Client
	Buf         *bytes.Buffer
	OriBuf      *bytes.Buffer
	FileHash    string
	Folder      string
	OriFolder   string
	Filename    string
	BucketName  string
}

type MinioFileReturnPaths struct {
	FilePath    string
	OriFilePath string
}

func MinioFileUpload(params *MinioFileUploadParams) (urls *MinioFileReturnPaths, err error) {

	fileName := fmt.Sprintf("%s-%s.webp", FileNameNoExt(params.Filename), params.FileHash)
	filePath := fmt.Sprintf("%s/%s", params.Folder, fileName)
	oriFilePath := fmt.Sprintf("%s/%s", params.OriFolder, params.Filename)

	reader := bytes.NewReader(params.Buf.Bytes())
	oriReader := bytes.NewReader(params.OriBuf.Bytes())

	_, err = params.MinioClient.PutObject(
		params.Ctx,
		params.BucketName,
		filePath,
		reader,
		reader.Size(),
		minio.PutObjectOptions{},
	)

	_, err = params.MinioClient.PutObject(
		params.Ctx,
		params.BucketName,
		oriFilePath,
		oriReader,
		oriReader.Size(),
		minio.PutObjectOptions{},
	)

	if err == nil {
		return &MinioFileReturnPaths{
			FilePath:    filePath,
			OriFilePath: oriFilePath,
		}, nil
	} else {
		return nil, err
	}
}
