package main

import (
	"os"
	"testing"
	"time"

	"github.com/alibabacloud-go/tea/tea"
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
		resp, err := emrC.DescribeEmrCluster(model.DescribeInput{
			Profile: tea.String("aws"),
			Region:  tea.String("cn-northwest-1"),
		})
		if err != nil {
			t.Error(err)
		}
		t.Log("aws", time.Since(startTime), len(resp))
	}

	{
		startTime := time.Now()
		resp, err := emrC.DescribeEmrCluster(model.DescribeInput{
			Profile: tea.String("tencent"),
			Region:  tea.String("ap-shanghai"),
		})
		if err != nil {
			t.Error(err)
		}
		t.Log("tencent", time.Since(startTime), len(resp))
	}
	{
		startTime := time.Now()
		resp, err := emrC.DescribeEmrCluster(model.DescribeInput{
			Profile: tea.String("tencent"),
			Region:  tea.String("na-ashburn"),
		})
		if err != nil {
			t.Error(err)
		}
		t.Log("tencent", time.Since(startTime), len(resp))
	}
}

func TestQueryEmr(t *testing.T) {
	{
		startTime := time.Now()
		resp, err := emrC.QueryEmrCluster(model.EmrFilter{
			Profile: tea.String("aws"),
		})
		if err != nil {
			t.Error(err)
		}
		t.Log("aws", time.Since(startTime), len(resp.Clusters))
	}

	{
		startTime := time.Now()
		resp, err := emrC.QueryEmrCluster(model.EmrFilter{
			Profile: tea.String("tencent"),
		})
		if err != nil {
			t.Error(err)
		}
		t.Log("tencent", time.Since(startTime), len(resp.Clusters))
	}
	{
		startTime := time.Now()
		resp, err := emrC.QueryEmrCluster(model.EmrFilter{
			Profile: tea.String("tencent"),
			Region:  tea.String("na-ashburn"),
		})
		if err != nil {
			t.Error(err)
		}
		t.Log("tencent", time.Since(startTime), len(resp.Clusters))
	}
}
