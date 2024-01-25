package io

import (
	"strings"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/xops-infra/multi-cloud-sdk/pkg/model"
)

func (c *awsClient) DescribeInstances(profile, region string, input model.DescribeInstancesInput) (model.InstanceResponse, error) {
	svc, err := c.io.GetAwsEc2Client(profile, region)
	if err != nil {
		return model.InstanceResponse{}, err
	}
	req := &ec2.DescribeInstancesInput{}
	if input.InstanceIds != nil {
		req.InstanceIds = input.InstanceIds
	}
	if input.Filters != nil {
		for _, filter := range input.Filters {
			req.Filters = append(req.Filters, &ec2.Filter{
				Name:   filter.Name,
				Values: filter.Values,
			})
		}
	}

	if input.NextMarker != nil {
		req.NextToken = input.NextMarker
	}

	if input.Size != nil {
		req.MaxResults = input.Size
	}

	out, err := svc.DescribeInstances(req)
	if err != nil {
		return model.InstanceResponse{}, err
	}
	var instances []model.Instance
	for _, reservation := range out.Reservations {
		for _, instance := range reservation.Instances {
			tags := model.AwsTagsToModelTags(instance.Tags)
			instances = append(instances, model.Instance{
				Profile:    profile,
				KeyName:    []*string{instance.KeyName},
				InstanceID: instance.InstanceId,
				Name:       tags.GetName(),
				Region:     instance.Placement.AvailabilityZone,
				Status:     model.ToInstanceStatus(strings.ToUpper(*instance.State.Name)),
				PublicIP:   []*string{instance.PublicIpAddress},
				PrivateIP:  []*string{instance.PrivateIpAddress},
				Tags:       tags,
				Owner:      tags.GetOwner(),
				Platform:   instance.PlatformDetails,
			})
		}
	}
	return model.InstanceResponse{
		Instances:  instances,
		NextMarker: out.NextToken,
	}, nil
}

func (c *awsClient) CreateInstance(profile, region string, input model.CreateInstanceInput) (model.CreateInstanceResponse, error) {
	panic("implement me")
}

func (c *awsClient) ModifyInstance(profile, region string, input model.ModifyInstanceInput) (model.ModifyInstanceResponse, error) {
	panic("implement me")
}

func (c *awsClient) DeleteInstance(profile, region string, input model.DeleteInstanceInput) (model.DeleteInstanceResponse, error) {
	panic("implement me")
}
