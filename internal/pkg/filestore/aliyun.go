package filestore

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type AliYunOssOptions struct {
	Client          interface{}
	Endpoint        string
	AccessKeyId     string
	AccessKeySecret string
	BucketName      string
}

func (e *AliYunOssOptions) NewOptions() error {
	client, err := oss.New(e.Endpoint, e.AccessKeyId, e.AccessKeySecret)
	if err != nil {
		return err
	}
	e.Client = client
	return nil
}

func (e *AliYunOssOptions) UpLoad(pathName string, localFile string) error {
	//获取连接
	val,ok := e.Client.(*oss.Client)
	if !ok {
		err := e.NewOptions()
		if err != nil {
			return err
		}
		val = e.Client.(*oss.Client)
	}
	// 获取存储空间。
	bucket, err := val.Bucket(e.BucketName)
	if err != nil {
		return err
	}
	// 设置分片大小为100 KB，指定分片上传并发数为3，并开启断点续传上传。
	// 其中<yourObjectName>与objectKey是同一概念，表示断点续传上传文件到OSS时需要指定包含文件后缀在内的完整路径，例如abc/efg/123.jpg。
	// "LocalFile"为filePath，100*1024为partSize。
	err = bucket.UploadFile(pathName, localFile, 100*1024,
		oss.Routines(3), oss.Checkpoint(true, ""))
	if err != nil {
		return err
	}
	return nil
}
