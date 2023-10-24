package tools

import (
	"encoding/base64"
	"os"
	"path/filepath"
)

func GetAbPath() (path string, err error) {
	str, err := os.Executable()
	if err != nil {
		return "", err
	}

	path, err = filepath.EvalSymlinks(filepath.Dir(str))
	if err != nil {
		return "", err
	}

	return path, nil
}

// base编码
func Base64Encode(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

// base64解码
func Base64Decode(str string) string {
	decodeString, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return ""
	}
	return string(decodeString)
}
