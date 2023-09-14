package oss_client

import (
	"context"
	"io"
)

// OssClientIface oss客户端
type OssClientIface interface {
	// Token 获取token
	Token(ctx context.Context, bucket string, objectKey string, opts ...Option) (*Token, error)
	// ObjectMeta 获取对象基本信息
	ObjectMeta(ctx context.Context, bucket string, objectKey string, opts ...Option) (*ObjectMeta, error)
	// PutObject 上传对象
	PutObject(ctx context.Context, bucket string, objectKey string, data io.Reader, size int64, opts ...Option) (string, error)
	// PutObjectFromFile 从本地文件上传对象
	PutObjectFromFile(ctx context.Context, bucket string, objectKey string, filePath string, opts ...Option) (string, error)
	// BatchCopy 批量复制
	BatchCopy(ctx context.Context, bucket string, copyKeys map[string]string, opts ...Option) error
	// GetObject 获取对象
	GetObject(ctx context.Context, bucket string, objectKey string, opts ...Option) ([]byte, error)
	// Copy 复制对象
	Copy(ctx context.Context, bucket string, srcObject string, dstObject string, opts ...Option) error
	// ListObjects 获取对象信息列表
	ListObjects(ctx context.Context, bucket string, opts ...Option) (ListObjectsResult, error)
	// MakePrivateURL 创建私有URL
	MakePrivateURL(ctx context.Context, bucket string, key string, deadline int64, opts ...Option) (string, error)
}

type OssConfig struct {
	AccessKey       string // key
	AccessSecret    string // 密钥
	Endpoint        string // 端点
	DownloadUrl     string // 下载地址
	RoleArn         string // 角色
	RoleSessionName string // 角色名称
	RegionId        string // 区域
}

type Token struct {
	Token           string `json:"token"`             // token信息
	AccessKey       string `json:"access_key"`        // key
	AccessKeySecret string `json:"access_key_secret"` // 密钥
	Expiration      int64  `json:"expiration"`        // 到期时间，时间戳
	Bucket          string `json:"bucket"`            // bucket
	Endpoint        string `json:"endpoint"`          // 端点
}

type ObjectMeta struct {
	Hash     string `json:"hash"`      // 文件hash
	FileSize int64  `json:"file_size"` // 文件大小
	PutTime  int64  `json:"putTime"`   // 上传时间
	Type     int    `json:"type"`      // 类型
}

type ListObjectsResult struct {
	Prefix         string             `json:"prefix"`
	Marker         string             `json:"marker"`
	Delimiter      string             `json:"delimiter"`
	IsTruncated    bool               `json:"is_truncated"`
	NextMarker     string             `json:"next_marker"`
	Objects        []ObjectProperties `json:"objects"`
	CommonPrefixes []string           `json:"common_prefixes"`
}

type ObjectProperties struct {
	Key  string `json:"key"`  // Object key
	Type string `json:"type"` // Object type
	Size int64  `json:"size"` // Object size
	ETag string `json:"etag"` // Object etag
}
