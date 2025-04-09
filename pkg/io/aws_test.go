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
		},
	}
	clientIo := io.NewCloudClient(profiles)
	AwsIo = io.NewAwsClient(clientIo)
}
func TestDescribeAwsInstanceTypes(t *testing.T) {
	timeStart := time.Now()
	instanceTypes, err := AwsIo.DescribeInstanceTypes("aws", "cn-northwest-1")
	if err != nil {
		t.Error(err)
	}
	for _, instanceType := range instanceTypes {
		fmt.Println(tea.Prettify(instanceType))
	}
	t.Log("Success.", time.Since(timeStart), len(instanceTypes))
}

func TestDescribeAwsImages(t *testing.T) {
	timeStart := time.Now()
	images, err := AwsIo.DescribeImages("aws", "cn-northwest-1", model.CommonFilter{})
	if err != nil {
		t.Error(err)
	}
	t.Log("Success.", time.Since(timeStart), len(images))
}

func TestQueryAwsKeyPairs(t *testing.T) {
	timeStart := time.Now()
	keyPairs, err := AwsIo.DescribeKeyPairs("aws", "cn-northwest-1", model.CommonFilter{})
	if err != nil {
		t.Error(err)
	}
	for _, keyPair := range keyPairs {
		fmt.Println(tea.Prettify(keyPair))
	}
	t.Log("Success.", time.Since(timeStart), len(keyPairs))
}

func TestQueryAwsSecurityGroups(t *testing.T) {
	timeStart := time.Now()
	securityGroups, err := AwsIo.QuerySecurityGroups("aws", "cn-northwest-1", model.CommonFilter{})
	if err != nil {
		t.Error(err)
	}
	for _, securityGroup := range securityGroups {
		fmt.Println(tea.Prettify(securityGroup))
	}
	t.Log("Success.", time.Since(timeStart), len(securityGroups))
}

func TestQueryEmrCluster(t *testing.T) {
	timeStart := time.Now()
	period := 24 * time.Hour
	filter := model.EmrFilter{
		Profile: tea.String("aws"),
		Region:  tea.String("us-east-1"),
		Period:  &period,
		ClusterStates: []model.EMRClusterStatus{
			model.EMRClusterRunning,
			model.EMRClusterWaiting,
			// model.EMRClusterTerminated,
		},
		// NextMarker: tea.String("xxx"),
	}
	// filter.ClusterStates = []model.EMRClusterStatus{model.EMRClusterRunning}
	resp, err := AwsIo.QueryEmrCluster(filter)
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
	clusters, err := AwsIo.DescribeEmrCluster(model.DescribeInput{
		Profile: tea.String("aws"),
		Region:  tea.String("us-east-1"),
		IDS:     []*string{tea.String("j-xxx")},
	})
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(tea.Prettify(clusters))
	t.Log("Success.", time.Since(timeStart), len(clusters))
}

func TestDescribeAwsVolumes(t *testing.T) {
	timeStart := time.Now()
	volumes, err := AwsIo.DescribeVolumes("aws", "cn-northwest-1", model.DescribeVolumesInput{
		// VolumeIDs: []string{"vol-xxx"},
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("Success.", time.Since(timeStart), len(volumes))
}

// TEST TestAwsDescribeInstances
func TestAwsDescribeInstances(t *testing.T) {
	filter := model.InstanceFilter{
		Size: tea.Int64(6),
	}
	instances, err := AwsIo.DescribeInstances("aws", "cn-northwest-1", filter.ToAwsDescribeInstancesInput())
	if err != nil {
		t.Error(err)
		return
	}
	for _, instance := range instances.Instances {
		fmt.Println(*instance.InstanceID, *instance.CreatedTime, tea.Prettify(instance.InstanceChargeType))
	}
	t.Log("Success.", len(instances.Instances))
}

func TestAwsDescribeInstancesAll(t *testing.T) {
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

// TestAwsCreateSqs
func TestAwsCreateSqs(t *testing.T) {
	err := AwsIo.CreateSqs("aws", "cn-northwest-1", model.CreateSqsRequest{
		QueueName: "zhoushoujian3",
		Type:      "fifo",
		Config: model.SqsConfig{
			VisibilityTimeout:  3600,
			MessageRetention:   86400,
			MaximumMessageSize: 262144,
			ReceiveWaitTime:    20,
			DelaySeconds:       0,
		},
		RedrivePolicy: &model.RedrivePolicy{
			MaxReceiveCount:     "5",
			DeadLetterTargetArn: "arn:aws-cn:sqs:cn-northwest-1:955466075186:zhoushoujian2.fifo",
		},
		Policy:     `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":"sqs:SendMessage","Resource":"*","Condition":{"ArnEquals":{"aws:SourceArn":"arn:aws:s3:*:*:*/*"}}}]}`,
		Encryption: false,
		Tags: map[string]string{
			"Owner": "zhoushoujian",
		},
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("Success.")
}
