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

type awsClient struct {
	io model.ClientIo
}

func NewAwsClient(io model.ClientIo) model.CloudIO {
	return &awsClient{
		io: io,
	}
}

func (c *awsClient) CreateTags(profile, region string, input model.CreateTagsInput) error {
	return fmt.Errorf("not support for aws")
}

func (c *awsClient) AddTagsToResource(profile, region string, input model.AddTagsInput) error {
	return fmt.Errorf("not support for aws")
}

func (c *awsClient) RemoveTagsFromResource(profile, region string, input model.RemoveTagsInput) error {
	return fmt.Errorf("not support for aws")
}

func (c *awsClient) ModifyTagsForResource(profile, region string, input model.ModifyTagsInput) error {
	return fmt.Errorf("not support for aws")
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
	if input.Limit != nil {
		param.MaxItems = tea.String(cast.ToString(input.Limit))
	}
	if input.NextMarker != nil {
		param.StartRecordName, param.StartRecordType = model.DecodeAwsNextMaker(*input.NextMarker)
	}
	// fmt.Println(tea.Prettify(param))

	var records []model.Record
	resp, err := client.ListResourceRecordSets(param)
	if err != nil {
		return model.DescribeRecordListResponse{}, err
	}
	var nextMarker string
pageLoop:
	for {
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
			// 解决httpDecode问题，比如 \\052
			record.Name = tea.String(strings.ReplaceAll(*record.Name, "\\052", "*"))
			subDomain := strings.TrimSuffix(*record.Name, fmt.Sprintf("%s.", *input.Domain))
			records = append(records, model.Record{
				SubDomain:  tea.String(strings.TrimSuffix(subDomain, ".")),
				TTL:        tea.Uint64(cast.ToUint64(record.TTL)),
				Weight:     tea.Uint64(cast.ToUint64(record.Weight)),
				RecordType: record.Type,
				Value:      tea.String(values),
				Status:     record.SetIdentifier,
				RecordId:   record.Name,
			})
			nextMarker = model.ToAwsNextMaker(record.Name, record.Type)
			if len(records) == int(*input.Limit+1) {
				break pageLoop
			}
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
		RecordList: records[:len(records)-1],
		NextMarker: tea.String(nextMarker),
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
