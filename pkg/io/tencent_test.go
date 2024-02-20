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
				"ap-beijing",
				"na-ashburn",
			},
		},
	}
	clientIo := io.NewCloudClient(profiles)
	TencentIo = io.NewTencentClient(clientIo)
}

func TestQueryTencentEmrCluster(t *testing.T) {
	timeStart := time.Now()
	filter := model.EmrFilter{
		Profile: tea.String("tencent"),
		Region:  tea.String("na-ashburn"),
	}
	instances, err := TencentIo.QueryEmrCluster(filter)
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
	instances, err := TencentIo.DescribeEmrCluster(model.DescribeInput{
		Profile: tea.String("tencent"),
		Region:  tea.String("na-ashburn"),
		IDS:     []*string{tea.String("emr-xxx")},
	})
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

func TestCreateTags(t *testing.T) {
	input := model.CreateTagsInput{
		Tags: model.Tags{
			{
				Key:   "CreateTime",
				Value: time.Now().Format("2006010215"),
			},
		},
	}
	err := TencentIo.CreateTags("tencent", "ap-shanghai", input)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("Success. %s", tea.Prettify(input))
}

func TestCreateInstance(t *testing.T) {
	resp, err := TencentIo.CreateInstance("tencent", "ap-shanghai", model.CreateInstanceInput{
		Name:             tea.String("multi-cloud-sdk-test"),
		ImageID:          tea.String("img-hdt9xxkt"),
		InstanceType:     tea.String("SA5.MEDIUM2"),
		KeyIds:           []*string{tea.String(os.Getenv("TEST_TENCENT_KEY_ID"))},
		Zone:             tea.String("ap-shanghai-5"),
		VpcID:            tea.String(os.Getenv("TEST_TENCENT_VPC_ID")),
		SubnetID:         tea.String(os.Getenv("TEST_TENCENT_SUBNET_ID")),
		SecurityGroupIDs: []*string{tea.String(os.Getenv("TEST_TENCENT_SECURITY_GROUP_ID"))},
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("Success. %s", tea.Prettify(resp))
}

func TestModifyInstance(t *testing.T) {
	instancesIds := []*string{tea.String("ins-iwh5ysbx")}

	// // StartInstance
	// {
	// 	resp, err := TencentIo.ModifyInstance("tencent", "ap-shanghai", model.ModifyInstanceInput{
	// 		Action:      model.StartInstance,
	// 		InstanceIDs: instancesIds,
	// 	})
	// 	if err != nil {
	// 		t.Error(err)
	// 		return
	// 	}
	// 	t.Logf("StartInstance Success. %s", tea.Prettify(resp))
	// 	time.Sleep(30 * time.Second)
	// }
	// // RebootInstance
	// {
	// 	resp, err := TencentIo.ModifyInstance("tencent", "ap-shanghai", model.ModifyInstanceInput{
	// 		Action:      model.RebootInstance,
	// 		InstanceIDs: instancesIds,
	// 	})
	// 	if err != nil {
	// 		t.Error(err)
	// 		return
	// 	}
	// 	t.Logf("RebootInstance Success. %s", tea.Prettify(resp))
	// 	time.Sleep(30 * time.Second)
	// }

	// ResetInstance
	{
		resp, err := TencentIo.ModifyInstance("tencent", "ap-shanghai", model.ModifyInstanceInput{
			Action:      model.ResetInstance,
			InstanceIDs: instancesIds,
		})
		if err != nil {
			t.Error(err)
			return
		}
		t.Logf("ResetInstance Success. %s", tea.Prettify(resp))
		time.Sleep(30 * time.Second)
	}

	// StopInstance
	// {
	// 	resp, err := TencentIo.ModifyInstance("tencent", "ap-shanghai", model.ModifyInstanceInput{
	// 		Action:      model.StopInstance,
	// 		InstanceIDs: instancesIds,
	// 	})
	// 	if err != nil {
	// 		t.Error(err)
	// 		return
	// 	}
	// 	t.Logf("StopInstance Success. %s", tea.Prettify(resp))
	// 	time.Sleep(30 * time.Second)
	// }

	// ChangeInstanceType
	// {
	// 	resp, err := TencentIo.ModifyInstance("tencent", "ap-shanghai", model.ModifyInstanceInput{
	// 		Action:       model.ChangeInstanceType,
	// 		InstanceIDs:  instancesIds,
	// 		InstanceType: tea.String("SA5.MEDIUM2"),
	// 	})
	// 	if err != nil {
	// 		t.Error(err)
	// 		return
	// 	}
	// 	t.Logf("ChangeInstanceType Success. %s", tea.Prettify(resp))
	// }
}

func TestChangeInstanceType(t *testing.T) {
	resp, err := TencentIo.ModifyInstance("tencent", "ap-shanghai", model.ModifyInstanceInput{
		Action:       model.ChangeInstanceType,
		InstanceIDs:  []*string{tea.String("ins-k7fdkyi1")},
		InstanceType: tea.String("SA5.2XLARGE32"),
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("Success. %s", tea.Prettify(resp))
}

// TestResetInstance
func TestResetInstance(t *testing.T) {
	resp, err := TencentIo.ModifyInstance("tencent", "ap-shanghai", model.ModifyInstanceInput{
		Action:      model.ResetInstance,
		InstanceIDs: []*string{tea.String("ins-xx")},
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("Success. %s", tea.Prettify(resp))
}

func TestDeleteInstance(t *testing.T) {
	resp, err := TencentIo.DeleteInstance("tencent", "ap-shanghai", model.DeleteInstanceInput{
		InstanceIds: []*string{tea.String("ins-xx")},
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("Success. %s", tea.Prettify(resp))
}

// TEST CreateSecurityGroupWithPolicies
func TestCreateSecurityGroupWithPolicies(t *testing.T) {
	resp, err := TencentIo.CreateSecurityGroupWithPolicies("tencent", "ap-beijing", model.CreateSecurityGroupWithPoliciesInput{
		GroupName:        tea.String("office-test"),
		GroupDescription: tea.String("multi-cloud-sdk-test"),
		PolicySet: model.PolicySet{
			Egress: []model.SecurityGroupPolicy{
				{
					Protocol:          tea.String("ALL"),
					Port:              tea.String("ALL"),
					CidrBlock:         tea.String("0.0.0.0/0"),
					PolicyDescription: tea.String("allow all"),
					Action:            tea.String("ACCEPT"),
				},
			},
			Ingress: []model.SecurityGroupPolicy{
				{
					Protocol:          tea.String("ALL"),
					Port:              tea.String("ALL"),
					CidrBlock:         tea.String("0.0.0.0/0"),
					PolicyDescription: tea.String("allow all"),
					Action:            tea.String("ACCEPT"),
				},
			},
		},
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("Success. %s", tea.Prettify(resp))
}
