package io

import (
	"fmt"
	"strings"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/spf13/cast"

	"github.com/xops-infra/multi-cloud-sdk/pkg/model"
)

type awsClient struct {
	io model.ClientIo
}

func NewAwsClient(io model.ClientIo) model.CloudIo {
	return &awsClient{
		io: io,
	}
}

func (c *awsClient) QueryInstances(profile, region string) ([]*model.Instance, error) {
	svc, err := c.io.GetAwsEc2Client(profile, region)
	if err != nil {
		return nil, err
	}
	input := &ec2.DescribeInstancesInput{
		MaxResults: aws.Int64(10),
	}
	var instances []*model.Instance
	for {
		out, err := svc.DescribeInstances(input)
		if err != nil {
			return nil, err
		}

		for _, reservation := range out.Reservations {
			for _, instance := range reservation.Instances {
				tags := model.AwsTagsToModelTags(instance.Tags)
				instances = append(instances, &model.Instance{
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
		if out.NextToken == nil {
			break
		}
		input.NextToken = out.NextToken
	}
	return instances, nil
}

func (c *awsClient) DescribeInstances(profile, region string, instanceIds []*string) ([]*model.Instance, error) {
	svc, err := c.io.GetAwsEc2Client(profile, region)
	if err != nil {
		return nil, err
	}
	input := &ec2.DescribeInstancesInput{
		InstanceIds: instanceIds,
	}
	out, err := svc.DescribeInstances(input)
	if err != nil {
		return nil, err
	}
	var instances []*model.Instance
	for _, reservation := range out.Reservations {
		for _, instance := range reservation.Instances {
			tags := model.AwsTagsToModelTags(instance.Tags)
			instances = append(instances, &model.Instance{
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
	return instances, nil
}

// QueryVpcs
func (c *awsClient) QueryVPC(profile, region string, input model.CommonQueryInput) ([]*model.VPC, error) {
	if input.CloudProvider != "" && input.CloudProvider != model.AWS {
		return nil, nil
	}
	if input.Account != "" && input.Account != profile {
		return nil, nil
	}
	if input.Region != "" && input.Region != region {
		return nil, nil
	}
	svc, err := c.io.GetAwsEc2Client(profile, region)
	if err != nil {
		return nil, err
	}
	var vpcs []*model.VPC
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
			vpcs = append(vpcs, &model.VPC{
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
func (c *awsClient) QuerySubnet(profile, region string, input model.CommonQueryInput) ([]*model.Subnet, error) {
	if !input.Filter(profile, region) {
		return nil, nil
	}
	svc, err := c.io.GetAwsEc2Client(profile, region)
	if err != nil {
		return nil, err
	}
	var subnets []*model.Subnet
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
			subnets = append(subnets, &model.Subnet{
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
func (c *awsClient) QueryEIP(profile, region string, input model.CommonQueryInput) ([]*model.EIP, error) {
	if !input.Filter(profile, region) {
		return nil, nil
	}
	svc, err := c.io.GetAwsEc2Client(profile, region)
	if err != nil {
		return nil, err
	}
	var eips []*model.EIP
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
		eips = append(eips, &model.EIP{
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
func (c *awsClient) QueryNAT(profile, region string, input model.CommonQueryInput) ([]*model.NAT, error) {
	if !input.Filter(profile, region) {
		return nil, nil
	}
	svc, err := c.io.GetAwsEc2Client(profile, region)
	if err != nil {
		return nil, err
	}
	var nats []*model.NAT
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
		nats = append(nats, &model.NAT{
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

// CommonOCR
func (c *awsClient) CommonOCR(profile, region string, input model.OcrRequest) (model.OcrResponse, error) {
	return model.OcrResponse{}, nil
}

// CreatePicture
func (c *awsClient) CreatePicture(profile, region string, input model.CreatePictureRequest) (model.CreatePictureResponse, error) {
	return model.CreatePictureResponse{}, nil
}

// GetPictureByName
func (c *awsClient) GetPictureByName(profile, region string, input model.CommonPictureRequest) (model.GetPictureByNameResponse, error) {
	return model.GetPictureByNameResponse{}, nil
}

// QueryPicture
func (c *awsClient) QueryPicture(profile, region string, input model.QueryPictureRequest) (model.QueryPictureResponse, error) {
	return model.QueryPictureResponse{}, nil
}

// DeletePicture
func (c *awsClient) DeletePicture(profile, region string, input model.CommonPictureRequest) (model.CommonPictureResponse, error) {
	return model.CommonPictureResponse{}, nil
}

// UpdatePicture
func (c *awsClient) UpdatePicture(profile, region string, input model.UpdatePictureRequest) (model.CommonPictureResponse, error) {
	return model.CommonPictureResponse{}, nil
}

// SearchPicture
func (c *awsClient) SearchPicture(profile, region string, input model.SearchPictureRequest) (model.SearchPictureResponse, error) {
	return model.SearchPictureResponse{}, nil
}

// DescribeDomainList
func (c *awsClient) DescribeDomainList(profile, region string, input model.DescribeDomainListRequest) (model.DescribeDomainListResponse, error) {
	client, err := c.io.GetAwsRoute53Client(profile, region)
	if err != nil {
		return model.DescribeDomainListResponse{}, err
	}

	params := &route53.ListHostedZonesInput{}
	var domains []*model.Domain

	for {
		resp, err := client.ListHostedZones(params)
		if err != nil {
			return model.DescribeDomainListResponse{}, err
		}
		for _, domain := range resp.HostedZones {
			if input.DomainKeyword != nil && *input.DomainKeyword != "" {
				if !strings.Contains(*domain.Name, *input.DomainKeyword) {
					continue
				}
			}
			domains = append(domains, &model.Domain{
				DomainId: domain.Id,
				Name:     domain.Name,
				Meta:     domain,
			})
		}
		if resp.IsTruncated == nil || !*resp.IsTruncated {
			break
		}
		params.Marker = resp.NextMarker
	}

	return model.DescribeDomainListResponse{
		DomainCountInfo: &model.DomainCountInfo{
			Total: tea.Int64(cast.ToInt64(len(domains))),
		},
		DomainList: domains,
	}, nil
}

// DescribeRecordList
func (c *awsClient) DescribeRecordList(profile, region string, input model.DescribeRecordListRequest) (model.DescribeRecordListResponse, error) {
	client, err := c.io.GetAwsRoute53Client(profile, region)
	if err != nil {
		return model.DescribeRecordListResponse{}, err
	}

	hostedZoneId, err := c.getHostedZoneIdByDomain(profile, region, input.Domain)
	if err != nil {
		return model.DescribeRecordListResponse{}, err
	}

	param := &route53.ListResourceRecordSetsInput{
		HostedZoneId: hostedZoneId,
	}

	var records []*model.Record
	for {
		resp, err := client.ListResourceRecordSets(param)
		if err != nil {
			return model.DescribeRecordListResponse{}, err
		}
		for _, record := range resp.ResourceRecordSets {
			var values string
			for _, value := range record.ResourceRecords {
				values = *value.Value
			}
			if input.RecordType != nil && *input.RecordType != "" && *input.RecordType != *record.Type {
				continue
			}
			if input.Keyword != nil && *input.Keyword != "" {
				if !strings.Contains(*record.Name, *input.Keyword) && !strings.Contains(values, *input.Keyword) {
					continue
				}
			}
			subDomain := strings.TrimSuffix(*record.Name, fmt.Sprintf(".%s.", *input.Domain))
			records = append(records, &model.Record{
				SubDomain:  tea.String(subDomain),
				TTL:        tea.Uint64(cast.ToUint64(record.TTL)),
				Weight:     tea.Uint64(cast.ToUint64(record.Weight)),
				RecordType: record.Type,
				Value:      tea.String(values),
				Status:     record.SetIdentifier,
				Meta:       record,
			})
		}
		if resp.IsTruncated == nil || !*resp.IsTruncated {
			break
		}
		param.StartRecordName = resp.NextRecordName
		param.StartRecordType = resp.NextRecordType
	}
	return model.DescribeRecordListResponse{
		RecordList: records,
		RecordCountInfo: &model.RecordCountInfo{
			Total: tea.Int64(cast.ToInt64(len(records))),
		},
	}, nil

}

// DescribeRecord
func (c *awsClient) DescribeRecord(profile, region string, input model.DescribeRecordRequest) (model.DescribeRecordResponse, error) {
	resp, err := c.DescribeRecordList(profile, region, model.DescribeRecordListRequest{
		Domain:  input.Domain,
		Keyword: input.SubDomain,
	})
	if err != nil {
		return model.DescribeRecordResponse{}, err
	}
	for _, record := range resp.RecordList {
		if *record.SubDomain == *input.SubDomain {
			if input.RecordType != nil && *input.RecordType != "" && *input.RecordType != *record.RecordType {
				continue
			}
			return model.DescribeRecordResponse{
				Record: record,
			}, nil
		}
	}
	return model.DescribeRecordResponse{}, fmt.Errorf("record not found")
}

// CreateRecord
func (c *awsClient) CreateRecord(profile, region string, input model.CreateRecordRequest) (model.CreateRecordResponse, error) {
	client, err := c.io.GetAwsRoute53Client(profile, region)
	if err != nil {
		return model.CreateRecordResponse{}, err
	}
	var ttl int64 = 300
	if input.TTL != nil {
		ttl = cast.ToInt64(input.TTL)
	}
	hostedZoneId, err := c.getHostedZoneIdByDomain(profile, region, input.Domain)
	if err != nil {
		return model.CreateRecordResponse{}, err
	}
	param := &route53.ChangeResourceRecordSetsInput{
		HostedZoneId: hostedZoneId,
		ChangeBatch: &route53.ChangeBatch{
			Changes: []*route53.Change{
				{
					Action: aws.String("CREATE"),
					ResourceRecordSet: &route53.ResourceRecordSet{
						Name: tea.String(fmt.Sprintf("%s.%s.", *input.SubDomain, *input.Domain)),
						Type: input.RecordType,
						TTL:  tea.Int64(ttl),
						ResourceRecords: []*route53.ResourceRecord{
							{
								Value: input.Value,
							},
						},
					},
				},
			},
			Comment: tea.String(fmt.Sprintf("%s, created by multi-cloud-sdk", *input.Info)),
		},
	}
	resp, err := client.ChangeResourceRecordSets(param)
	if err != nil {
		return model.CreateRecordResponse{}, err
	}
	return model.CreateRecordResponse{
		RecordId: resp.ChangeInfo.Id,
		Meta:     resp.ChangeInfo,
	}, nil
}

// getHostedZoneIdByDomain
func (c *awsClient) getHostedZoneIdByDomain(profile, region string, domain *string) (*string, error) {
	resp, err := c.DescribeDomainList(profile, region, model.DescribeDomainListRequest{
		DomainKeyword: domain,
	})
	if err != nil {
		return nil, err
	}
	domain = tea.String(strings.TrimSuffix(*domain, "."))

	for _, _domain := range resp.DomainList {
		if *_domain.Name == *domain+"." {
			return _domain.DomainId, nil
		}
	}
	return nil, fmt.Errorf("domain not found")
}

// ModifyRecord
// ignoreType 腾讯云修改需要一起提供记录类型，aws不需要，所以不处理
func (c *awsClient) ModifyRecord(profile, region string, ignoreType bool, input model.ModifyRecordRequest) (model.ModifyRecordResponse, error) {
	cloudClient, err := c.io.GetAwsRoute53Client(profile, region)
	if err != nil {
		return model.ModifyRecordResponse{}, err
	}
	hostedZoneId, err := c.getHostedZoneIdByDomain(profile, region, input.Domain)
	if err != nil {
		return model.ModifyRecordResponse{}, err
	}
	var ttl int64 = 300
	if input.TTL != nil {
		ttl = cast.ToInt64(input.TTL)
	}
	param := &route53.ChangeResourceRecordSetsInput{
		HostedZoneId: hostedZoneId,
		ChangeBatch: &route53.ChangeBatch{
			Changes: []*route53.Change{
				{
					Action: aws.String("UPSERT"),
					ResourceRecordSet: &route53.ResourceRecordSet{
						Name: tea.String(fmt.Sprintf("%s.%s.", *input.SubDomain, *input.Domain)),
						Type: input.RecordType,
						TTL:  tea.Int64(ttl),
						ResourceRecords: []*route53.ResourceRecord{
							{
								Value: input.Value,
							},
						},
					},
				},
			},
		},
	}
	resp, err := cloudClient.ChangeResourceRecordSets(param)
	if err != nil {
		return model.ModifyRecordResponse{}, err
	}
	return model.ModifyRecordResponse{
		RecordId: resp.ChangeInfo.Id,
		Meta:     resp.ChangeInfo,
	}, nil

}

// DeleteDns
func (c *awsClient) DeleteRecord(profile, region string, input model.DeleteRecordRequest) (model.CommonDnsResponse, error) {
	client, err := c.io.GetAwsRoute53Client(profile, region)
	if err != nil {
		return model.CommonDnsResponse{}, err
	}
	hostedZoneId, err := c.getHostedZoneIdByDomain(profile, region, input.Domain)
	if err != nil {
		return model.CommonDnsResponse{}, err
	}
	recordResp, err := c.DescribeRecord(profile, region, model.DescribeRecordRequest{
		Domain:     input.Domain,
		SubDomain:  input.SubDomain,
		RecordType: input.RecordType,
	})
	if err != nil {
		return model.CommonDnsResponse{}, err
	}
	param := &route53.ChangeResourceRecordSetsInput{
		HostedZoneId: hostedZoneId,
		ChangeBatch: &route53.ChangeBatch{
			Changes: []*route53.Change{
				{
					Action: aws.String("DELETE"),
					ResourceRecordSet: &route53.ResourceRecordSet{
						Name: tea.String(fmt.Sprintf("%s.%s.", *input.SubDomain, *input.Domain)),
						Type: recordResp.Record.RecordType,
						TTL:  tea.Int64(cast.ToInt64(recordResp.Record.TTL)),
						ResourceRecords: []*route53.ResourceRecord{
							{
								Value: recordResp.Record.Value,
							},
						},
					},
				},
			},
			Comment: tea.String("deleted by multi-cloud-sdk"),
		},
	}
	resp, err := client.ChangeResourceRecordSets(param)
	if err != nil {
		return model.CommonDnsResponse{}, err
	}
	return model.CommonDnsResponse{
		Meta: resp.ChangeInfo,
	}, nil
}
