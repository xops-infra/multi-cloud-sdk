package io

import (
	"fmt"
	"strings"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/spf13/cast"
	"github.com/xops-infra/multi-cloud-sdk/pkg/model"
)

// DescribeDomainList
func (c *awsClient) DescribeDomainList(profile, region string, input model.DescribeDomainListRequest) (model.DescribeDomainListResponse, error) {
	client, err := c.io.GetAwsRoute53Client(profile, region)
	if err != nil {
		return model.DescribeDomainListResponse{}, err
	}

	params := &route53.ListHostedZonesInput{}
	var domains []model.Domain

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
			domains = append(domains, model.Domain{
				DomainId: domain.Id,
				Name:     tea.String(strings.TrimSuffix(*domain.Name, ".")),
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

func (c *awsClient) DescribeRecordListWithPages(profile, region string, input model.DescribeRecordListWithPageRequest) (model.ListRecordsPageResponse, error) {
	if input.Domain == nil {
		return model.ListRecordsPageResponse{}, fmt.Errorf("domain,region is required")
	}
	client, err := c.io.GetAwsRoute53Client(profile, region)
	if err != nil {
		return model.ListRecordsPageResponse{}, err
	}

	params := &route53.ListResourceRecordSetsInput{
		HostedZoneId: input.Domain,
		MaxItems:     tea.String("100"),
	}
	domain, err := c.getHostedZoneIdByDomain(profile, region, input.Domain)
	if err != nil {
		return model.ListRecordsPageResponse{}, err
	}
	params.HostedZoneId = domain.DomainId

	if input.Limit != nil {
		params.MaxItems = tea.String(cast.ToString(input.Limit))
	}
	var pageNum int64
	var resp model.ListRecordsPageResponse

	err = client.ListResourceRecordSetsPages(params,
		func(page *route53.ListResourceRecordSetsOutput, lastPage bool) bool {
			// fmt.Printf("---%d---\n %s", pageNum, tea.Prettify(page))
			if input.Page == nil {
				input.Page = tea.Int64(1)
			}
			if pageNum == *input.Page-1 {
				for _, record := range page.ResourceRecordSets {
					var recordValue *string
					if record.ResourceRecords == nil {
						recordValue = record.AliasTarget.DNSName
					} else {
						// TODO: 有多个值的情况
						for _, value := range record.ResourceRecords {
							recordValue = value.Value
						}
					}
					// 解决httpDecode问题，比如 * -> \\052 @ -> \\100 # -> \\043
					record.Name = aws.String(strings.ReplaceAll(aws.StringValue(record.Name), "\\052", "*"))
					record.Name = aws.String(strings.ReplaceAll(*record.Name, "\\100", "@"))
					record.Name = aws.String(strings.ReplaceAll(*record.Name, "\\043", "#"))
					subDomain := strings.TrimSuffix(*record.Name, fmt.Sprintf("%s.", *domain.Name))
					resp.RecordList = append(resp.RecordList, model.Record{
						SubDomain:  aws.String(strings.TrimSuffix(subDomain, ".")),
						TTL:        tea.Uint64(cast.ToUint64(record.TTL)),
						Weight:     tea.Uint64(cast.ToUint64(record.Weight)),
						RecordType: record.Type,
						Value:      recordValue,
						Status:     record.SetIdentifier,
						RecordId:   record.Name,
					})
				}
				if len(resp.RecordList) == cast.ToInt(params.MaxItems) {
					resp.NextPage = tea.Int64(pageNum + 2)
					resp.PrePage = tea.Int64(pageNum)
					if pageNum == 0 {
						resp.PrePage = nil
					}
				}
				return false
			} else {
				pageNum++
				return true
			}
		})
	if err != nil {
		return resp, err
	}
	return resp, nil
}

// DescribeRecordList
func (c *awsClient) DescribeRecordList(profile, region string, input model.DescribeRecordListRequest) (model.DescribeRecordListResponse, error) {

	client, err := c.io.GetAwsRoute53Client(profile, region)
	if err != nil {
		return model.DescribeRecordListResponse{}, err
	}

	param := &route53.ListResourceRecordSetsInput{
		HostedZoneId: input.Domain,
	}
	domain, err := c.getHostedZoneIdByDomain(profile, region, input.Domain)
	if err != nil {
		return model.DescribeRecordListResponse{}, err
	}
	param.HostedZoneId = domain.DomainId

	var records []model.Record
	resp, err := client.ListResourceRecordSets(param)
	if err != nil {
		return model.DescribeRecordListResponse{}, err
	}
	for {
		for _, record := range resp.ResourceRecordSets {
			var values string
			for _, value := range record.ResourceRecords {
				values = *value.Value
			}
			if input.Keyword != nil && *input.Keyword != "" {
				if !strings.Contains(*record.Name, *input.Keyword) {
					continue
				}
			}
			// 解决httpDecode问题，比如 \\052
			record.Name = tea.String(strings.ReplaceAll(*record.Name, "\\052", "*"))
			subDomain := strings.TrimSuffix(*record.Name, fmt.Sprintf("%s.", *domain.Name))
			records = append(records, model.Record{
				SubDomain:  tea.String(strings.TrimSuffix(subDomain, ".")),
				TTL:        tea.Uint64(cast.ToUint64(record.TTL)),
				Weight:     tea.Uint64(cast.ToUint64(record.Weight)),
				RecordType: record.Type,
				Value:      tea.String(values),
				Status:     record.SetIdentifier,
				RecordId:   record.Name,
			})
		}

		if resp.IsTruncated == nil || !*resp.IsTruncated {
			break
		}
		param.StartRecordName = resp.NextRecordName
		param.StartRecordType = resp.NextRecordType
		param.StartRecordIdentifier = resp.NextRecordIdentifier
		resp, err = client.ListResourceRecordSets(param)
		if err != nil {
			return model.DescribeRecordListResponse{}, err
		}
	}
	return model.DescribeRecordListResponse{
		Total:      cast.ToInt64(len(records)),
		RecordList: records,
	}, nil

}

// DescribeRecord 完全匹配
func (c *awsClient) DescribeRecord(profile, region string, input model.DescribeRecordRequest) (model.Record, error) {
	resp, err := c.DescribeRecordList(profile, region, model.DescribeRecordListRequest{
		Domain:  input.Domain,
		Keyword: input.SubDomain,
	})
	if err != nil {
		return model.Record{}, err
	}
	for _, record := range resp.RecordList {
		if *record.SubDomain == *input.SubDomain {
			if input.RecordType != nil && *input.RecordType != "" && *input.RecordType != *record.RecordType {
				continue
			}
			return record, nil
		}
	}
	return model.Record{}, fmt.Errorf("record not found")
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
	domain, err := c.getHostedZoneIdByDomain(profile, region, input.Domain)
	if err != nil {
		return model.CreateRecordResponse{}, err
	}
	param := &route53.ChangeResourceRecordSetsInput{
		HostedZoneId: domain.DomainId,
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
			Comment: tea.String(fmt.Sprintf("%s, created by multi-cloud-sdk", tea.StringValue(input.Info))),
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

// getHostedZoneIdByDomain domain 为域名或者hostedzoneId
func (c *awsClient) getHostedZoneIdByDomain(profile, region string, domain *string) (*model.Domain, error) {
	if domain == nil {
		return nil, fmt.Errorf("domain is required")
	}
	resp, err := c.DescribeDomainList(profile, region, model.DescribeDomainListRequest{})
	if err != nil {
		return nil, err
	}
	for _, _domain := range resp.DomainList {
		if strings.HasPrefix(*domain, "/hostedzone/") {
			if *_domain.DomainId == *domain {
				return &_domain, nil
			}
		} else {
			if *_domain.Name == *domain {
				return &_domain, nil
			}
		}
	}
	return nil, fmt.Errorf("domain not found")
}

// ModifyRecord
// ignoreType 腾讯云修改需要一起提供记录类型，aws不需要，所以不处理
func (c *awsClient) ModifyRecord(profile, region string, ignoreType bool, input model.ModifyRecordRequest) error {
	cloudClient, err := c.io.GetAwsRoute53Client(profile, region)
	if err != nil {
		return err
	}
	resp, err := c.getHostedZoneIdByDomain(profile, region, input.Domain)
	if err != nil {
		return err
	}
	var ttl int64 = 300
	if input.TTL != nil {
		ttl = cast.ToInt64(input.TTL)
	}
	param := &route53.ChangeResourceRecordSetsInput{
		HostedZoneId: resp.DomainId,
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
	_, err = cloudClient.ChangeResourceRecordSets(param)
	if err != nil {
		return err
	}
	return nil

}

// DeleteDns
func (c *awsClient) DeleteRecord(profile, region string, input model.DeleteRecordRequest) (model.CommonDnsResponse, error) {
	if input.Domain == nil || input.SubDomain == nil || input.RecordType == nil {
		return model.CommonDnsResponse{}, fmt.Errorf("domain, subDomain, recordType is required")
	}
	client, err := c.io.GetAwsRoute53Client(profile, region)
	if err != nil {
		return model.CommonDnsResponse{}, err
	}
	hostedZoneId, err := c.getHostedZoneIdByDomain(profile, region, input.Domain)
	if err != nil {
		return model.CommonDnsResponse{}, err
	}
	record, err := c.DescribeRecord(profile, region, model.DescribeRecordRequest{
		Domain:     input.Domain,
		SubDomain:  input.SubDomain,
		RecordType: input.RecordType,
	})
	if err != nil {
		return model.CommonDnsResponse{}, err
	}
	param := &route53.ChangeResourceRecordSetsInput{
		HostedZoneId: hostedZoneId.DomainId,
		ChangeBatch: &route53.ChangeBatch{
			Changes: []*route53.Change{
				{
					Action: aws.String("DELETE"),
					ResourceRecordSet: &route53.ResourceRecordSet{
						Name: tea.String(fmt.Sprintf("%s.%s.", *input.SubDomain, *input.Domain)),
						Type: record.RecordType,
						TTL:  tea.Int64(cast.ToInt64(record.TTL)),
						ResourceRecords: []*route53.ResourceRecord{
							{
								Value: record.Value,
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
