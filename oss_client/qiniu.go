package oss_client

import (
	"context"
	"errors"
	"io"

	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/spf13/cast"
)

type qiniuClient struct {
	opts OssConfig
}

func NewQiniuClient(opts OssConfig) OssClientIface {
	return &qiniuClient{opts: opts}
}

// Token 获取token
func (p *qiniuClient) Token(ctx context.Context, bucket string, objectKey string, opts ...Option) (*Token, error) {
	option := ForOption(opts...)
	token := p.getToken(bucket, objectKey, option.TokenExpiry)
	return &Token{Token: token}, nil
}

// ObjectMeta 获取对象基本信息
func (p *qiniuClient) ObjectMeta(ctx context.Context, bucket string, objectKey string, opts ...Option) (*ObjectMeta, error) {
	mac := auth.New(p.opts.AccessKey, p.opts.AccessSecret)

	cfg := storage.Config{UseHTTPS: false}
	bucketManager := storage.NewBucketManager(mac, &cfg)
	fileInfo, err := bucketManager.Stat(bucket, objectKey)
	if err != nil {
		return nil, err
	}

	out := &ObjectMeta{
		Hash:     fileInfo.Hash,
		FileSize: fileInfo.Fsize,
		PutTime:  fileInfo.PutTime,
		Type:     fileInfo.Type,
	}
	return out, nil
}

// PutObject 上传对象
func (p *qiniuClient) PutObject(ctx context.Context, bucket string, objectKey string, data io.Reader, size int64, opts ...Option) (string, error) {
	upToken := p.getToken(bucket, objectKey, 3600)
	cfg := storage.Config{UseCdnDomains: true, Zone: &storage.ZoneHuanan}

	// 构建表单上传的对象
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}

	err := formUploader.Put(context.Background(), &ret, upToken, objectKey, data, size, nil)
	if err != nil {
		return "", err
	}

	return p.opts.DownloadUrl + "/" + ret.Key, nil
}

func (p *qiniuClient) PutObjectFromFile(ctx context.Context, bucket string, objectKey string, filePath string, opts ...Option) (string, error) {
	upToken := p.getToken(bucket, objectKey, 3600)
	cfg := storage.Config{UseCdnDomains: true}

	// 构建表单上传的对象
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}

	err := formUploader.PutFile(ctx, &ret, upToken, objectKey, filePath, nil)
	if err != nil {
		return "", err
	}

	return p.opts.DownloadUrl + "/" + ret.Key, nil
}

// GetObject 获取对象
func (p *qiniuClient) GetObject(ctx context.Context, bucket string, objectKey string, opts ...Option) ([]byte, error) {

	return nil, nil
}

// BatchCopy 批量复制,
func (p *qiniuClient) BatchCopy(ctx context.Context, bucket string, copyKeys map[string]string, opts ...Option) error {
	mac := auth.New(p.opts.AccessKey, p.opts.AccessSecret)

	cfg := storage.Config{
		// 是否使用https域名进行资源管理
		UseHTTPS: false,
	}
	cfg.Zone = &storage.ZoneHuanan
	bucketManager := storage.NewBucketManager(mac, &cfg)

	//每个batch的操作数量不可以超过1000个，如果总数量超过1000，需要分批发送
	copyOps := make([]string, 0, len(copyKeys))
	for srcKey, destKey := range copyKeys {
		copyOps = append(copyOps, storage.URICopy(bucket, srcKey, bucket, destKey, true))
	}

	rets, err := bucketManager.Batch(copyOps)
	if err != nil {
		return err
	}

	var errStr string
	for _, ret := range rets {
		if ret.Code != 200 {
			errStr = ret.Data.Error
			break
		}
	}

	if errStr != "" {
		return errors.New(errStr)
	}

	return nil
}

func (p *qiniuClient) Copy(ctx context.Context, bucket string, srcObject string, dstObject string, opts ...Option) error {
	mac := auth.New(p.opts.AccessKey, p.opts.AccessSecret)

	cfg := storage.Config{
		// 是否使用https域名进行资源管理
		UseHTTPS: false,
	}
	cfg.Zone = &storage.ZoneHuanan
	bucketManager := storage.NewBucketManager(mac, &cfg)

	return bucketManager.Copy(bucket, srcObject, bucket, dstObject, false)
}

func (p *qiniuClient) ListObjects(ctx context.Context, bucket string, opts ...Option) (ListObjectsResult, error) {
	op := &Options{}
	// 解析参数信息
	for _, o := range opts {
		o(op)
	}
	mac := auth.New(p.opts.AccessKey, p.opts.AccessSecret)

	cfg := storage.Config{UseHTTPS: false}
	bucketManager := storage.NewBucketManager(mac, &cfg)

	outResult := ListObjectsResult{
		Prefix:    op.Prefix,
		Marker:    op.Marker,
		Delimiter: op.Delimiter,
	}

	items, commonPrefixes, nextMarker, hasNext, err := bucketManager.ListFiles(bucket, op.Prefix, op.Delimiter, op.Marker, op.Limit)
	if err != nil {
		return outResult, err
	}

	outResult.IsTruncated = hasNext
	outResult.NextMarker = nextMarker
	outResult.CommonPrefixes = commonPrefixes
	for _, v := range items {
		outResult.Objects = append(outResult.Objects, ObjectProperties{
			Key:  v.Key,
			Type: cast.ToString(v.Type),
			Size: v.Fsize,
			ETag: v.Hash,
		})
	}

	return outResult, nil
}

// MakePrivateURL 创建私有URL
func (p *qiniuClient) MakePrivateURL(ctx context.Context, bucketName string, key string, deadline int64, opts ...Option) (string, error) {
	return "", nil
}

func (p *qiniuClient) getToken(bucket string, objectKey string, expiry int) string {
	putPolicy := storage.PutPolicy{
		Scope:   bucket + ":" + objectKey,
		Expires: uint64(expiry),
	}

	mac := qbox.NewMac(p.opts.AccessKey, p.opts.AccessSecret)

	return putPolicy.UploadToken(mac)
}

func (p *qiniuClient) Options() OssConfig {
	return p.opts
}
