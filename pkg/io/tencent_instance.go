package io

import (
	"fmt"
	"time"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/spf13/cast"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	"github.com/xops-infra/multi-cloud-sdk/pkg/model"
)

// Instance
func (c *tencentClient) DescribeInstances(profile, region string, input model.DescribeInstancesInput) (model.InstanceResponse, error) {
	instances := make(map[string]model.Instance, 0)
	client, err := c.io.GetTencentCvmClient(profile, region)
	if err != nil {
		return model.InstanceResponse{}, err
	}
	request := cvm.NewDescribeInstancesRequest()

	request.InstanceIds = input.InstanceIds
	for _, filter := range input.Filters {
		request.Filters = append(request.Filters, &cvm.Filter{
			Name:   filter.Name,
			Values: filter.Values,
		})
	}
	var pageSize int64 = 20
	if input.Size != nil {
		pageSize = *input.Size
		request.Limit = input.Size
	}

	response, err := client.DescribeInstances(request)
	if err != nil {
		return model.InstanceResponse{}, err
	}

	total_cvm := *response.Response.TotalCount
	pages_all := total_cvm / pageSize
	pages := *common.Int64Ptr(0)
	for pages <= pages_all {
		request.Limit = common.Int64Ptr(pageSize)
		if pages > 0 {
			request.Offset = common.Int64Ptr(pageSize*pages - 1)
		}
		response, err := client.DescribeInstances(request)
		if err != nil {
			return model.InstanceResponse{}, err
		}
		for _, instanceSet := range response.Response.InstanceSet {
			if _, ok := instances[*instanceSet.InstanceId]; ok {
				continue
			}
			createdTime, err := time.Parse(time.RFC3339, tea.StringValue(instanceSet.CreatedTime))
			if err != nil {
				return model.InstanceResponse{}, err
			}
			instances[*instanceSet.InstanceId] = model.Instance{
				Profile:            profile,
				KeyIDs:             instanceSet.LoginSettings.KeyIds,
				Zone:               instanceSet.Placement.Zone,
				InstanceID:         instanceSet.InstanceId,
				Name:               instanceSet.InstanceName,
				Region:             &region,
				Status:             model.ToInstanceStatus(*instanceSet.InstanceState),
				PublicIP:           instanceSet.PublicIpAddresses,
				PrivateIP:          instanceSet.PrivateIpAddresses,
				Tags:               model.TencentTagsToModelTags(instanceSet.Tags),
				Owner:              model.TencentTagsToModelTags(instanceSet.Tags).GetOwner(),
				InstanceType:       instanceSet.InstanceType,
				Platform:           instanceSet.OsName,
				InstanceChargeType: instanceSet.InstanceChargeType,
				CreatedTime:        &createdTime,
			}
		}
		pages = pages + 1
	}
	instance := make([]model.Instance, 0)
	for _, v := range instances {
		instance = append(instance, v)
	}

	return model.InstanceResponse{
		Instances:  instance,
		NextMarker: nil,
	}, nil
}

func (c *tencentClient) CreateInstance(profile, region string, input model.CreateInstanceInput) (model.CreateInstanceResponse, error) {
	client, err := c.io.GetTencentCvmClient(profile, region)
	if err != nil {
		return model.CreateInstanceResponse{}, err
	}
	response, err := client.RunInstances(input.ToTencentRunInstancesRequest())
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return model.CreateInstanceResponse{}, fmt.Errorf("an api error has returned: %s", err)
	}
	if err != nil {
		return model.CreateInstanceResponse{}, err
	}
	return model.CreateInstanceResponse{
		Meta:        response.ToJsonString(),
		InstanceIds: response.Response.InstanceIdSet,
	}, nil
}

