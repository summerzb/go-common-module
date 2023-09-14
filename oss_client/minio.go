package oss_client

import (
	"context"
	"io"
	"io/ioutil"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type minioClient struct {
	opts OssConfig
}

func NewMinioClient(opts OssConfig) OssClientIface {
	return &minioClient{opts: opts}
}

// Token 获取token
func (p *minioClient) Token(ctx context.Context, bucket string, objectKey string, opts ...Option) (*Token, error) {
	op := ForOption(opts...)
	cred, err := credentials.NewSTSAssumeRole(p.opts.Endpoint, credentials.STSAssumeRoleOptions{
		AccessKey:       p.opts.AccessKey,    // sts授权，这里是用户名
		SecretKey:       p.opts.AccessSecret, // sts授权，这里是用户密码
		DurationSeconds: op.TokenExpiry,
	})
	if err != nil {
		return nil, err
	}

	value, err := cred.Get()
	if err != nil {
		return nil, err
	}

	return &Token{
		Token:           value.SessionToken,
		AccessKey:       value.AccessKeyID,
		AccessKeySecret: value.SecretAccessKey,
	}, nil
}

// ObjectMeta 获取对象基本信息
func (p *minioClient) ObjectMeta(ctx context.Context, bucket string, objectKey string, opts ...Option) (*ObjectMeta, error) {
	op := ForOption(opts...)
	client, err := p.client(op.SessionToken)
	if err != nil {
		return nil, err
	}

	objInfo, err := client.StatObject(ctx, bucket, objectKey, minio.StatObjectOptions{})
	if err != nil {
		return nil, err
	}
	metaInfo := &ObjectMeta{
		FileSize: objInfo.Size,
		PutTime:  objInfo.LastModified.Unix(),
		Type:     0,
	}
	return metaInfo, nil
}

// PutObject 上传对象
func (p *minioClient) PutObject(ctx context.Context, bucket string, objectKey string, data io.Reader, size int64, opts ...Option) (string, error) {
	op := ForOption(opts...)
	client, err := p.client(op.SessionToken)
	if err != nil {
		return "", err
	}

	_, err = client.PutObject(ctx, bucket, objectKey, data, size, minio.PutObjectOptions{})
	if err != nil {
		return "", err
	}

	return p.opts.DownloadUrl + "/" + objectKey, nil
}

// PutObjectFromFile 从本地文件上传对象
func (p *minioClient) PutObjectFromFile(ctx context.Context, bucket string, objectKey string, filePath string, opts ...Option) (string, error) {
	op := ForOption(opts...)
	client, err := p.client(op.SessionToken)
	if err != nil {
		return "", err
	}

	_, err = client.FPutObject(ctx, bucket, objectKey, filePath, minio.PutObjectOptions{})
	if err != nil {
		return "", err
	}

	return p.opts.DownloadUrl + "/" + objectKey, nil
}

// BatchCopy 批量复制
func (p *minioClient) BatchCopy(ctx context.Context, bucket string, copyKeys map[string]string, opts ...Option) error {
	op := ForOption(opts...)
	client, err := p.client(op.SessionToken)
	if err != nil {
		return err
	}

	for src, dst := range copyKeys {
		srcOpts := minio.CopySrcOptions{
			Bucket: bucket,
			Object: src,
		}

		// Destination object
		dstOpts := minio.CopyDestOptions{
			Bucket: bucket,
			Object: dst,
		}
		// Copy object call
		_, err := client.CopyObject(ctx, dstOpts, srcOpts)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *minioClient) Copy(ctx context.Context, bucketName string, srcObject string, dstObject string, opts ...Option) error {
	op := ForOption(opts...)
	client, err := p.client(op.SessionToken)
	if err != nil {
		return err
	}
	srcOpts := minio.CopySrcOptions{
		Bucket: bucketName,
		Object: srcObject,
	}

	// Destination object
	dstOpts := minio.CopyDestOptions{
		Bucket: bucketName,
		Object: dstObject,
	}
	_, err = client.CopyObject(ctx, dstOpts, srcOpts)
	return err
}

// GetObject 获取对象
func (p *minioClient) GetObject(ctx context.Context, bucket string, objectKey string, opts ...Option) ([]byte, error) {
	op := ForOption(opts...)
	client, err := p.client(op.SessionToken)
	if err != nil {
		return nil, err
	}

	object, err := client.GetObject(ctx, bucket, objectKey, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	defer object.Close()

	return ioutil.ReadAll(object)
}

func (p *minioClient) ListObjects(ctx context.Context, bucket string, opts ...Option) (ListObjectsResult, error) {
	op := ForOption(opts...)
	client, err := p.client(op.SessionToken)
	if err != nil {
		return ListObjectsResult{}, err
	}

	var result ListObjectsResult
	listOptions := minio.ListObjectsOptions{Prefix: op.Prefix}
	for object := range client.ListObjects(ctx, bucket, listOptions) {
		result.Objects = append(result.Objects, ObjectProperties{
			Key:  object.Key,
			Type: "",
			Size: object.Size,
			ETag: object.ETag,
		})
	}

	return result, nil
}

// MakePrivateURL 创建私有URL
func (p *minioClient) MakePrivateURL(ctx context.Context, bucketName string, key string, deadline int64, opts ...Option) (string, error) {
	return "", nil
}

func (p *minioClient) client(token string) (*minio.Client, error) {
	client, err := minio.New(p.opts.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(p.opts.AccessKey, p.opts.AccessSecret, token),
		Secure: false,
	})
	if err != nil {
		return nil, err
	}

	return client, nil
}

// process 进度
type process struct {
	ConsumedBytes int64            // 当前数量
	TotalBytes    int64            // 总数
	Listener      ProgressListener // 监听接口
}

func (p *process) Read(data []byte) (n int, err error) {
	n = len(data)
	p.ConsumedBytes += int64(n)
	p.Listener.ProgressChanged(&ProgressEvent{
		ConsumedBytes: p.ConsumedBytes,
		TotalBytes:    p.TotalBytes,
	})
	return n, nil
}
