package oss

import (
	"errors"
	"os"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/gogf/gf/util/guid"
)

var Service = ossService{}

type ossService struct{}

type BaseParams struct {
	Endpoint        string
	AccessKeyId     string
	AccessKeySecret string
	UseCname        bool // 是否绑定自定义域名
	Bucket          string
	ObjectKey       string // 目标地址
}
type UploadParams struct {
	*BaseParams
	ForbidOverWrite bool
	ACL             string   // 文件权限
	LocalFilePath   string   // （可选）本地文件地址
	FileStream      *os.File // （可选）文件流
}
type DownloadParams struct {
	*BaseParams
	LocalFilePath string   // 本地文件地址
	FileStream    *os.File // （可选）文件流
}

type SignUrlParams struct {
	*BaseParams
	ACL         string // (可选)上传时权限控制参数
	Timeout     int64
	ContentType string
}

type ListObjectsParams struct {
	*BaseParams
	Page    int64
	MaxKeys int
	Marker  string
	Prefix  string
}

type DeleteObjectsParams struct {
	*BaseParams
	ObjectKeys []string
}

func (s *ossService) InitClient(Endpoint string, AccessKeyId string, AccessKeySecret string, UseCname bool) (client *oss.Client, err error) {
	client, err = oss.New(Endpoint, AccessKeyId, AccessKeySecret, oss.UseCname(UseCname), oss.Timeout(30, 120))
	return
}

func (s *ossService) InitBucket(Endpoint string, AccessKeyId string, AccessKeySecret string, UseCname bool, Bucket string) (bucket *oss.Bucket, err error) {
	client, err := s.InitClient(Endpoint, AccessKeyId, AccessKeySecret, UseCname)
	if err != nil {
		return nil, errors.New("初始化OSS失败，请检查密钥或配置，错误信息：" + err.Error())
	}
	bucket, err = client.Bucket(Bucket)
	if err != nil {
		return bucket, errors.New("设置OSS存储空间失败，错误信息：" + err.Error())
	}
	return
}

// UploadLocalFile 上传本地文件
func (s *ossService) UploadLocalFile(params *UploadParams) (err error) {
	bucket, err := s.InitBucket(params.Endpoint, params.AccessKeyId, params.AccessKeySecret, params.UseCname, params.Bucket)
	if err != nil {
		return err
	}
	err = bucket.PutObjectFromFile(params.ObjectKey, params.LocalFilePath, oss.ForbidOverWrite(params.ForbidOverWrite), oss.ObjectACL(oss.ACLType(params.ACL)))
	if err != nil {
		return errors.New("上传文件失败，错误信息：" + err.Error())
	}
	return nil
}

// UploadFileStream 上传文件流
func (s *ossService) UploadFileStream(params *UploadParams) (err error) {
	bucket, err := s.InitBucket(params.Endpoint, params.AccessKeyId, params.AccessKeySecret, params.UseCname, params.Bucket)
	if err != nil {
		return err
	}
	err = bucket.PutObject(params.ObjectKey, params.FileStream, oss.ForbidOverWrite(params.ForbidOverWrite), oss.ObjectACL(oss.ACLType(params.ACL)))
	if err != nil {
		return errors.New("上传文件失败，错误信息：" + err.Error())
	}
	return nil
}

// GetObjectToFile 下载文件到本地
func (s *ossService) GetObjectToFile(params *DownloadParams) (err error) {
	bucket, err := s.InitBucket(params.Endpoint, params.AccessKeyId, params.AccessKeySecret, params.UseCname, params.Bucket)
	if err != nil {
		return err
	}
	err = bucket.GetObjectToFile(params.ObjectKey, params.LocalFilePath)
	if err != nil {
		return errors.New("文件下载失败，错误信息：" + err.Error())
	}
	return nil
}

// CreateSignedPutUrl 创建临时上传签名URL
func (s *ossService) CreateSignedPutUrl(params *SignUrlParams) (signedURL string, err error) {
	bucket, err := s.InitBucket(params.Endpoint, params.AccessKeyId, params.AccessKeySecret, params.UseCname, params.Bucket)
	if err != nil {
		return "", err
	}
	options := []oss.Option{
		oss.ObjectACL(oss.ACLType(params.ACL)),
		oss.ContentType(params.ContentType),
		oss.ResponseContentType("json"),
	}
	signedURL, err = bucket.SignURL(params.ObjectKey, oss.HTTPPut, params.Timeout, options...)
	if err != nil {
		return "", errors.New("OSS签名授权失败，错误信息：" + err.Error())
	}
	return
}

// CreateSignedGetUrl 创建临时下载/访问签名URL
func (s *ossService) CreateSignedGetUrl(params *SignUrlParams) (signedURL string, err error) {
	bucket, err := s.InitBucket(params.Endpoint, params.AccessKeyId, params.AccessKeySecret, params.UseCname, params.Bucket)
	if err != nil {
		return "", err
	}
	signedURL, err = bucket.SignURL(params.ObjectKey, oss.HTTPGet, params.Timeout)
	if err != nil {
		return "", errors.New("OSS签名授权失败，错误信息：" + err.Error())
	}
	return
}

// ListBucketObjects 分页列举OSS文件
func (s *ossService) ListBucketObjects(params *ListObjectsParams) (lsRes oss.ListObjectsResult, err error) {
	bucket, err := s.InitBucket(params.Endpoint, params.AccessKeyId, params.AccessKeySecret, params.UseCname, params.Bucket)
	if err != nil {
		return
	}
	marker := oss.Marker(params.Marker)
	lsRes, err = bucket.ListObjects(oss.MaxKeys(params.MaxKeys), marker, oss.Prefix(params.Prefix))
	return
}

// DeleteObject 删除单个OSS存储文件
func (s *ossService) DeleteObject(params *BaseParams) (err error) {
	bucket, err := s.InitBucket(params.Endpoint, params.AccessKeyId, params.AccessKeySecret, params.UseCname, params.Bucket)
	if err != nil {
		return err
	}
	err = bucket.DeleteObject(params.ObjectKey)
	if err != nil {
		return errors.New("删除OSS文件失败，错误信息：" + err.Error())
	}
	return
}

// DeleteObjects 删除多个OSS存储文件
func (s *ossService) DeleteObjects(params *DeleteObjectsParams) (err error) {
	bucket, err := s.InitBucket(params.Endpoint, params.AccessKeyId, params.AccessKeySecret, params.UseCname, params.Bucket)
	if err != nil {
		return err
	}
	_, err = bucket.DeleteObjects(params.ObjectKeys, oss.DeleteObjectsQuiet(true))
	if err != nil {
		return errors.New("删除OSS文件失败，错误信息：" + err.Error())
	}
	return
}

// RandomObjectKey 随机唯一文件名
func (s *ossService) RandomObjectKey(prefix string, suffix string) string {
	randomKey := guid.S()
	return prefix + randomKey + "." + suffix
}