// 查询可用区列表
func (c *tencentClient) QueryRegions(profile, region string) (*cvm.DescribeZonesResponse, error) {
	svc, err := c.io.GetTencentCvmClient(profile, region)
	if err != nil {
		return nil, err
	}
	// 实例化一个请求对象,每个接口都会对应一个request对象
	request := cvm.NewDescribeZonesRequest()
	// 返回的resp是一个DescribeZonesResponse的实例，与请求对象对应
	response, err := svc.DescribeZones(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return nil, fmt.Errorf("an api error has returned: %s", err)
	}
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (c *tencentClient) ModifyInstance(profile, region string, input model.ModifyInstanceInput) (model.ModifyInstanceResponse, error) {
	instanceIDs := make([]*string, 0)
	for _, instanceID := range input.InstanceIDs {
		if instanceID == "" {
			continue
		}
		instanceIDs = append(instanceIDs, tea.String(instanceID))
	}
	switch input.Action {
	case model.StartInstance:
		return c.StartInstance(profile, region, instanceIDs)
	case model.StopInstance:
		return c.StopInstance(profile, region, instanceIDs)
	case model.RebootInstance:
		return c.RebootInstance(profile, region, instanceIDs)
	case model.ResetInstance:
		return c.ResetInstance(profile, region, instanceIDs)
	case model.ChangeInstanceType:
		if input.InstanceType == nil {
			return model.ModifyInstanceResponse{}, fmt.Errorf("instance type is required")
		}
		return c.ChangeInstanceType(profile, region, instanceIDs, input.InstanceType)
	case model.ChangeInstanceTags:
		if input.ModifyTagsInput == nil {
			return model.ModifyInstanceResponse{}, fmt.Errorf("modify tags input is required")
		}
		for _, instanceID := range input.InstanceIDs {
			err := c.ModifyTagsForResource(profile, region, model.ModifyTagsInput{
				InstanceId: tea.String(instanceID),
				Key:        input.ModifyTagsInput.Key,
				Value:      input.ModifyTagsInput.Value,
			})
			if err != nil {
				return model.ModifyInstanceResponse{}, err
			}
		}
		return model.ModifyInstanceResponse{}, nil
	default:
		return model.ModifyInstanceResponse{}, fmt.Errorf("unsupported action: %s", input.Action)
	}
}

func (c *tencentClient) StartInstance(profile, region string, instances []*string) (model.ModifyInstanceResponse, error) {
	client, err := c.io.GetTencentCvmClient(profile, region)
	if err != nil {
		return model.ModifyInstanceResponse{}, err
	}
	request := cvm.NewStartInstancesRequest()
	request.InstanceIds = instances
	response, err := client.StartInstances(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return model.ModifyInstanceResponse{}, fmt.Errorf("an api error has returned: %s", err)
	}
	if err != nil {
		return model.ModifyInstanceResponse{}, err
	}
	return model.ModifyInstanceResponse{
		Meta: response.ToJsonString(),
	}, nil
}

func (c *tencentClient) StopInstance(profile, region string, instances []*string) (model.ModifyInstanceResponse, error) {
	client, err := c.io.GetTencentCvmClient(profile, region)
	if err != nil {
		return model.ModifyInstanceResponse{}, err
	}
	request := cvm.NewStopInstancesRequest()
	request.InstanceIds = instances
	response, err := client.StopInstances(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return model.ModifyInstanceResponse{}, fmt.Errorf("an api error has returned: %s", err)
	}
	if err != nil {
		return model.ModifyInstanceResponse{}, err
	}
	return model.ModifyInstanceResponse{
		Meta: response.ToJsonString(),
	}, nil
}

func (c *tencentClient) RebootInstance(profile, region string, instances []*string) (model.ModifyInstanceResponse, error) {
	client, err := c.io.GetTencentCvmClient(profile, region)
	if err != nil {
		return model.ModifyInstanceResponse{}, err
	}
	request := cvm.NewRebootInstancesRequest()
	request.InstanceIds = instances
	response, err := client.RebootInstances(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return model.ModifyInstanceResponse{}, fmt.Errorf("an api error has returned: %s", err)
	}
	if err != nil {
		return model.ModifyInstanceResponse{}, err
	}
	return model.ModifyInstanceResponse{
		Meta: response.ToJsonString(),
	}, nil
}

func (c *tencentClient) ResetInstance(profile, region string, instanceIDs []*string) (model.ModifyInstanceResponse, error) {
	resp, err := c.DescribeInstances(profile, region, model.DescribeInstancesInput{
		InstanceIds: instanceIDs,
	})
	if err != nil {
		return model.ModifyInstanceResponse{}, err
	}
	if len(resp.Instances) != 1 {
		return model.ModifyInstanceResponse{}, fmt.Errorf("only support reset one instance at a time")
	}
	instance := resp.Instances[0]
	client, err := c.io.GetTencentCvmClient(profile, region)
	if err != nil {
		return model.ModifyInstanceResponse{}, err
	}

	request := cvm.NewResetInstanceRequest()
	request.InstanceId = instance.InstanceID
	request.LoginSettings = &cvm.LoginSettings{
		KeyIds: instance.KeyIDs,
		// KeepImageLogin: tea.String("TRUE"), // 不支持公有镜像
	}
	response, err := client.ResetInstance(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return model.ModifyInstanceResponse{}, fmt.Errorf("an api error has returned: %s", err)
	}
	if err != nil {
		return model.ModifyInstanceResponse{}, err
	}
	return model.ModifyInstanceResponse{
		Meta: response.ToJsonString(),
	}, nil
}

// 默认关闭强制关机
func (c *tencentClient) ChangeInstanceType(profile, region string, instances []*string, instanceType *string) (model.ModifyInstanceResponse, error) {
	client, err := c.io.GetTencentCvmClient(profile, region)
	if err != nil {
		return model.ModifyInstanceResponse{}, err
	}
	request := cvm.NewResetInstancesTypeRequest()
	request.InstanceIds = instances
	request.InstanceType = instanceType
	request.ForceStop = tea.Bool(false)
	response, err := client.ResetInstancesType(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return model.ModifyInstanceResponse{}, fmt.Errorf("an api error has returned: %s", err)
	}
	if err != nil {
		return model.ModifyInstanceResponse{}, err
	}
	return model.ModifyInstanceResponse{
		Meta: response.ToJsonString(),
	}, nil
}

func (c *tencentClient) DeleteInstance(profile, region string, input model.DeleteInstanceInput) (model.DeleteInstanceResponse, error) {
	client, err := c.io.GetTencentCvmClient(profile, region)
	if err != nil {
		return model.DeleteInstanceResponse{}, err
	}
	response, err := client.TerminateInstances(input.ToTencentTerminateInstancesRequest())
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return model.DeleteInstanceResponse{}, fmt.Errorf("an api error has returned: %s", err)
	}
	if err != nil {
		return model.DeleteInstanceResponse{}, err
	}
	return model.DeleteInstanceResponse{
		Meta: response.ToJsonString(),
	}, nil
}

