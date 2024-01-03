package main

import (
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"

	"github.com/xops-infra/multi-cloud-sdk/pkg/io"
	"github.com/xops-infra/multi-cloud-sdk/pkg/model"
	"github.com/xops-infra/multi-cloud-sdk/pkg/service"
)

var emrC model.EmrContact

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	profiles := []model.ProfileConfig{
		{
			Name:  "aws",
			Cloud: model.AWS,
			AK:    os.Getenv("AWS_ACCESS_KEY_ID"),
			SK:    os.Getenv("AWS_SECRET_ACCESS_KEY"),
			Regions: []string{
				"cn-northwest-1",
			},
		},
		{
			Name:  "tencent",
			Cloud: model.TENCENT,
			AK:    os.Getenv("TENCENT_ACCESS_KEY"),
			SK:    os.Getenv("TENCENT_SECRET_KEY"),
			Regions: []string{
				"ap-shanghai",
				"na-ashburn",
			},
		},
	}
	cloudIo := io.NewCloudClient(profiles)
	aws := io.NewAwsClient(cloudIo)
	tencent := io.NewTencentClient(cloudIo)
	emrC = service.NewEmrService(profiles, aws, tencent)
}

func TestDescribeEmr(t *testing.T) {
	{
		startTime := time.Now()
		resp, err := emrC.DescribeEmrCluster("aws", "cn-northwest-1", []*string{})
		if err != nil {
			t.Error(err)
		}
		t.Log("aws", time.Since(startTime), len(resp))
	}

	{
		startTime := time.Now()
		resp, err := emrC.DescribeEmrCluster("tencent", "ap-shanghai", []*string{})
		if err != nil {
			t.Error(err)
		}
		t.Log("tencent", time.Since(startTime), len(resp))
	}
	{
		startTime := time.Now()
		resp, err := emrC.DescribeEmrCluster("tencent", "na-ashburn", []*string{})
		if err != nil {
			t.Error(err)
		}
		t.Log("tencent", time.Since(startTime), len(resp))
	}
}

func TestQueryEmr(t *testing.T) {
	{
		startTime := time.Now()
		resp, err := emrC.QueryEmrCluster("aws", "cn-northwest-1", model.EmrFilter{})
		if err != nil {
			t.Error(err)
		}
		t.Log("aws", time.Since(startTime), len(resp.Clusters))
	}

	{
		startTime := time.Now()
		resp, err := emrC.QueryEmrCluster("tencent", "ap-shanghai", model.EmrFilter{})
		if err != nil {
			t.Error(err)
		}
		t.Log("tencent", time.Since(startTime), len(resp.Clusters))
	}
	{
		startTime := time.Now()
		resp, err := emrC.QueryEmrCluster("tencent", "na-ashburn", model.EmrFilter{})
		if err != nil {
			t.Error(err)
		}
		t.Log("tencent", time.Since(startTime), len(resp.Clusters))
	}
}
