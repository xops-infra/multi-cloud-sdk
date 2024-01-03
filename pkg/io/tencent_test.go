package io_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/joho/godotenv"

	"github.com/xops-infra/multi-cloud-sdk/pkg/io"
	"github.com/xops-infra/multi-cloud-sdk/pkg/model"
)

var TencentIo model.CloudIO

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	profiles := []model.ProfileConfig{
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
	clientIo := io.NewCloudClient(profiles)
	TencentIo = io.NewTencentClient(clientIo)
}

func TestQueryTencentEmrCluster(t *testing.T) {
	timeStart := time.Now()
	filter := model.EmrFilter{}
	instances, err := TencentIo.QueryEmrCluster("tencent", "na-ashburn", filter)
	if err != nil {
		t.Error(err)
		return
	}
	for _, instance := range instances.Clusters {
		fmt.Println(tea.Prettify(instance))
	}
	t.Log("Success.", time.Since(timeStart), len(instances.Clusters))
}

func TestDescribeTencentEmrCluster(t *testing.T) {
	timeStart := time.Now()
	instances, err := TencentIo.DescribeEmrCluster("tencent", "na-ashburn",
		[]*string{tea.String("emr-kthwjob1"), tea.String("emr-bxss1pm3")})
	if err != nil {
		t.Error(err)
		return
	}
	for _, instance := range instances {
		fmt.Println(tea.Prettify(instance))
	}
	t.Log("Success.", time.Since(timeStart), len(instances))
}

func TestDescribeInstances(t *testing.T) {
	{
		timeStart := time.Now()
		filter := model.InstanceFilter{}
		filter.PrivateIp = tea.String(os.Getenv("TEST_TENCENT_PRIVATE_IP"))
		instances, err := TencentIo.DescribeInstances("tencent", "ap-shanghai", filter.ToTxDescribeInstancesInput())
		if err != nil {
			t.Error(err)
			return
		}
		t.Log("PrivateIp Success.", time.Since(timeStart), len(instances.Instances))
	}
	{
		timeStart := time.Now()
		filter := model.InstanceFilter{}
		filter.PublicIp = tea.String(os.Getenv("TEST_TENCENT_PUBLIC_IP"))
		instances, err := TencentIo.DescribeInstances("tencent", "ap-shanghai", filter.ToTxDescribeInstancesInput())
		if err != nil {
			t.Error(err)
			return
		}
		t.Log("PublicIp Success.", time.Since(timeStart), len(instances.Instances))
	}
	{
		timeStart := time.Now()
		filter := model.InstanceFilter{}
		filter.Owner = tea.String("zhoushoujian")
		instances, err := TencentIo.DescribeInstances("tencent", "ap-shanghai", filter.ToTxDescribeInstancesInput())
		if err != nil {
			t.Error(err)
			return
		}
		t.Log("Owner Success.", time.Since(timeStart), len(instances.Instances))
	}
	{
		timeStart := time.Now()
		filter := model.InstanceFilter{}
		filter.IDs = []*string{tea.String(os.Getenv("TEST_TENCENT_ID"))}
		instances, err := TencentIo.DescribeInstances("tencent", "ap-shanghai", filter.ToTxDescribeInstancesInput())
		if err != nil {
			t.Error(err)
			return
		}
		t.Log("ID Success.", time.Since(timeStart), len(instances.Instances))
	}
	{
		timeStart := time.Now()
		filter := model.InstanceFilter{}
		filter.Name = tea.String(os.Getenv("TEST_TENCENT_NAME"))
		instances, err := TencentIo.DescribeInstances("tencent", "ap-shanghai", filter.ToTxDescribeInstancesInput())
		if err != nil {
			t.Error(err)
			return
		}
		t.Log("Name Success.", time.Since(timeStart), len(instances.Instances))
	}
	{
		timeStart := time.Now()
		filter := model.InstanceFilter{}
		filter.Status = model.InstanceStatusStopped.TString()
		instances, err := TencentIo.DescribeInstances("tencent", "ap-shanghai", filter.ToTxDescribeInstancesInput())
		if err != nil {
			t.Error(err)
			return
		}
		t.Log("Status Success.", time.Since(timeStart), len(instances.Instances))
	}
}