func (c *tencentClient) DescribeKeyPairs(profile, region string, input model.CommonFilter) ([]model.KeyPair, error) {
	client, err := c.io.GetTencentCvmClient(profile, region)
	if err != nil {
		return nil, err
	}
	keyPairs := make([]model.KeyPair, 0)
	request := cvm.NewDescribeKeyPairsRequest()
	// request.Limit = tea.Int64(2)
	for {
		response, err := client.DescribeKeyPairs(request)
		if _, ok := err.(*errors.TencentCloudSDKError); ok {
			return nil, fmt.Errorf("an api error has returned: %s", err)
		}
		if err != nil {
			return nil, err
		}
		for _, keyPair := range response.Response.KeyPairSet {
			keyPairs = append(keyPairs, model.KeyPair{
				ID:        *keyPair.KeyId,
				Name:      *keyPair.KeyName,
				PublicKey: *keyPair.PublicKey,
			})
		}
		if response.Response.TotalCount != nil && *response.Response.TotalCount == int64(len(keyPairs)) {
			break
		}
		request.Offset = tea.Int64(cast.ToInt64(len(keyPairs)))
		// fmt.Println(cast.ToInt64(len(keyPairs)))
	}
	return keyPairs, nil
}

func (c *tencentClient) DescribeImages(profile, region string, input model.CommonFilter) ([]model.Image, error) {
	client, err := c.io.GetTencentCvmClient(profile, region)
	if err != nil {
		return nil, err
	}
	images := make([]model.Image, 0)
	request := cvm.NewDescribeImagesRequest()
	// request.Limit = tea.Uint64(10)
	for {
		response, err := client.DescribeImages(request)
		if _, ok := err.(*errors.TencentCloudSDKError); ok {
			return nil, fmt.Errorf("an api error has returned: %s", err)
		}
		if err != nil {
			return nil, err
		}
		for _, image := range response.Response.ImageSet {
			images = append(images, model.Image{
				ID:       *image.ImageId,
				Name:     *image.ImageName,
				Arch:     *image.Architecture,
				Platform: *image.Platform,
			})
		}
		if response.Response.TotalCount != nil && *response.Response.TotalCount == int64(len(images)) {
			break
		}
		request.Offset = tea.Uint64(cast.ToUint64(len(images)))
		// fmt.Println(cast.ToUint64(len(images)))
	}
	return images, nil
}

