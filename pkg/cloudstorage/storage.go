package cloudstorage

import (
	"mime/multipart"
)

type Storage interface {
	// 上传对象-通过本地文件
	UploadObject(objectKey string, localFile string, publicRead bool) error

	// 上传对象-通过文件内容
	UploadByContent(objectKey string, content []byte, publicRead bool) error

	// 上传对象-通过multipart.FileHeader
	UploadByMultipartFileHeader(objectKey string, fileHeader *multipart.FileHeader, publicRead bool) error

	// 获取对象
	GetObject(objectKey string) ([]byte, error)

	// 复制对象
	CopyObject(srcObjectKey string, destObjectKey string, publicRead bool) error

	// 删除单个对象
	DeleteObject(objectKey string) error

	//删除多个对象
	DeleteObjects(objectKeys []string) error
}
