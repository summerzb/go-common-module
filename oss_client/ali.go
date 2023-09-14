package oss_client

import (
	"context"
	"fmt"
	"io"
	"path"
	"strconv"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/sts"
	aliOss "github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/jinzhu/now"
)

type aliClient struct {
	opts OssConfig
}

func NewAliClient(opts OssConfig) OssClientIface {
	return &aliClient{
		opts: opts,
	}
}

// Token 获取token
func (p *aliClient) Token(ctx context.Context, bucket string, objectKey string, opts ...Option) (*Token, error) {
	option := ForOption(opts...)

	//设置调用者（RAM用户或RAM角色）的AccessKey ID和AccessKey Secret。
	client, _ := sts.NewClientWithAccessKey(p.opts.RegionId, p.opts.AccessKey, p.opts.AccessSecret)

	//构建请求对象。
	request := sts.CreateAssumeRoleRequest()
	request.Scheme = "https"

	//设置参数。
	request.RoleArn = p.opts.RoleArn
	request.RoleSessionName = p.opts.RoleSessionName
	request.DurationSeconds = requests.NewInteger(option.TokenExpiry)
	request.Policy = fmt.Sprintf(`{
		"Version": "1",
		"Statement": [
			{
				"Action": [
					"oss:PutObject"
				],
				"Effect": "Allow",
				"Resource": [
					"acs:oss:*:*:%v/%v/*"
				]
			}
		]
	}`, bucket, path.Dir(objectKey))

	//发起请求，并得到响应。
	response, err := client.AssumeRole(request)
	if err != nil {
		return nil, err
	}

	expiration, err := now.Parse(response.Credentials.Expiration)
	if err != nil {
		return nil, err
	}

	token := &Token{
		Token:           response.Credentials.SecurityToken,
		AccessKey:       response.Credentials.AccessKeyId,
		AccessKeySecret: response.Credentials.AccessKeySecret,
		Expiration:      expiration.Unix(),
		Bucket:          bucket,
		Endpoint:        p.opts.DownloadUrl,
	}

	return token, nil
}

// ObjectMeta 获取对象基本信息
func (p *aliClient) ObjectMeta(ctx context.Context, bucketName string, objectKey string, opts ...Option) (*ObjectMeta, error) {
	bucket, err := p.getBucket(ctx, bucketName)
	if err != nil {
		return nil, err
	}

	props, err := bucket.GetObjectDetailedMeta(objectKey)
	if err != nil {
		return nil, err
	}

	metaInfo := &ObjectMeta{}
	metaInfo.Hash = props.Get(aliOss.HTTPHeaderContentMD5)
	metaInfo.FileSize, err = strconv.ParseInt(props.Get(aliOss.HTTPHeaderContentLength), 10, 64)
	if err != nil {
		return nil, err
	}

	return metaInfo, nil
}

// PutObject 上传对象
func (p *aliClient) PutObject(ctx context.Context, bucketName string, objectKey string, data io.Reader, size int64, opts ...Option) (string, error) {
	bucket, err := p.getBucket(ctx, bucketName)
	if err != nil {
		return "", err
	}

	if err = bucket.PutObject(objectKey, data); err != nil {
		return "", err
	}

	return p.opts.DownloadUrl + "/" + objectKey, nil
}

func (p *aliClient) PutObjectFromFile(ctx context.Context, bucketName string, objectKey string, filePath string, opts ...Option) (string, error) {
	bucket, err := p.getBucket(ctx, bucketName)
	if err != nil {
		return "", err
	}

	if err = bucket.PutObjectFromFile(objectKey, filePath); err != nil {
		return "", err
	}

	return p.opts.DownloadUrl + "/" + objectKey, nil
}

// GetObject 获取对象
func (p *aliClient) GetObject(ctx context.Context, bucketName string, objectKey string, opts ...Option) ([]byte, error) {
	bucket, err := p.getBucket(ctx, bucketName)
	if err != nil {
		return nil, err
	}

	body, err := bucket.GetObject(objectKey)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	return io.ReadAll(body)
}

// BatchCopy 批量复制
func (p *aliClient) BatchCopy(ctx context.Context, bucketName string, copyKeys map[string]string, opts ...Option) error {
	bucket, err := p.getBucket(ctx, bucketName)
	if err != nil {
		return err
	}
	// todo: 目前sdk不支持批量复制，采用循环复制，后面可以改为多线程复制
	for src, dest := range copyKeys {
		_, err = bucket.CopyObject(src, dest)
		if err != nil {
			return err
		}
	}

	return nil
}

// Copy 复制对象
func (p *aliClient) Copy(ctx context.Context, bucketName string, srcObject string, dstObject string, opts ...Option) error {
	bucket, err := p.getBucket(ctx, bucketName)
	if err != nil {
		return err
	}

	_, err = bucket.CopyObject(srcObject, dstObject)
	return err
}

func (p *aliClient) ListObjects(ctx context.Context, bucketName string, opts ...Option) (ListObjectsResult, error) {
	bucket, err := p.getBucket(ctx, bucketName)
	if err != nil {
		return ListObjectsResult{}, err
	}

	op := ForOption(opts...)
	listResp, err := bucket.ListObjects(aliOss.Prefix(op.Prefix), aliOss.Marker(op.Marker), aliOss.Delimiter(op.Delimiter), aliOss.MaxKeys(op.Limit))
	if err != nil {
		return ListObjectsResult{}, err
	}

	out := ListObjectsResult{
		Prefix:         op.Prefix,
		Marker:         op.Marker,
		Delimiter:      op.Delimiter,
		IsTruncated:    listResp.IsTruncated,
		NextMarker:     listResp.NextMarker,
		Objects:        nil,
		CommonPrefixes: listResp.CommonPrefixes,
	}
	for _, v := range listResp.Objects {
		out.Objects = append(out.Objects, ObjectProperties{
			Key:  v.Key,
			Type: v.Type,
			Size: v.Size,
			ETag: v.ETag,
		})
	}

	return out, nil
}

// MakePrivateURL 创建私有URL
func (p *aliClient) MakePrivateURL(ctx context.Context, bucketName string, key string, deadline int64, opts ...Option) (string, error) {
	return "", nil
}

func (p *aliClient) getBucket(ctx context.Context, bucketName string) (*aliOss.Bucket, error) {
	// NewAliClient client
	client, err := aliOss.New(p.opts.Endpoint, p.opts.AccessKey, p.opts.AccessSecret)
	if err != nil {
		return nil, err
	}

	// Get bucket
	bucket, err := client.Bucket(bucketName)

	return bucket, err
}
