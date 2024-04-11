package io_test

import (
	"testing"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/xops-infra/multi-cloud-sdk/pkg/model"
)

func TestCreateS3Bucket(t *testing.T) {
	err := AwsIo.CreateBucket("aws-us", "us-east-1", model.CreateBucketRequest{
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
	resp, err := AwsIo.ListBucket("aws-us", "us-east-1", model.ListBucketRequest{
		KeyWord: tea.String("test-bucket-zsj"),
	})
	if err != nil {
		t.Error(err)
	}
	t.Log(tea.Prettify(resp))
}

// TestGetObjectPregisn
func TestS3GetObjectPregisn(t *testing.T) {
	resp, err := AwsIo.GetObjectPregisn("aws-us", "us-east-1", model.ObjectPregisnRequest{
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
