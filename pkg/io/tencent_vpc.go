package io

import (
	"fmt"
	"strconv"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/spf13/cast"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	tencentVpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	"github.com/xops-infra/multi-cloud-sdk/pkg/model"
)

func (c *tencentClient) QueryVPC(profile, region string, input model.CommonFilter) ([]model.VPC, error) {
	client, err := c.io.GetTencentVpcClient(profile, region)
	if err != nil {
		return nil, err
	}
	request := tencentVpc.NewDescribeVpcsRequest()
	if input.ID != "" {
		request.Filters = []*tencentVpc.Filter{
			{
				Name:   common.StringPtr("vpc-id"),
				Values: []*string{common.StringPtr(input.ID)},
			},
		}
	}
	response, err := client.DescribeVpcs(request)
	if err != nil {
		return nil, err
	}
	var vpcs []model.VPC
	for _, vpc := range response.Response.VpcSet {
		vpcs = append(vpcs, model.VPC{
			ID:            *vpc.VpcId,
			Region:        region,
			Account:       profile,
			CloudProvider: model.TENCENT,
			Tags:          model.TencentVpcTagsFmt(vpc.TagSet),
			IsDefault:     *vpc.IsDefault,
			CidrBlock:     *vpc.CidrBlock,
		})
	}

	return vpcs, nil
}

// QuerySubnet
func (c *tencentClient) QuerySubnet(profile, region string, input model.CommonFilter) ([]model.Subnet, error) {
	client, err := c.io.GetTencentVpcClient(profile, region)
	if err != nil {
		return nil, err
	}
	request := tencentVpc.NewDescribeSubnetsRequest()
	if input.ID != "" {
		request.Filters = []*tencentVpc.Filter{
			{
				Name:   common.StringPtr("subnet-id"),
				Values: []*string{common.StringPtr(input.ID)},
			},
		}
	}
	response, err := client.DescribeSubnets(request)
	if err != nil {
		return nil, err
	}
	var subnets []model.Subnet
	for _, subnet := range response.Response.SubnetSet {
		createTime, _ := model.TimeParse(*subnet.CreatedTime)
		subnets = append(subnets, model.Subnet{
			ID:                      subnet.SubnetId,
			Region:                  region,
			Account:                 profile,
			CloudProvider:           model.TENCENT,
			Tags:                    model.TencentVpcTagsFmt(subnet.TagSet),
			VpcID:                   subnet.VpcId,
			Name:                    subnet.SubnetName,
			CidrBlock:               subnet.CidrBlock,
			IsDefault:               subnet.IsDefault,
			Zone:                    subnet.Zone,
			RouteTableId:            subnet.RouteTableId,
			CreatedTime:             &createTime,
			AvailableIpAddressCount: cast.ToInt64(subnet.AvailableIpAddressCount),
			NetworkAclId:            subnet.NetworkAclId,
		})
	}
	return subnets, nil
}

// QueryEIP
func (c *tencentClient) QueryEIP(profile, region string, input model.CommonFilter) ([]model.EIP, error) {
	client, err := c.io.GetTencentVpcClient(profile, region)
	if err != nil {
		return nil, err
	}
	request := tencentVpc.NewDescribeAddressesRequest()
	if input.ID != "" {
		request.Filters = []*tencentVpc.Filter{
			{
				Name:   common.StringPtr("address-id"),
				Values: []*string{common.StringPtr(input.ID)},
			},
		}
	}
	response, err := client.DescribeAddresses(request)
	if err != nil {
		return nil, err
	}
	var eips []model.EIP
	for _, eip := range response.Response.AddressSet {
		createTime, _ := model.TimeParse(*eip.CreatedTime)
		eips = append(eips, model.EIP{
			ID:                 eip.AddressId,
			Region:             region,
			Account:            profile,
			CloudProvider:      model.TENCENT,
			Tags:               model.TencentVpcTagsFmt(eip.TagSet),
			Name:               eip.AddressName,
			Status:             eip.AddressStatus,
			AddressIp:          eip.AddressIp,
			InstanceId:         eip.InstanceId,
			CreatedTime:        &createTime,
			Bandwidth:          tea.Int64(cast.ToInt64(eip.Bandwidth)),
			InternetChargeType: eip.InternetChargeType,
		})
	}
	return eips, nil
}

