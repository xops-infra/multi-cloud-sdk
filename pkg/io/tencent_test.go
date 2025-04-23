package io_test

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
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
		}, {
			Name:  "tx-dev",
			Cloud: model.TENCENT,
			AK:    os.Getenv("TX_DEV_ID"),
			SK:    os.Getenv("TX_DEV_KEY"),
		},
	}
	clientIo := io.NewCloudClient(profiles)
	TencentIo = io.NewTencentClient(clientIo)
}

// test DescribeInstancePrice
func TestDescribeInstancePrice(t *testing.T) {
	price, err := TencentIo.DescribeInstancePrice("tencent", "ap-beijing", model.DescribeInstancePriceInput{
		InstanceType: tea.String("S5.6XLARGE96"),
		ImageId:      tea.String("img-1u6l2i9l"),
		Period:       tea.Int64(1),
		Zone:         tea.String("ap-beijing-7"),
		SystemDisk: &model.Disk{
			Size: tea.Int64(100),
		},
		// DataDisks: []model.Disk{
		// 	{Size: tea.Int64(100)},
		// 	{Size: tea.Int64(100)},
		// 	{Size: tea.Int64(100)},
		// },
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("Success. %s", tea.Prettify(price))
}

// test ModifyTagsForResource
func TestModifyTagsForResource(t *testing.T) {
	err := TencentIo.ModifyTagsForResource("tencent", "ap-shanghai", model.ModifyTagsInput{
		InstanceId: tea.String("ins-1oy0zn7n"),
		Key:        tea.String("ExpireTime"),
		Value:      tea.String("2025123112"),
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("Success.")
}

// Test AddTagsToResource
func TestAddTagsToResource(t *testing.T) {
	err := TencentIo.AddTagsToResource("tencent", "ap-shanghai", model.AddTagsInput{
		InstanceIds: []*string{tea.String("ins-hs9y7cqv")},
		Tags: model.Tags{
			{
				Key:   "ExpireTime",
				Value: "2025123122",
			},
		},
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("Success.")
}

// TEST QuerySubnet
func TestQuerySubnet(t *testing.T) {
	timeStart := time.Now()
	subnets, err := TencentIo.QuerySubnet("tencent", "ap-shanghai", model.CommonFilter{})
	if err != nil {
		t.Error(err)
		return
	}
	for _, subnet := range subnets {
		fmt.Println(tea.Prettify(subnet))
	}
	t.Log("Success.", time.Since(timeStart), len(subnets))
}

func TestDescribeInstanceTypes(t *testing.T) {
	timeStart := time.Now()
	instanceTypes, err := TencentIo.DescribeInstanceTypes("tencent", "ap-shanghai")
	if err != nil {
		t.Error(err)
		return
	}
	for _, instanceType := range instanceTypes {
		fmt.Println(tea.Prettify(instanceType))
	}
	t.Log("Success.", time.Since(timeStart), len(instanceTypes))
}
func TestDescribeKeyPairs(t *testing.T) {
	timeStart := time.Now()
	keyPairs, err := TencentIo.DescribeKeyPairs("tencent", "ap-shanghai", model.CommonFilter{})
	if err != nil {
		t.Error(err)
		return
	}
	for _, keyPair := range keyPairs {
		fmt.Println(tea.Prettify(keyPair))
	}
	t.Log("Success.", time.Since(timeStart), len(keyPairs))
}

func TestDescribeImages(t *testing.T) {
	timeStart := time.Now()
	images, err := TencentIo.DescribeImages("tencent", "ap-shanghai", model.CommonFilter{})
	if err != nil {
		t.Error(err)
		return
	}
	for _, image := range images {
		fmt.Println(tea.Prettify(image))
	}
	t.Log("Success.", time.Since(timeStart), len(images))
}

func TestQuerySecurityGroups(t *testing.T) {
	timeStart := time.Now()
	filter := model.CommonFilter{}
	securityGroups, err := TencentIo.QuerySecurityGroups("tencent", "ap-shanghai", filter)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("Success.", time.Since(timeStart), len(securityGroups))
}

func TestQueryTencentEmrCluster(t *testing.T) {
	timeStart := time.Now()
	filter := model.EmrFilter{
		Profile: tea.String("tencent"),
		Region:  tea.String("na-ashburn"),
		ClusterStates: []model.EMRClusterStatus{
			model.EMRClusterRunning,
		},
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
		Region:  tea.String("ap-shanghai"),
		// IDS:     []*string{tea.String("emr-alhn4h4s")},
	})
	if err != nil {
		t.Error(err)
		return
	}
	for _, instance := range instances {
		fmt.Println(tea.Prettify(instance))
	}
	t.Log("Success.", time.Since(timeStart), len(instances), instances)
}

func TestCreateEmrCluster(t *testing.T) {
	timeStart := time.Now()
	input := model.CreateEmrClusterInput{
		Name: tea.String("test"),
		Tags: model.Tags{
			{
				Key:   "Owner",
				Value: "zhoushoujian",
			},
		},
		InstanceChargeType: model.POSTPAID_BY_HOUR.String(),
		APPs:               []*string{tea.String("Spark")},
		ResourceSpec: &model.ResourceSpec{
			HA:     tea.Bool(false),
			VPC:    tea.String("vpc-gjljk6e8"),
			Subnet: tea.String("subnet-j94dsqaj"),
			SgId:   tea.String("sg-2qt3di24"),
			KeyID:  tea.String("skey-gyqojq9d"),
			MasterResourceSpec: &model.EMRInstaceSpec{
				InstanceCount: tea.Int64(1),
				InstanceType:  tea.String("TS5.2XLARGE32"),
				DiskType:      tea.String("CLOUD_SSD"),
				DiskNum:       tea.Int64(0),
				DiskSize:      tea.Int64(40),
				RootSize:      tea.Int64(40),
			},
			CoreResourceSpec: &model.EMRInstaceSpec{
				InstanceCount: tea.Int64(2),
				InstanceType:  tea.String("TS5.2XLARGE32"),
				DiskType:      tea.String("CLOUD_SSD"),
				DiskNum:       tea.Int64(0),
				DiskSize:      tea.Int64(40),
				RootSize:      tea.Int64(40),
			},
		},
	}
	instances, err := TencentIo.CreateEmrCluster("tencent", "ap-shanghai", input)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("Success.", time.Since(timeStart), instances)
}

func TestDescribeVolumes(t *testing.T) {
	timeStart := time.Now()
	volumes, err := TencentIo.DescribeVolumes("tencent", "ap-beijing", model.DescribeVolumesInput{
		VolumeIDs: []string{"disk-703vadur"},
	})
	if err != nil {
		t.Error(err)
		return
	}
	for _, volume := range volumes {
		fmt.Println(tea.Prettify(volume))
	}
	t.Log("Success.", time.Since(timeStart), len(volumes))
}

func TestListInstance(t *testing.T) {
	timeStart := time.Now()
	filter := model.InstanceFilter{}
	instances, err := TencentIo.DescribeInstances("tencent", "ap-beijing", filter.ToTxDescribeInstancesInput())
	if err != nil {
		t.Error(err)
		return
	}
	for _, instance := range instances.Instances {
		fmt.Println(*instance.Zone, *instance.CreatedTime, tea.Prettify(instance.InstanceChargeType))
	}
	t.Log("Success.", time.Since(timeStart), len(instances.Instances))
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

func TestChangeInstanceType(t *testing.T) {
	resp, err := TencentIo.ModifyInstance("tencent", "ap-shanghai", model.ModifyInstanceInput{
		Action:       model.ChangeInstanceType,
		InstanceIDs:  []string{"ins-k7fdkyi1"},
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
		InstanceIDs: []string{"ins-xx"},
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
		GroupName:        tea.String("-test"),
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

func TestCreateSecurityGroupWithPolicies1(t *testing.T) {

	// 读取本地json文件
	file, err := os.Open("/tmp/sg.json")
	if err != nil {
		t.Error(err)
		return
	}
	defer file.Close()

	// 解析json文件
	var sgs []map[string]string
	err = json.NewDecoder(file).Decode(&sgs)
	if err != nil {
		t.Error(err)
		return
	}

	// 创建sg
	var ingress []model.SecurityGroupPolicy
	for _, sg := range sgs {
		ingress = append(ingress, model.SecurityGroupPolicy{
			Protocol:          tea.String("TCP"),
			Port:              tea.String("21,22"),
			CidrBlock:         tea.String(sg["CidrIp"]),
			PolicyDescription: tea.String(sg["Description"]),
			Action:            tea.String("ACCEPT"),
		})
	}
	resp, err := TencentIo.CreateSecurityGroupWithPolicies("tencent", "na-ashburn", model.CreateSecurityGroupWithPoliciesInput{
		GroupName:        tea.String("0066_aws_cp_ftp"),
		GroupDescription: tea.String("multi-cloud-sdk"),
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
			Ingress: ingress,
		},
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("Success. %s", tea.Prettify(resp))
}

// TEST CreateSecurityGroupPolicies
func TestCreateSecurityGroupPolicies(t *testing.T) {
	allowAll := strings.Split(os.Getenv("TestCreateSecurityGroupPoliciesCidr"), ",")
	ingress := []model.SecurityGroupPolicy{}
	for _, cidr := range allowAll {
		if cidr == "" {
			continue
		}
		ingress = append(ingress, model.SecurityGroupPolicy{
			Protocol:          tea.String("ALL"),
			Port:              tea.String("ALL"),
			CidrBlock:         tea.String(cidr),
			PolicyDescription: tea.String("allow all for office" + "(mcs)"),
			Action:            tea.String("ACCEPT"),
		})
	}
	fmt.Println(tea.Prettify(ingress))
	resp, err := TencentIo.CreateSecurityGroupPolicies("tencent", "ap-beijing", model.CreateSecurityGroupPoliciesInput{
		SecurityGroupId: tea.String("sg-xxx"),
		PolicySet: model.PolicySet{
			Ingress: ingress,
		},
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("Success. %s", tea.Prettify(resp))
}
