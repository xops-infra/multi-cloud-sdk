package model_test

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/xops-infra/multi-cloud-sdk/pkg/model"
)

var lifecycleFile = "../../lifecycle.json"

// TEST ToCOSLifecycle
func TestToCOSLifecycle(t *testing.T) {
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
		Bucket:     tea.String("examplebucket-1250000000"),
		Lifecycles: lifecycleNewJson,
	}

	cosLc, err := req.ToCOSLifecycle()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(tea.Prettify(cosLc))
}

// TEST ToS3Lifecycle
func TestToS3Lifecycle(t *testing.T) {
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
		Bucket:     tea.String("examplebucket-1250000000"),
		Lifecycles: lifecycleNewJson,
	}

	s3Lc, err := req.ToS3Lifecycle()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(tea.Prettify(s3Lc))
}
