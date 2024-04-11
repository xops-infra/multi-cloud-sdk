package io

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/xops-infra/multi-cloud-sdk/pkg/model"
)

// 比官方多了个创建bucket的tags功能。
func (c *awsClient) CreateBucket(profile, region string, input model.CreateBucketRequest) error {
	client, err := c.io.GetAWSS3Client(profile, region)
	if err != nil {
		return err
	}

	_, err = client.CreateBucket(&s3.CreateBucketInput{
		Bucket: input.BucketName,
	})
	if err != nil {
		return err
	}

	// add tags
	if len(input.Tags) > 0 {
		_, err = client.PutBucketTagging(&s3.PutBucketTaggingInput{
			Bucket:  input.BucketName,
			Tagging: &s3.Tagging{TagSet: input.Tags.ToAWSS3Tags()},
		})
		if err != nil {
			return fmt.Errorf("create bucket success. put bucket tags failed: %v", err)
		}
	}

	return nil
}

// 删除走后台人工吧，接口不支持。
func (c *awsClient) DeleteBucket(profile, region string, input model.DeleteBucketRequest) (model.DeleteBucketResponse, error) {
	panic("implement me")
	// client, err := c.io.GetAWSS3Client(profile, region)
	// if err != nil {
	// 	return model.DeleteBucketResponse{}, err

	// }
	// resp, err := client.DeleteBucket(&s3.DeleteBucketInput{
	// 	Bucket: input.BucketName,
	// })
	// if err != nil {
	// 	return model.DeleteBucketResponse{}, err

	// }
	// return model.DeleteBucketResponse{
	// 	Meta: resp.GoString(),
	// }, nil
}

// 比官方多了个查询桶标签和地域的功能。
func (c *awsClient) ListBucket(profile, region string, input model.ListBucketRequest) (model.ListBucketResponse, error) {
	client, err := c.io.GetAWSS3Client(profile, region)
	if err != nil {
		return model.ListBucketResponse{}, err
	}
	resp, err := client.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		return model.ListBucketResponse{}, err
	}
	var buckets []*model.Bucket
	wg := sync.WaitGroup{}
	for _, bucket := range resp.Buckets {
		if input.KeyWord != nil && *input.KeyWord != "" {
			if !strings.Contains(*bucket.Name, *input.KeyWord) {
				continue
			}
		}
		var createTime string
		if bucket.CreationDate != nil {
			createTime = bucket.CreationDate.Local().Format(time.DateTime)
		}

		newBucket := &model.Bucket{
			Name:       tea.StringValue(bucket.Name),
			CreateTime: createTime,
		}
		wg.Add(1)
		go func(bucket *model.Bucket) {
			// 查询一下桶地域
			defer wg.Done()
			locationResp, err := client.GetBucketLocation(&s3.GetBucketLocationInput{
				Bucket: &bucket.Name,
			})
			if err != nil {
				if strings.Contains(err.Error(), "AccessDenied") {
					locationResp = &s3.GetBucketLocationOutput{
						LocationConstraint: tea.String("AccessDenied"),
					}
				}
				return
			}
			if locationResp.LocationConstraint == nil {
				bucket.Location = region
			} else {
				bucket.Location = tea.StringValue(locationResp.LocationConstraint)
			}

			// 查询一下桶标签
			tagResp, err := client.GetBucketTagging(&s3.GetBucketTaggingInput{
				Bucket: &bucket.Name,
			})
			if err != nil {
				if strings.Contains(err.Error(), "NoSuchTagSet") {
					tagResp = &s3.GetBucketTaggingOutput{}
				} else if strings.Contains(err.Error(), "AccessDenied") {
					tagResp = &s3.GetBucketTaggingOutput{
						TagSet: []*s3.Tag{
							{Key: tea.String("Error"), Value: tea.String("AccessDenied")},
						},
					}
				} else {
					return
				}
			}
			bucket.Tags = model.NewTagsFromAWSS3Tags(tagResp.TagSet)
		}(newBucket)

		buckets = append(buckets, newBucket)
	}
	wg.Wait()
	return model.ListBucketResponse{
		Buckets: buckets,
		Total:   int64(len(buckets)),
	}, nil
}

func (c *awsClient) GetObjectPregisn(profile, region string, req model.ObjectPregisnRequest) (model.ObjectPregisnResponse, error) {
	client, err := c.io.GetAWSS3Client(profile, region)
	if err != nil {
		return model.ObjectPregisnResponse{}, err
	}
	return c.getObjectPregisn(client, req)
}

func (c *awsClient) GetObjectPregisnWithAKSK(ak, sk, region string, req model.ObjectPregisnRequest) (model.ObjectPregisnResponse, error) {
	cre := credentials.NewStaticCredentials(ak, sk, "")
	session, err := session.NewSession(aws.NewConfig().WithCredentials(cre))
	if err != nil {
		return model.ObjectPregisnResponse{}, err
	}
	session.Config.Region = aws.String(region)
	client := s3.New(session)
	return c.getObjectPregisn(client, req)
}

func (c *awsClient) getObjectPregisn(client *s3.S3, req model.ObjectPregisnRequest) (model.ObjectPregisnResponse, error) {
	// head object
	_, err := client.HeadObject(&s3.HeadObjectInput{
		Bucket: req.Bucket,
		Key:    req.Key,
	})
	if err != nil {
		return model.ObjectPregisnResponse{}, err
	}

	// request object
	resp, _ := client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: req.Bucket,
		Key:    req.Key,
	})
	timeD := 3600 * time.Second // 1 hour
	if req.Expire != nil {
		timeD = time.Duration(*req.Expire) * time.Second
	}
	urlPresign, err := resp.Presign(timeD)
	if err != nil {
		return model.ObjectPregisnResponse{}, err
	}
	return model.ObjectPregisnResponse{
		Url: urlPresign,
	}, nil
}