// QueryNAT
func (c *tencentClient) QueryNAT(profile, region string, input model.CommonFilter) ([]model.NAT, error) {
	client, err := c.io.GetTencentVpcClient(profile, region)
	if err != nil {
		return nil, err
	}
	request := tencentVpc.NewDescribeNatGatewaysRequest()
	if input.ID != "" {
		request.Filters = []*tencentVpc.Filter{
			{
				Name:   common.StringPtr("nat-gateway-id"),
				Values: []*string{common.StringPtr(input.ID)},
			},
		}
	}
	response, err := client.DescribeNatGateways(request)
	if err != nil {
		return nil, err
	}
	var nats []model.NAT
	for _, nat := range response.Response.NatGatewaySet {
		createTime, _ := model.TimeParse(*nat.CreatedTime)
		nats = append(nats, model.NAT{
			ID:            *nat.NatGatewayId,
			Region:        region,
			Account:       profile,
			CloudProvider: model.TENCENT,
			Tags:          model.TencentVpcTagsFmt(nat.TagSet),
			Name:          *nat.NatGatewayName,
			Status:        *nat.State,
			VpcID:         *nat.VpcId,
			Zone:          nat.Zone,
			SubnetID:      *nat.SubnetId,
			CreatedTime:   createTime,
		})
	}
	return nats, nil
}

func (c *tencentClient) CreateSecurityGroupWithPolicies(profile, region string, input model.CreateSecurityGroupWithPoliciesInput) (model.CreateSecurityGroupWithPoliciesResponse, error) {
	client, err := c.io.GetTencentVpcClient(profile, region)
	if err != nil {
		return model.CreateSecurityGroupWithPoliciesResponse{}, err
	}
	request := tencentVpc.NewCreateSecurityGroupWithPoliciesRequest()
	request.GroupName = input.GroupName
	request.GroupDescription = input.GroupDescription
	// request.ProjectId
	request.SecurityGroupPolicySet = input.PolicySet.ToTencentPolicySet()
	response, err := client.CreateSecurityGroupWithPolicies(request)
	if err != nil {
		return model.CreateSecurityGroupWithPoliciesResponse{}, err
	}
	return model.CreateSecurityGroupWithPoliciesResponse{
		Data: response,
	}, nil
}

func (c *tencentClient) CreateSecurityGroupPolicies(profile, region string, input model.CreateSecurityGroupPoliciesInput) (model.CreateSecurityGroupPoliciesResponse, error) {
	client, err := c.io.GetTencentVpcClient(profile, region)
	if err != nil {
		return model.CreateSecurityGroupPoliciesResponse{}, err
	}
	// 实例化一个请求对象,每个接口都会对应一个request对象
	request := tencentVpc.NewCreateSecurityGroupPoliciesRequest()
	request.SecurityGroupPolicySet = input.PolicySet.ToTencentPolicySet()
	request.SecurityGroupId = input.SecurityGroupId

	// 返回的resp是一个CreateSecurityGroupPoliciesResponse的实例，与请求对象对应
	response, err := client.CreateSecurityGroupPolicies(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return model.CreateSecurityGroupPoliciesResponse{}, fmt.Errorf("an api error has returned: %s", err)
	}
	if err != nil {
		return model.CreateSecurityGroupPoliciesResponse{}, err
	}
	return model.CreateSecurityGroupPoliciesResponse{
		Result: response,
	}, nil
}

func (c *tencentClient) QuerySecurityGroups(profile, region string, input model.CommonFilter) ([]model.SecurityGroup, error) {
	client, err := c.io.GetTencentVpcClient(profile, region)
	if err != nil {
		return nil, err
	}
	request := tencentVpc.NewDescribeSecurityGroupsRequest()
	// request.Limit = tea.String("10")
	securityGroups := make([]model.SecurityGroup, 0)
	for {
		response, err := client.DescribeSecurityGroups(request)
		if err != nil {
			return nil, err
		}
		for _, securityGroup := range response.Response.SecurityGroupSet {
			securityGroups = append(securityGroups, model.SecurityGroup{
				ID:   *securityGroup.SecurityGroupId,
				Name: *securityGroup.SecurityGroupName,
			})
		}
		if response.Response.TotalCount != nil && *response.Response.TotalCount == uint64(len(securityGroups)) {
			// fmt.Println("total count", *response.Response.TotalCount, "len", len(securityGroups))
			break
		}
		request.Offset = tea.String(strconv.Itoa(len(securityGroups)))
	}
	return securityGroups, nil
}
