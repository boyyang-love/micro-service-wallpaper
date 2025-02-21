package helper

import (
	"crypto/md5"
	"fmt"
)

func ImageFileHashByBytes(img []byte) (hash string) {
	return fmt.Sprintf("%x", md5.Sum(img))
}