func (c *tencentClient) DescribeInstanceTypes(profile, region string) ([]model.InstanceType, error) {
	client, err := c.io.GetTencentCvmClient(profile, region)
	if err != nil {
		return nil, err
	}
	request := cvm.NewDescribeInstanceTypeConfigsRequest()
	response, err := client.DescribeInstanceTypeConfigs(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return nil, fmt.Errorf("an api error has returned: %s", err)
	}
	if err != nil {
		return nil, err
	}
	instanceTypes := make([]model.InstanceType, 0)
	for _, instanceType := range response.Response.InstanceTypeConfigSet {
		instanceTypes = append(instanceTypes, model.InstanceType{
			Type: tea.StringValue(instanceType.InstanceType),
			CPU:  tea.Int64Value(instanceType.CPU),
			Mem:  tea.Int64Value(instanceType.Memory),
		})
	}
	return instanceTypes, nil
}

func (t *tencentClient) DescribeInstancePrice(profile, region string, input model.DescribeInstancePriceInput) (model.DescribeInstancePriceResponse, error) {
	client, err := t.io.GetTencentCvmClient(profile, region)
	if err != nil {
		return model.DescribeInstancePriceResponse{}, err
	}

	request := cvm.NewInquiryPriceRunInstancesRequest()
	request.Placement = &cvm.Placement{
		Zone: input.Zone,
	}
	request.ImageId = input.ImageId
	request.InstanceChargeType = common.StringPtr("PREPAID")
	request.InstanceChargePrepaid = &cvm.InstanceChargePrepaid{
		Period: input.Period,
	}
	request.InstanceType = input.InstanceType
	request.SystemDisk = &cvm.SystemDisk{
		DiskType: input.SystemDisk.Type,
		DiskSize: input.SystemDisk.Size,
	}
	for _, dataDisk := range input.DataDisks {
		request.DataDisks = append(request.DataDisks, &cvm.DataDisk{
			DiskType: dataDisk.Type,
			DiskSize: dataDisk.Size,
		})
	}

	response, err := client.InquiryPriceRunInstances(request)
	if err != nil {
		return model.DescribeInstancePriceResponse{}, err
	}

	return model.DescribeInstancePriceResponse{
		OriginalPrice: response.Response.Price.InstancePrice.OriginalPrice,
		DiscountPrice: response.Response.Price.InstancePrice.DiscountPrice,
		Discount:      response.Response.Price.InstancePrice.Discount,
		Currency:      tea.String("CNY"),
	}, nil
}
