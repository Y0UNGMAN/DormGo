package utils

import (
	"fmt"
	"mime/multipart"
	"path"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

func UploadFile(file *multipart.FileHeader, dirName string) (string, error) {
	//1. 创建oss client客户端
	client, err := oss.New(
		viper.GetString("oss.Endpoint"),
		viper.GetString("oss.AccessKeyId"),
		viper.GetString("oss.AccessKeySecret"),
	)
	if err != nil {
		return "", err
	}

	//2. 获取Bucket实例
	bucket, err := client.Bucket(viper.GetString("oss.BucketName"))
	if err != nil {
		return "", err
	}
	// 3. 打开上传的文件流
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// 4. 构造唯一文件名 (防止用户文件名冲突)
	// 最终路径类似: images/20231127/550e8400-e29b-41d4-a716-446655440000.jpg
	ext := path.Ext(file.Filename)            // 获取后缀，如 .jpg
	datePath := time.Now().Format("20060102") // 按日期分文件夹
	filename := fmt.Sprintf("%s/%s/%s%s", dirName, datePath, uuid.New().String(), ext)

	err = bucket.PutObject(filename, src)
	if err != nil {
		return "", err
	}

	return viper.GetString("oss.BaseURL") + "/" + filename, nil
}
