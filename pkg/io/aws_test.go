package io_test

import (
	"os"
	"testing"
	"time"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/joho/godotenv"

	"github.com/xops-infra/multi-cloud-sdk/pkg/io"
	"github.com/xops-infra/multi-cloud-sdk/pkg/model"
)

var AwsIo model.CloudIO

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
	}
	clientIo := io.NewCloudClient(profiles)
	AwsIo = io.NewAwsClient(clientIo)
}

func TestAwsDescribeInstances(t *testing.T) {
	{
		timeStart := time.Now()
		filter := model.InstanceFilter{}
		filter.PrivateIp = tea.String(os.Getenv("TEST_AWS_PRIVATE_IP"))
		instances, err := AwsIo.DescribeInstances("aws", "cn-northwest-1", filter.ToAwsDescribeInstancesInput())
		if err != nil {
			t.Error(err)
			return
		}
		t.Log("PrivateIp Success.", time.Since(timeStart), len(instances.Instances))
	}
	{
		timeStart := time.Now()
		filter := model.InstanceFilter{}
		filter.PublicIp = tea.String(os.Getenv("TEST_AWS_PUBLIC_IP"))
		instances, err := AwsIo.DescribeInstances("aws", "cn-northwest-1", filter.ToAwsDescribeInstancesInput())
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
		instances, err := AwsIo.DescribeInstances("aws", "cn-northwest-1", filter.ToAwsDescribeInstancesInput())
		if err != nil {
			t.Error(err)
			return
		}
		t.Log("Owner Success.", time.Since(timeStart), len(instances.Instances))
	}
	{
		timeStart := time.Now()
		filter := model.InstanceFilter{}
		filter.IDs = []*string{tea.String(os.Getenv("TEST_AWS_ID"))}
		instances, err := AwsIo.DescribeInstances("aws", "cn-northwest-1", filter.ToAwsDescribeInstancesInput())
		if err != nil {
			t.Error(err)
			return
		}
		t.Log("ID Success.", time.Since(timeStart), len(instances.Instances))
	}
	{
		timeStart := time.Now()
		filter := model.InstanceFilter{}
		filter.Name = tea.String(os.Getenv("TEST_AWS_NAME"))
		instances, err := AwsIo.DescribeInstances("aws", "cn-northwest-1", filter.ToAwsDescribeInstancesInput())
		if err != nil {
			t.Error(err)
			return
		}
		t.Log("Name Success.", time.Since(timeStart), len(instances.Instances))
	}
	{
		timeStart := time.Now()
		filter := model.InstanceFilter{}
		filter.Status = model.InstanceStatusRunning.TString()
		instances, err := AwsIo.DescribeInstances("aws", "cn-northwest-1", filter.ToAwsDescribeInstancesInput())
		if err != nil {
			t.Error(err)
			return
		}
		t.Log("Status Success.", time.Since(timeStart), len(instances.Instances))
	}
}
