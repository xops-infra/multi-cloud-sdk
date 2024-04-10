package io_test

import (
	"testing"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/xops-infra/multi-cloud-sdk/pkg/model"
)

// CreateBucket
func TestCreateBucket(t *testing.T) {
	t.Log("CreateBucket")
	err := TencentIo.CreateBucket(profile, "ap-beijing", model.CreateBucketRequest{
		BucketName: tea.String("test-bucket-zsj-1-xxx"),
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
