package io_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/xops-infra/multi-cloud-sdk/pkg/model"
)

var lifecycleFile = "../../lifecycle.json"

func TestCreateCOSLifecycle(t *testing.T) {
	// 解析 jason 文件
	data, err := os.ReadFile(lifecycleFile)
	if err != nil {
		panic(err)
	}
	var lifecycleNewJson []model.Lifecycle
	err = json.Unmarshal(data, &lifecycleNewJson)
	if err != nil {
		panic(err)
	}
	req := model.CreateBucketLifecycleRequest{
		Bucket:     tea.String("zhoushoujiantest-1251949819"),
		Lifecycles: lifecycleNewJson,
	}

	err = TencentIo.CreateBucketLifecycle(profile, "na-ashburn", req)
	if err != nil {
		t.Error(err)
	}

}

// TEST GetBucketLifecycle
func TestGetTencentBucketLifecycle(t *testing.T) {
	resp, err := TencentIo.GetBucketLifecycle(profile, "ap-shanghai", model.GetBucketLifecycleRequest{})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(tea.Prettify(resp))
}

// TEST CreateBucketLifecycle
func TestCreateBucketLifecycle(t *testing.T) {
	err := TencentIo.CreateBucketLifecycle(profile, "ap-shanghai", model.CreateBucketLifecycleRequest{
		Bucket: tea.String("examplebucket-1250000000"),
		Lifecycles: []model.Lifecycle{
			{
				ID: tea.String("OPS_BASE"),
				Filter: &model.LifecycleFilter{
					Prefix: tea.String(""), // "" 整个桶；
				},
				AbortIncompleteMultipartUpload: &model.LifecycleAbortIncompleteMultipartUpload{
					DaysAfterInitiation: tea.Int(20),
				}, // 20 天删除碎片
				// NoncurrentVersionExpiration: &model.LifecycleNoncurrentVersionExpiration{
				// 	NoncurrentDays: tea.Int(360),
				// }, // 当前版本文件删除：360 天
				// NoncurrentVersionTransitions: []model.LifecycleNoncurrentVersionTransition{
				// 	{
				// 		NoncurrentDays: tea.Int(30),
				// 		StorageClass:   tea.String("STANDARD_IA"),
				// 	},
				// }, // 历史版本沉降至低频存储：30 天
				// Transitions: []model.LifecycleTransition{
				// 	{
				// 		Days:         tea.Int(40),
				// 		StorageClass: tea.String("STANDARD_IA"),
				// 	},
				// }, // 当前版本文件沉降至低频存储：40 天
				// Expiration: &model.LifecycleExpiration{
				// 	Days: tea.Int(180),
				// }, // 当前版本文件删除：180 天
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
