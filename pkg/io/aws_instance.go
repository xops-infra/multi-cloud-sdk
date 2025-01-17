package io

import (
	"fmt"
	"strings"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/aws/aws-sdk-go/aws"
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

	var instances []model.Instance

	var out *ec2.DescribeInstancesOutput

	for {
		// 打印 out 占用内存
		out, err = svc.DescribeInstances(req)
		if err != nil {
			return model.InstanceResponse{}, err
		}
		for _, reservation := range out.Reservations {
			for _, instance := range reservation.Instances {
				tags := model.AwsTagsToModelTags(instance.Tags)
				instances = append(instances, model.Instance{
					Profile:    profile,
					KeyIDs:     []*string{instance.KeyName},
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
		if out.NextToken == nil {
			break
		}
		req.NextToken = out.NextToken
	}

	out = nil

	return model.InstanceResponse{Instances: instances}, nil

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

func (c *awsClient) DescribeKeyPairs(profile, region string, input model.CommonFilter) ([]model.KeyPair, error) {
	svc, err := c.io.GetAwsEc2Client(profile, region)
	if err != nil {
		return nil, err
	}
	var keyPairs []model.KeyPair
	_input := &ec2.DescribeKeyPairsInput{}
	if input.ID != "" {
		_input.KeyPairIds = []*string{aws.String(input.ID)}
	}
	out, err := svc.DescribeKeyPairs(_input)
	if err != nil {
		return nil, err
	}
	for _, keyPair := range out.KeyPairs {
		keyPairs = append(keyPairs, model.KeyPair{
			ID:        tea.StringValue(keyPair.KeyPairId),
			Name:      tea.StringValue(keyPair.KeyName),
			PublicKey: tea.StringValue(keyPair.PublicKey),
		})
	}
	return keyPairs, nil
}

// 太多了不做获取
func (c *awsClient) DescribeImages(profile, region string, input model.CommonFilter) ([]model.Image, error) {
	svc, err := c.io.GetAwsEc2Client(profile, region)
	if err != nil {
		return nil, err
	}

	_input := &ec2.DescribeImagesInput{}
	if input.ID != "" {
		_input.ImageIds = []*string{aws.String(input.ID)}
	}
	// _input.MaxResults = tea.Int64(100)
	images := make([]model.Image, 0)
	for {
		out, err := svc.DescribeImages(_input)
		if err != nil {
			return nil, err
		}
		for _, image := range out.Images {
			images = append(images, model.Image{
				ID:       tea.StringValue(image.ImageId),
				Name:     tea.StringValue(image.Name),
				Arch:     tea.StringValue(image.Architecture),
				Platform: tea.StringValue(image.Platform),
			})
		}
		if out.NextToken == nil {
			break
		}
		_input.NextToken = out.NextToken
		fmt.Println(tea.Prettify(out.NextToken))
	}
	return images, nil
}
