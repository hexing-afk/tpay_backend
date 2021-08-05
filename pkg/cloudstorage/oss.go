package cloudstorage

import (
	"bytes"
	"io/ioutil"
	"mime/multipart"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type OssStorage struct {
	Endpoint        string
	AccessKeyId     string
	AccessKeySecret string
	Bucket          string
}

func NewOssStorage(accessKeyId, accessKeySecret, endpoint, bucket string) Storage {
	return &OssStorage{
		Endpoint:        endpoint,
		AccessKeyId:     accessKeyId,
		AccessKeySecret: accessKeySecret,
		Bucket:          bucket,
	}
}

// 获取存储空间
func (o *OssStorage) GetBucket() (*oss.Bucket, error) {
	client, err := oss.New(o.Endpoint, o.AccessKeyId, o.AccessKeySecret)
	if err != nil {
		return nil, err
	}

	// 获取存储空间。
	bucket, err := client.Bucket(o.Bucket)
	if err != nil {
		return nil, err
	}

	return bucket, nil
}

func (o *OssStorage) UploadObject(objectKey string, localFile string, publicRead bool) error {
	bucket, err := o.GetBucket()
	if err != nil {
		return err
	}

	var options []oss.Option
	if publicRead {
		options = append(options, oss.ObjectACL(oss.ACLPublicRead))
	}

	return bucket.PutObjectFromFile(objectKey, localFile, options...)
}

func (o *OssStorage) UploadByContent(objectKey string, content []byte, publicRead bool) error {
	bucket, err := o.GetBucket()
	if err != nil {
		return err
	}

	fd := bytes.NewReader(content)

	var options []oss.Option
	if publicRead {
		options = append(options, oss.ObjectACL(oss.ACLPublicRead))
	}

	return bucket.PutObject(objectKey, fd, options...)
}

func (o *OssStorage) UploadByMultipartFileHeader(objectKey string, fileHeader *multipart.FileHeader, publicRead bool) error {
	bucket, err := o.GetBucket()
	if err != nil {
		return err
	}

	// 读取数据流
	fd, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer fd.Close()

	var options []oss.Option
	if publicRead {
		options = append(options, oss.ObjectACL(oss.ACLPublicRead))
	}

	return bucket.PutObject(objectKey, fd, options...)
}

func (o *OssStorage) GetObject(objectKey string) ([]byte, error) {
	bucket, err := o.GetBucket()
	if err != nil {
		return nil, err
	}

	result, err := bucket.GetObject(objectKey)
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(result)
}

func (o *OssStorage) CopyObject(srcObjectKey string, destObjectKey string, publicRead bool) error {
	bucket, err := o.GetBucket()
	if err != nil {
		return err
	}

	var options []oss.Option
	if publicRead {
		options = append(options, oss.ObjectACL(oss.ACLPublicRead))
	}

	_, rerr := bucket.CopyObject(srcObjectKey, destObjectKey, options...)
	return rerr
}

func (o *OssStorage) DeleteObject(objectKey string) error {
	bucket, err := o.GetBucket()
	if err != nil {
		return err
	}

	return bucket.DeleteObject(objectKey)
}

func (o *OssStorage) DeleteObjects(objectKeys []string) error {
	bucket, err := o.GetBucket()
	if err != nil {
		return err
	}

	_, err = bucket.DeleteObjects(objectKeys)
	return err
}
