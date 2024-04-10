package io

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/xops-infra/multi-cloud-sdk/pkg/model"
)

// Host: <BucketName-APPID>.cos.<Region>.myqcloud.com，其中 <BucketName-APPID> 为带 APPID 后缀的存储桶名字，例如 examplebucket-1250000000
// appid不考虑在配置中获取，考虑云配置的一致性，特性的东西不应该放在配置中。
func (c *tencentClient) CreateBucket(profile, region string, input model.CreateBucketRequest) error {
	if input.BucketName == nil || region == "" {
		return fmt.Errorf("bucket name or region is empty")
	}
	// 判断新建桶是否带了appid，类似 examplebucket-1250000000，严格判断appid类型
	// 获取最后一段字符串，判断是否为数字，长度是否为10
	bucketName := *input.BucketName
	strArry := strings.Split(bucketName, "-")
	appid := strArry[len(strArry)-1]
	if len(appid) != 10 {
		return fmt.Errorf("appid is invalid")
	}
	// 组成的完整请求域名字符数总计最多60个字符
	if len(*input.BucketName) > 60 {
		return fmt.Errorf("bucket name is too long")
	}
	host := fmt.Sprintf("https://%s.cos.%s.myqcloud.com", *input.BucketName, region)
	client, err := c.io.GetTencentCosClient(profile, host)
	if err != nil {
		return err
	}
	_, err = client.Bucket.Put(context.Background(), &cos.BucketPutOptions{
		XCosACL: "private",
	})
	if err != nil {
		return err
	}

	// add tags
	if len(input.Tags) > 0 {
		_, err = client.Bucket.PutTagging(context.Background(), &cos.BucketPutTaggingOptions{
			TagSet: input.Tags.ToTencentCosTags(),
		})
		if err != nil {
			return fmt.Errorf("create bucket success. put bucket tags failed: %v", err)
		}
	}

	return nil
}

// 删除走后台人工吧，接口不支持。
func (c *tencentClient) DeleteBucket(profile, region string, input model.DeleteBucketRequest) (model.DeleteBucketResponse, error) {
	panic("implement me")
}

// Host: 查询全部存储桶列表指定为 service.cos.myqcloud.com，查询特定地域下的存储桶列表指定为 cos.<Region>.myqcloud.com，其中 <Region> 为 COS 的可用地域
func (c *tencentClient) ListBucket(profile, region string, input model.ListBucketRequest) (model.ListBucketResponse, error) {
	host := "http://service.cos.myqcloud.com"
	if region != "" {
		host = fmt.Sprintf("http://cos.%s.myqcloud.com", region)
	}
	client, err := c.io.GetTencentCosClient(profile, host)
	if err != nil {
		return model.ListBucketResponse{}, err
	}
	opt := &cos.ServiceGetOptions{
		// Region:  region,
		MaxKeys: 20,
	}

	result, _, err := client.Service.Get(context.Background(), opt)
	if err != nil {
		return model.ListBucketResponse{}, err
	}

	var buckets []*model.Bucket
	for {
		wg := sync.WaitGroup{}
		for _, bucket := range result.Buckets {
			fmt.Println(tea.Prettify(bucket))
			if input.KeyWord != nil && *input.KeyWord != "" {
				if !strings.Contains(bucket.Name, *input.KeyWord) {
					continue
				}
			}
			createTime, _ := time.Parse(time.RFC3339, bucket.CreationDate)
			newBucket := &model.Bucket{
				Name:       bucket.Name,
				CreateTime: createTime.Local().Format(time.DateTime),
				Location:   bucket.Region,
			}
			wg.Add(1)
			go func(bucket *model.Bucket) {
				defer wg.Done()
				// get tags
				client.BaseURL.BucketURL = &url.URL{
					Scheme: "https",
					Host:   fmt.Sprintf("%s.cos.%s.myqcloud.com", bucket.Name, bucket.Location),
				}
				result, _, err := client.Bucket.GetTagging(context.Background())
				if err != nil {
					return
				}
				bucket.Tags = model.NewTagsFromTencentCosTags(result.TagSet)

			}(newBucket)

			buckets = append(buckets, newBucket)
		}
		wg.Wait()
		if result.IsTruncated {
			opt.Marker = result.NextMarker
			result, _, err = client.Service.Get(context.Background(), opt)
			if err != nil {
				return model.ListBucketResponse{}, err
			}
		} else {
			break
		}
	}

	return model.ListBucketResponse{
		Buckets: buckets,
		Total:   int64(len(buckets)),
	}, nil
}
