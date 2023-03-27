package util

import (
	"os"
)

func HasUTF8Dom(bytes []byte) bool {
	return len(bytes) >= 3 && bytes[0] == 0xEF && bytes[1] == 0xBB && bytes[2] == 0xBF
}

// IsSrtFile 是否为srt文件
func IsSrtFile(strFilepath string) bool {
	if len(strFilepath) <= 4 {
		return false
	}
	return strFilepath[len(strFilepath)-3:] == "srt"
}

// IsFileOrDirExisted 检测文件或路径是否存在
func IsFileOrDirExisted(strFilepath string) error {
	fd, err := os.Open(strFilepath)
	defer func() {
		if fd != nil {
			_ = fd.Close()
		}
	}()
	if err != nil {
		return err
	}
	return nil
}
