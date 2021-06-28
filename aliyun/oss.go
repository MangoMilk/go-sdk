package aliyun

import (
	"bytes"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

type OssConfig struct {
	AccessKeyID     string
	AccessKeySecret string
	Endpoint        string
	Bucket          string
}

type OssStore struct {
	conf  *OssConfig
	store *oss.Client
}

func NewOssStore(conf *OssConfig) (*OssStore, error) {
	cli, err := oss.New(conf.Endpoint, conf.AccessKeyID, conf.AccessKeySecret)
	if err != nil {
		return nil, err
	}

	return &OssStore{
		conf:  conf,
		store: cli,
	}, nil
}

// 上传字符串
func (s *OssStore) UploadText(content, remotePath string) error {
	bucket, getBucketErr := s.store.Bucket(s.conf.Bucket)
	if getBucketErr != nil {
		return getBucketErr
	}

	// 指定存储类型为标准存储，缺省也为标准存储。
	//storageType := oss.ObjectStorageClass(oss.StorageStandard)
	// 指定访问权限为公共读，缺省为继承bucket的权限。
	//objectAcl := oss.ObjectACL(oss.ACLPublicRead)

	return bucket.PutObject(s.fixPathPrefix(remotePath), strings.NewReader(content)) //, storageType, objectAcl)
}

// 上传Byte数组
func (s *OssStore) UploadBytes(content []byte, remotePath string) error {
	bucket, getBucketErr := s.store.Bucket(s.conf.Bucket)
	if getBucketErr != nil {
		return getBucketErr
	}

	return bucket.PutObject(s.fixPathPrefix(remotePath), bytes.NewReader(content))
}

// 上传本地文件
func (s *OssStore) UploadFile(localPath, remotePath string) error {
	bucket, getBucketErr := s.store.Bucket(s.conf.Bucket)
	if getBucketErr != nil {
		return getBucketErr
	}

	return bucket.PutObjectFromFile(s.fixPathPrefix(remotePath), localPath)
}

// 上传文件流
func (s *OssStore) UploadStream(fd *os.File, remotePath string) error {
	bucket, getBucketErr := s.store.Bucket(s.conf.Bucket)
	if getBucketErr != nil {
		return getBucketErr
	}

	return bucket.PutObject(s.fixPathPrefix(remotePath), fd)
}

// 删除单个文件
func (s *OssStore) DeleteOne(remotePath string) error {
	bucket, getBucketErr := s.store.Bucket(s.conf.Bucket)
	if getBucketErr != nil {
		return getBucketErr
	}

	return bucket.DeleteObject(s.fixPathPrefix(remotePath))
}

// 删除多个文件
func (s *OssStore) DeleteBatch(remotePaths []string) error {
	bucket, getBucketErr := s.store.Bucket(s.conf.Bucket)
	if getBucketErr != nil {
		return getBucketErr
	}

	for k, v := range remotePaths {
		remotePaths[k] = s.fixPathPrefix(v)
	}

	_, err := bucket.DeleteObjects(remotePaths, oss.DeleteObjectsQuiet(true))
	//delRes
	return err
}

func (s *OssStore) GetContent(remotePath string) ([]byte, error) {
	bucket, getBucketErr := s.store.Bucket(s.conf.Bucket)
	if getBucketErr != nil {
		return nil, getBucketErr
	}

	body, err := bucket.GetObject(s.fixPathPrefix(remotePath))
	if err != nil {
		return nil, err
	}

	// 数据读取完成后，获取的流必须关闭，否则会造成连接泄漏，导致请求无连接可用，程序无法正常工作。
	defer body.Close()

	return ioutil.ReadAll(body)
}

// 下载文件到缓存
func (s *OssStore) GetContentToBuf(remotePath string, buf *bytes.Buffer) error {
	bucket, getBucketErr := s.store.Bucket(s.conf.Bucket)
	if getBucketErr != nil {
		return getBucketErr
	}

	body, err := bucket.GetObject(s.fixPathPrefix(remotePath))
	if err != nil {
		return err
	}

	defer body.Close()

	if _, copyErr := io.Copy(buf, body); copyErr != nil {
		return copyErr
	}

	return nil
}

// 下载文件到文件流
func (s *OssStore) GetContentToStream(remotePath string, fd *os.File) error {
	bucket, getBucketErr := s.store.Bucket(s.conf.Bucket)
	if getBucketErr != nil {
		return getBucketErr
	}

	body, err := bucket.GetObject(s.fixPathPrefix(remotePath))
	if err != nil {
		return err
	}

	defer body.Close()

	if _, copyErr := io.Copy(fd, body); copyErr != nil {
		return copyErr
	}

	if _, seekErr := fd.Seek(0, 0); seekErr != nil {
		return seekErr
	}

	return nil
}

// 下载文件到本地
func (s *OssStore) Download(remotePath, localPath string) error {
	bucket, getBucketErr := s.store.Bucket(s.conf.Bucket)
	if getBucketErr != nil {
		return getBucketErr
	}

	return bucket.GetObjectToFile(s.fixPathPrefix(remotePath), localPath)
}

// 重置路径前缀
func (s *OssStore) fixPathPrefix(path string) string {
	prefix := path[0:1]
	if prefix == "/" {
		path = path[1:]
	}
	return path
}
