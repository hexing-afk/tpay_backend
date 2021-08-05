package cloudstorage

import (
	"bytes"
	"errors"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Storage struct {
	Region          string
	AccessKeyId     string
	AccessKeySecret string
	Bucket          string
}

func NewS3Storage(accessKeyId, accessKeySecret, region, bucket string) Storage {
	return &S3Storage{
		Region:          region,
		AccessKeyId:     accessKeyId,
		AccessKeySecret: accessKeySecret,
		Bucket:          bucket,
	}
}

// 获取上传的session
func (s *S3Storage) GetSession() (*session.Session, error) {
	return session.NewSession(&aws.Config{
		Region:      aws.String(s.Region),
		Credentials: credentials.NewStaticCredentials(s.AccessKeyId, s.AccessKeySecret, ""),
	})
}

func (s *S3Storage) UploadObject(objectKey string, localFile string, publicRead bool) error {
	// 获取session
	sess, err := s.GetSession()
	if err != nil {
		return err
	}

	// 打开文件
	file, err := os.Open(localFile)
	if err != nil {
		return err
	}
	defer file.Close()

	// 获取文件信息
	fileInfo, ferr := file.Stat()
	if ferr != nil {
		return ferr
	}

	// 将文件读入buffer
	buffer := make([]byte, fileInfo.Size())
	file.Read(buffer)

	input := &s3.PutObjectInput{
		Bucket:        aws.String(s.Bucket),
		Key:           aws.String(objectKey),
		Body:          bytes.NewReader(buffer),
		ContentLength: aws.Int64(fileInfo.Size()),
		ContentType:   aws.String(http.DetectContentType(buffer)),
	}

	if publicRead { // 可以公开访问
		input.ACL = aws.String(s3.ObjectCannedACLPublicRead)
	}

	_, rerr := s3.New(sess).PutObject(input)

	return rerr
}

func (s *S3Storage) UploadByContent(objectKey string, content []byte, publicRead bool) error {
	// 获取session
	sess, err := s.GetSession()
	if err != nil {
		return err
	}

	input := &s3.PutObjectInput{
		Bucket:        aws.String(s.Bucket),
		Key:           aws.String(objectKey),
		Body:          bytes.NewReader(content),
		ContentLength: aws.Int64(int64(len(content))),
		ContentType:   aws.String(http.DetectContentType(content)),
	}

	if publicRead { // 可以公开访问
		input.ACL = aws.String(s3.ObjectCannedACLPublicRead)
	}

	_, rerr := s3.New(sess).PutObject(input)
	return rerr
}

func (s *S3Storage) UploadByMultipartFileHeader(objectKey string, fileHeader *multipart.FileHeader, publicRead bool) error {
	// 获取session
	sess, err := s.GetSession()
	if err != nil {
		return err
	}

	// 获取文件内容
	fileContent, err := fileHeader.Open()
	if err != nil {
		return err
	}

	// 读取文件内容
	bytesContent, err := ioutil.ReadAll(fileContent)
	if err != nil {
		return err
	}

	input := &s3.PutObjectInput{
		Bucket:        aws.String(s.Bucket),
		Key:           aws.String(objectKey),
		Body:          bytes.NewReader(bytesContent),
		ContentLength: aws.Int64(fileHeader.Size),
		ContentType:   aws.String(http.DetectContentType(bytesContent)),
	}

	if publicRead { // 可以公开访问
		input.ACL = aws.String(s3.ObjectCannedACLPublicRead)
	}

	_, rerr := s3.New(sess).PutObject(input)

	return rerr
}

func (s *S3Storage) GetObject(objectKey string) ([]byte, error) {
	// 获取session
	sess, err := s.GetSession()
	if err != nil {
		return nil, err
	}

	input := &s3.GetObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(objectKey),
	}

	result, err := s3.New(sess).GetObject(input)
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(result.Body)
}

func (s *S3Storage) CopyObject(srcObjectKey string, destObjectKey string, publicRead bool) error {
	// 获取session
	sess, err := s.GetSession()
	if err != nil {
		return err
	}

	input := &s3.CopyObjectInput{
		Bucket:     aws.String(s.Bucket),
		CopySource: aws.String("/" + s.Bucket + "/" + srcObjectKey),
		Key:        aws.String(destObjectKey),
	}

	if publicRead { // 可以公开访问
		input.ACL = aws.String(s3.ObjectCannedACLPublicRead)
	}

	_, rerr := s3.New(sess).CopyObject(input)
	return rerr
}

func (s *S3Storage) DeleteObject(objectKey string) error {
	// 获取session
	sess, err := s.GetSession()
	if err != nil {
		return err
	}

	_, rerr := s3.New(sess).DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(objectKey),
	})

	return rerr
}

func (s *S3Storage) DeleteObjects(objectKeys []string) error {
	// 获取session
	sess, err := s.GetSession()
	if err != nil {
		return err
	}
	if len(objectKeys) == 0 {
		return errors.New("objectKeys不能为空")
	}

	var objList []*s3.ObjectIdentifier
	for _, v := range objectKeys {
		objList = append(objList, &s3.ObjectIdentifier{Key: aws.String(v)})
	}

	_, err = s3.New(sess).DeleteObjects(&s3.DeleteObjectsInput{
		Bucket: aws.String(s.Bucket),
		Delete: &s3.Delete{
			Objects: objList,
		},
	})
	return err
}
