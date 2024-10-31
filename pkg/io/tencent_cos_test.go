package io_test

import (
	"testing"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/xops-infra/multi-cloud-sdk/pkg/model"
)

// TEST CreateBucketLifecycle
func TestCreateBucketLifecycle(t *testing.T) {
	t.Log("CreateBucketLifecycle")
	err := TencentIo.CreateBucketLifecycle(profile, "na-ashburn", model.CreateBucketLifecycleRequest{
		Bucket: tea.String("examplebucket-1250000000"),
		Lifecycles: []model.Lifecycle{
			{
				ID:     tea.String("huggingface模型定期删除"),
				Filter: &model.LifecycleFilter{Prefix: tea.String("hg/")},
				Status: tea.Bool(true),
				Expiration: &model.LifecycleExpiration{
					Days: tea.Int(6),
				},
			},
			{
				ID:     tea.String("OPS_BASE"),
				Filter: &model.LifecycleFilter{},
				Status: tea.Bool(true),
				AbortIncompleteMultipartUpload: &model.LifecycleAbortIncompleteMultipartUpload{
					DaysAfterInitiation: tea.Int(40),
				},
			},
		},
	})
	if err != nil {
		t.Error(err)
	}
}

// CreateBucket
func TestCreateBucket(t *testing.T) {
	t.Log("CreateBucket")
	err := TencentIo.CreateBucket(profile, "ap-beijing", model.CreateBucketRequest{
		BucketName: tea.String("examplebucket-1250000000"),
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

// DeleteBucket

// ListBucket
func TestListBucket(t *testing.T) {
	t.Log("ListBucket")
	resp, err := TencentIo.ListBucket(profile, "", model.ListBucketRequest{
		KeyWord: tea.String(""),
	})
	if err != nil {
		t.Error(err)
	}
	t.Log(tea.Prettify(resp))
}

// GetObjectPregisn
func TestGetObjectPregisn(t *testing.T) {
	t.Log("GetObjectPregisn")
	resp, err := TencentIo.GetObjectPregisn(profile, "ap-shanghai", model.ObjectPregisnRequest{
		Bucket: tea.String("examplebucket-1250000000"),
		Key:    tea.String("test.txt"),
		Expire: tea.Int64(3600),
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(resp.Url)
}
