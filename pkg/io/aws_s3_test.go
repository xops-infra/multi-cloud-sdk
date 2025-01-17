package io_test

import (
	"fmt"
	"testing"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/xops-infra/multi-cloud-sdk/pkg/model"
)

func TestGetBucketLifecycle(t *testing.T) {
	resp, err := AwsIo.GetBucketLifecycle("aws", "cn-northwest-1", model.GetBucketLifecycleRequest{
		Bucket: tea.String("ops-9554"),
	})
	if err != nil {
		t.Error(err)
	}
	fmt.Println(tea.Prettify(resp))
}

func TestAwsCreateBucketLifecycle(t *testing.T) {
	err := AwsIo.CreateBucketLifecycle("aws", "us-east-1", model.CreateBucketLifecycleRequest{
		Bucket: tea.String("zhoushoujiantest"),
		Lifecycles: []model.Lifecycle{
			{
				ID: tea.String("删除桶数据"),
				NoncurrentVersionExpiration: &model.LifecycleNoncurrentVersionExpiration{
					NoncurrentDays: tea.Int(1),
				},
				AbortIncompleteMultipartUpload: &model.LifecycleAbortIncompleteMultipartUpload{
					DaysAfterInitiation: tea.Int(3),
				},
				Filter: &model.LifecycleFilter{
					Prefix: tea.String(""),
				},
			},
		},
	})
	if err != nil {
		t.Error(err)
	}
}

func TestCreateS3Bucket(t *testing.T) {
	err := AwsIo.CreateBucket("aws", "us-east-1", model.CreateBucketRequest{
		BucketName: tea.String("test-bucket-zsj-1"),
		Tags: model.Tags{
			{Key: "Owner", Value: "zhoushoujian"},
			{Key: "Team", Value: "ops"},
			{Key: "Env", Value: "abcszsj"},
		},
	})
	if err != nil {
		t.Error(err)
	}
}

// TestListS3Bucket
func TestListS3Bucket(t *testing.T) {
	resp, err := AwsIo.ListBucket("aws", "us-east-1", model.ListBucketRequest{
		KeyWord: tea.String("test-bucket-zsj"),
	})
	if err != nil {
		t.Error(err)
	}
	t.Log(tea.Prettify(resp))
}

// TestGetObjectPregisn
func TestS3GetObjectPregisn(t *testing.T) {
	resp, err := AwsIo.GetObjectPregisn("aws", "us-east-1", model.ObjectPregisnRequest{
		Bucket: tea.String("zhoushoujiantest"),
		Key:    tea.String("xxx.pdf"),
		Expire: tea.Int64(3600),
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(resp.Url)
}
