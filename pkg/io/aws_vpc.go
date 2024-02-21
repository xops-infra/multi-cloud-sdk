package io

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/xops-infra/multi-cloud-sdk/pkg/model"
)

// QueryVpcs
func (c *awsClient) QueryVPC(profile, region string, input model.CommonFilter) ([]model.VPC, error) {
	svc, err := c.io.GetAwsEc2Client(profile, region)
	if err != nil {
		return nil, err
	}
	var vpcs []model.VPC
	_input := &ec2.DescribeVpcsInput{}
	if input.ID != "" {
		_input.VpcIds = []*string{aws.String(input.ID)}
	}
	for {
		out, err := svc.DescribeVpcs(_input)
		if err != nil {
			return nil, err
		}
		for _, vpc := range out.Vpcs {
			vpcs = append(vpcs, model.VPC{
				ID:            aws.StringValue(vpc.VpcId),
				Tags:          model.AwsTagsToModelTags(vpc.Tags),
				Region:        region,
				CloudProvider: model.AWS,
				Account:       profile,
				IsDefault:     aws.BoolValue(vpc.IsDefault),
				CidrBlock:     aws.StringValue(vpc.CidrBlock),
			})
		}
		if out.NextToken == nil {
			break
		}
		_input.NextToken = out.NextToken
	}
	return vpcs, nil
}

// QuerySubnet
func (c *awsClient) QuerySubnet(profile, region string, input model.CommonFilter) ([]model.Subnet, error) {
	svc, err := c.io.GetAwsEc2Client(profile, region)
	if err != nil {
		return nil, err
	}
	var subnets []model.Subnet
	_input := &ec2.DescribeSubnetsInput{}
	if input.ID != "" {
		_input.SubnetIds = []*string{aws.String(input.ID)}
	}
	for {
		out, err := svc.DescribeSubnets(_input)
		if err != nil {
			return nil, err
		}
		for _, subnet := range out.Subnets {
			tags := model.AwsTagsToModelTags(subnet.Tags)
			subnets = append(subnets, model.Subnet{
				ID:            subnet.SubnetId,
				Tags:          tags,
				Name:          tags.GetName(),
				Region:        region,
				CloudProvider: model.AWS,
				Account:       profile,
				CidrBlock:     subnet.CidrBlock,
				VpcID:         subnet.VpcId,
				Zone:          subnet.AvailabilityZone,
				IsDefault:     subnet.DefaultForAz,
				// CreatedTime:   ,
				AvailableIpAddressCount: aws.Int64Value(subnet.AvailableIpAddressCount),
				// NetworkAclId:            aws.StringValue(subnet.ac),
				// RouteTableId: aws.StringValue(subnet.RouteTableId),
			})
		}
		if out.NextToken == nil {
			break
		}
		_input.NextToken = out.NextToken
	}
	return subnets, nil
}

// QueryEIP
func (c *awsClient) QueryEIP(profile, region string, input model.CommonFilter) ([]model.EIP, error) {
	svc, err := c.io.GetAwsEc2Client(profile, region)
	if err != nil {
		return nil, err
	}
	var eips []model.EIP
	_input := &ec2.DescribeAddressesInput{}
	if input.ID != "" {
		_input.AllocationIds = []*string{aws.String(input.ID)}
	}
	out, err := svc.DescribeAddresses(_input)
	if err != nil {
		return nil, err
	}
	for _, address := range out.Addresses {
		tags := model.AwsTagsToModelTags(address.Tags)
		eips = append(eips, model.EIP{
			ID:            address.AllocationId,
			Tags:          tags,
			Name:          tags.GetName(),
			Region:        region,
			CloudProvider: model.AWS,
			Account:       profile,
			// Status:             aws.StringValue(address.),
			AddressIp:          address.PublicIp,
			InstanceId:         address.InstanceId,
			NetworkInterfaceId: aws.StringValue(address.NetworkInterfaceId),
			PrivateAddressIp:   aws.StringValue(address.PrivateIpAddress),
			// Bandwidth:          aws.Int64Value(address.Bandwidth),
			// InternetChargeType: aws.StringValue(address.InternetChargeType),
		})
	}
	return eips, nil
}

// QueryNAT
func (c *awsClient) QueryNAT(profile, region string, input model.CommonFilter) ([]model.NAT, error) {
	svc, err := c.io.GetAwsEc2Client(profile, region)
	if err != nil {
		return nil, err
	}
	var nats []model.NAT
	_input := &ec2.DescribeNatGatewaysInput{}
	if input.ID != "" {
		_input.NatGatewayIds = []*string{aws.String(input.ID)}
	}
	out, err := svc.DescribeNatGateways(_input)
	if err != nil {
		return nil, err
	}
	for _, nat := range out.NatGateways {
		tags := model.AwsTagsToModelTags(nat.Tags)
		nats = append(nats, model.NAT{
			ID:            aws.StringValue(nat.NatGatewayId),
			Tags:          tags,
			Name:          *tags.GetName(),
			Region:        region,
			CloudProvider: model.AWS,
			Account:       profile,
			VpcID:         aws.StringValue(nat.VpcId),
			CreatedTime:   *nat.CreateTime,
			Status:        aws.StringValue(nat.State),
			// Zone:        aws.StringValue(nat.AvailabilityZone),
			SubnetID: aws.StringValue(nat.SubnetId),
		})
	}
	return nats, nil
}

// CreateSecurityGroupWithPolicies
func (c *awsClient) CreateSecurityGroupWithPolicies(profile, region string, input model.CreateSecurityGroupWithPoliciesInput) (model.CreateSecurityGroupWithPoliciesResponse, error) {
	panic("implement me")
}

func (c *awsClient) CreateSecurityGroupPolicies(profile, region string, input model.CreateSecurityGroupPoliciesInput) (model.CreateSecurityGroupPoliciesResponse, error) {
	panic("implement me")
}
