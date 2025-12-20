package utils

import (
	"fmt"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

// SaveFileLocal 将上传文件保存到本地 static 目录，返回可访问的相对 URL
func SaveFileLocal(file *multipart.FileHeader, dirName string) (string, error) {
	// 确保目录存在: ./static/<dirName>/<yyyymmdd>/
	datePath := time.Now().Format("20060102")
	dstDir := filepath.Join("static", dirName, datePath)
	if err := os.MkdirAll(dstDir, os.ModePerm); err != nil {
		return "", err
	}

	ext := path.Ext(file.Filename)
	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	dstPath := filepath.Join(dstDir, filename)

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	out, err := os.Create(dstPath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	if _, err := src.Seek(0, 0); err != nil {
		// ignore seek error
	}
	if _, err := out.ReadFrom(src); err != nil {
		return "", err
	}

	// 返回可通过静态路由访问的 URL 路径
	urlPath := "/static/" + dirName + "/" + datePath + "/" + filename
	return urlPath, nil
}
