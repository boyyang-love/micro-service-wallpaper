package helper

import (
	"path"
	"strings"
)

func FileNameNoExt(fileName string) string {
	return strings.TrimSuffix(path.Base(fileName), path.Ext(fileName))
}
