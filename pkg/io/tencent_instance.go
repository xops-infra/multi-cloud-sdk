package io

import (
	"fmt"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	"github.com/xops-infra/multi-cloud-sdk/pkg/model"
)

// Instance
func (c *tencentClient) DescribeInstances(profile, region string, input model.DescribeInstancesInput) (model.InstanceResponse, error) {
	var instances []model.Instance

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
			instances = append(instances, model.Instance{
				Profile:    profile,
				KeyIDs:     instanceSet.LoginSettings.KeyIds,
				InstanceID: instanceSet.InstanceId,
				Name:       instanceSet.InstanceName,
				Region:     instanceSet.Placement.Zone,
				Status:     model.ToInstanceStatus(*instanceSet.InstanceState),
				PublicIP:   instanceSet.PublicIpAddresses,
				PrivateIP:  instanceSet.PrivateIpAddresses,
				Tags:       model.TencentTagsToModelTags(instanceSet.Tags),
				Owner:      model.TencentTagsToModelTags(instanceSet.Tags).GetOwner(),
				Platform:   instanceSet.OsName,
			})
		}
		pages = pages + 1
	}

	return model.InstanceResponse{
		Instances:  instances,
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
	switch input.Action {
	case model.StartInstance:
		return c.StartInstance(profile, region, input.InstanceIDs)
	case model.StopInstance:
		return c.StopInstance(profile, region, input.InstanceIDs)
	case model.RebootInstance:
		return c.RebootInstance(profile, region, input.InstanceIDs)
	case model.ResetInstance:
		return c.ResetInstance(profile, region, input.InstanceIDs)
	case model.ChangeInstanceType:
		if input.InstanceType == nil {
			return model.ModifyInstanceResponse{}, fmt.Errorf("instance type is required")
		}
		return c.ChangeInstanceType(profile, region, input.InstanceIDs, input.InstanceType)
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
