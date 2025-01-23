package helper

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"github.com/nickalie/go-webpbin"
)

type CompressedImage struct {
	Hash    string
	Buf     *bytes.Buffer
	OriBuf  *bytes.Buffer
	Size    int64
	OriSize int64
}

func Image2Webp(fileBytes []byte, quality uint) (compressedImage *CompressedImage, err error) {

	buf := new(bytes.Buffer)
	oriBuf := bytes.NewBuffer(fileBytes)
	reader := new(bytes.Reader)
	oriRender := new(bytes.Reader)

	if err = webpbin.NewCWebP().
		Quality(quality).
		Input(bytes.NewReader(fileBytes)).
		Output(buf).
		Run(); err != nil {
		return nil, err
	}

	_, err = reader.Read(buf.Bytes())
	if err != nil {
		return nil, err
	}

	_, err = oriRender.Read(oriBuf.Bytes())
	if err != nil {
		return nil, err
	}

	return &CompressedImage{
		Hash:    MakeHashByBytes(fileBytes),
		Buf:     buf,
		OriBuf:  oriBuf,
		Size:    reader.Size(),
		OriSize: oriRender.Size(),
	}, nil

}

func MakeHashByBytes(bytes []byte) (hash string) {
	return fmt.Sprintf("%x", md5.Sum(bytes))
}
