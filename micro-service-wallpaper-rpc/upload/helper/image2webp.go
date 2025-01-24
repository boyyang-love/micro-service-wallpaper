package helper

import (
	"bytes"
	"crypto/md5"
	"errors"
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

	if err = webpbin.NewCWebP().
		Quality(quality).
		Input(bytes.NewReader(fileBytes)).
		Output(buf).
		Run(); err != nil {
		return nil, errors.New("webp compress failed")
	}

	reader := bytes.NewReader(buf.Bytes())
	oriRender := bytes.NewReader(oriBuf.Bytes())
	
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
