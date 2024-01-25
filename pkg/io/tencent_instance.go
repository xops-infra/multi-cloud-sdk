package io

import (
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
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
				KeyName:    instanceSet.LoginSettings.KeyIds,
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
	panic("implement me")
}

func (c *tencentClient) ModifyInstance(profile, region string, input model.ModifyInstanceInput) (model.ModifyInstanceResponse, error) {
	panic("implement me")
}

func (c *tencentClient) DeleteInstance(profile, region string, input model.DeleteInstanceInput) (model.DeleteInstanceResponse, error) {
	panic("implement me")
}
