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
				"us-east-1",
			},
		}, {
			Name:  "aws-us",
			Cloud: model.AWS,
			AK:    os.Getenv("AWS_US_ACCESS_KEY_ID"),
			SK:    os.Getenv("AWS_US_SECRET_ACCESS_KEY"),
			Regions: []string{
				"us-east-1",
			},
		},
	}
	clientIo := io.NewCloudClient(profiles)
	AwsIo = io.NewAwsClient(clientIo)
}

func TestQueryEmrCluster(t *testing.T) {
	timeStart := time.Now()
	period := 24 * time.Hour
	filter := model.EmrFilter{
		Period: &period,
		ClusterStates: []model.EMRClusterStatus{
			model.EMRClusterRunning,
			model.EMRClusterWaiting,
			// model.EMRClusterTerminated,
		},
		// NextMarker: tea.String("xxx"),
	}
	// filter.ClusterStates = []model.EMRClusterStatus{model.EMRClusterRunning}
	resp, err := AwsIo.QueryEmrCluster("aws-us", "us-east-1", filter)
	if err != nil {
		t.Error(err)
		return
	}
	for _, cluster := range resp.Clusters {
		fmt.Println(tea.Prettify(cluster))
	}
	if resp.NextMarker != nil {
		t.Log("NextMarker:", *resp.NextMarker)
	}
	t.Log("Success.", time.Since(timeStart), len(resp.Clusters))
}

func TestDescribeEmrCluster(t *testing.T) {
	timeStart := time.Now()
	ids := []*string{tea.String("j-xxxx")}
	clusters, err := AwsIo.DescribeEmrCluster("aws-us", "us-east-1", ids)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(tea.Prettify(clusters))
	t.Log("Success.", time.Since(timeStart), len(clusters))
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
