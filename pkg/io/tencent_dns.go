package io

import (
	"fmt"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/spf13/cast"
	dnspod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323"
	"github.com/xops-infra/multi-cloud-sdk/pkg/model"
)

// DescribeDomainList
func (c *tencentClient) DescribeDomainList(profile, region string, input model.DescribeDomainListRequest) (model.DescribeDomainListResponse, error) {
	client, err := c.io.GetTencentDnsPodClient(profile, region)
	if err != nil {
		return model.DescribeDomainListResponse{}, err
	}
	request := dnspod.NewDescribeDomainListRequest()
	request.Type = tea.String("ALL")
	request.Keyword = input.DomainKeyword

	response, err := client.DescribeDomainList(request)
	if err != nil {
		return model.DescribeDomainListResponse{}, err
	}
	var domains []model.Domain
	for _, domain := range response.Response.DomainList {
		domains = append(domains, model.Domain{
			DomainId: tea.String(cast.ToString(domain.DomainId)),
			Name:     domain.Name,
			Meta:     domain,
		})

	}
	return model.DescribeDomainListResponse{
		RequestId:  response.Response.RequestId,
		DomainList: domains,
		DomainCountInfo: &model.DomainCountInfo{
			Total: tea.Int64(cast.ToInt64(response.Response.DomainCountInfo.AllTotal)),
		},
	}, nil
}

// DescribeRecordList
func (c *tencentClient) DescribeRecordList(profile, region string, input model.DescribeRecordListRequest) (model.DescribeRecordListResponse, error) {
	client, err := c.io.GetTencentDnsPodClient(profile, region)
	if err != nil {
		return model.DescribeRecordListResponse{}, err
	}
	request := dnspod.NewDescribeRecordListRequest()
	request.Domain = input.Domain
	// request.Subdomain = input.Keyword
	request.RecordType = input.RecordType
	request.Keyword = input.Keyword
	request.Limit = tea.Uint64(100)
	if input.Limit != nil {
		request.Limit = tea.Uint64(cast.ToUint64(*input.Limit))
	}

	if input.NextMarker != nil {
		request.Offset = tea.Uint64((cast.ToUint64(model.DecodeTencentNextMaker(*input.NextMarker)) - 1) * uint64(*request.Limit))
	}

	response, err := client.DescribeRecordList(request)
	if err != nil {
		return model.DescribeRecordListResponse{}, err
	}
	var records []model.Record
	for _, record := range response.Response.RecordList {
		records = append(records, model.Record{
			RecordId:   tea.String(cast.ToString(record.RecordId)),
			SubDomain:  record.Name,
			RecordType: record.Type,
			Value:      record.Value,
			Status:     record.Status,
			UpdatedOn:  record.UpdatedOn,
			TTL:        record.TTL,
			RecordLine: record.Line,
			Remark:     record.Remark,
			Weight:     record.Weight,
		})
	}
	var nextMaker *string
	if input.NextMarker != nil {
		// 如果偏移量小于总数则说明还有，nextMarker为偏移量加1
		// fmt.Println(*response.Response.RecordCountInfo.TotalCount, *request.Offset)
		if *response.Response.RecordCountInfo.TotalCount > (*request.Offset + uint64(*input.Limit)) {
			nextMaker = model.ToTencentNextMaker(cast.ToString(cast.ToInt(model.DecodeTencentNextMaker(*input.NextMarker)) + 1))
		} else {
			nextMaker = nil
		}
	} else if *response.Response.RecordCountInfo.TotalCount > uint64(*request.Limit) {
		nextMaker = model.ToTencentNextMaker("2")
	}

	return model.DescribeRecordListResponse{
		RecordList: records,
		NextMarker: nextMaker,
	}, nil
}

// DescribeRecord
func (c *tencentClient) DescribeRecord(profile, region string, input model.DescribeRecordRequest) (model.DescribeRecordResponse, error) {
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
func (c *tencentClient) CreateRecord(profile, region string, input model.CreateRecordRequest) (model.CreateRecordResponse, error) {
	client, err := c.io.GetTencentDnsPodClient(profile, region)
	if err != nil {
		return model.CreateRecordResponse{}, err
	}
	// 实例化一个请求对象,每个接口都会对应一个request对象
	request := dnspod.NewCreateRecordRequest()
	request.Domain = input.Domain
	request.SubDomain = input.SubDomain
	request.RecordType = input.RecordType
	request.Value = input.Value
	request.RecordLine = tea.String("默认")
	request.Remark = input.Info
	if input.TTL != nil {
		request.TTL = input.TTL
	} else {
		request.TTL = tea.Uint64(600)
	}
	// 返回的resp是一个CreatePrivateZoneRecordResponse的实例，与请求对象对应
	response, err := client.CreateRecord(request)
	if err != nil {
		return model.CreateRecordResponse{}, err
	}
	return model.CreateRecordResponse{
		RecordId: tea.String(cast.ToString(response.Response.RecordId)),
		Meta:     response.Response,
	}, nil
}

// ModifyRecord
// ignoreType 是否开启忽略 recordType,
// true 注意这里会删除所有相同 subDomain 的记录，然后创建新的记录
// false 如果 recordType 不同，会报没找到记录
func (c *tencentClient) ModifyRecord(profile, region string, ignoreType bool, input model.ModifyRecordRequest) (model.ModifyRecordResponse, error) {
	client, err := c.io.GetTencentDnsPodClient(profile, region)
	if err != nil {
		return model.ModifyRecordResponse{}, err
	}
	if ignoreType {
		resp, err := c.DescribeRecordList(profile, region, model.DescribeRecordListRequest{
			Domain:  input.Domain,
			Keyword: input.SubDomain,
		})
		if err != nil {
			return model.ModifyRecordResponse{}, err
		}
		var delDomain []map[string]interface{}
		for _, record := range resp.RecordList {
			if *record.SubDomain == *input.SubDomain {
				_, err := c.DeleteRecord(profile, region, model.DeleteRecordRequest{
					Domain:     input.Domain,
					SubDomain:  input.SubDomain,
					RecordType: record.RecordType,
				})
				if err != nil {
					return model.ModifyRecordResponse{}, fmt.Errorf("delete record error: %v", err)
				}
				delDomain = append(delDomain, map[string]interface{}{
					"recordId":   record.RecordId,
					"recordType": record.RecordType,
					"subDomain":  record.SubDomain,
					"ttl":        record.TTL,
					"value":      record.Value,
				})
			}
		}

		if delDomain == nil {
			return model.ModifyRecordResponse{}, fmt.Errorf("record not found")
		}

		createInput := model.CreateRecordRequest{
			Domain:     input.Domain,
			SubDomain:  input.SubDomain,
			RecordType: input.RecordType,
			Value:      input.Value,
			TTL:        tea.Uint64(600),
			Info:       input.Info,
		}
		if input.TTL != nil {
			createInput.TTL = input.TTL
		}
		createResp, err := c.CreateRecord(profile, region, createInput)
		if err != nil {
			return model.ModifyRecordResponse{}, fmt.Errorf("create record error: %v", err)
		}
		return model.ModifyRecordResponse{
			RecordId: createResp.RecordId,
			Meta: map[string]interface{}{
				"info":      "enable ignoreType",
				"delDomain": delDomain,
			},
		}, nil
	} else {
		recordId, err := c.getRecordIdBySubDomain(profile, region, *input.SubDomain, *input.Domain, *input.RecordType)
		if err != nil {
			return model.ModifyRecordResponse{}, err
		}

		request := dnspod.NewModifyRecordRequest()
		request.RecordId = recordId
		request.Domain = input.Domain
		request.SubDomain = input.SubDomain
		request.RecordType = input.RecordType
		request.Value = input.Value
		request.RecordLine = tea.String("默认")
		request.TTL = input.TTL
		request.Weight = input.Weight

		if input.Status != nil {
			if *input.Status {
				request.Status = tea.String("ENABLE")
			} else {
				request.Status = tea.String("DISABLE")
			}
		}

		response, err := client.ModifyRecord(request)
		if err != nil {
			return model.ModifyRecordResponse{}, err
		}
		return model.ModifyRecordResponse{
			RecordId: tea.String(cast.ToString(response.Response.RecordId)),
			Meta:     response.Response,
		}, nil
	}
}

// getRecordIdBySubDomain
func (c *tencentClient) getRecordIdBySubDomain(profile, region, subDomain, domain, recordType string) (*uint64, error) {
	resp, err := c.DescribeRecordList(profile, region, model.DescribeRecordListRequest{
		Domain:  &domain,
		Keyword: &subDomain,
	})
	if err != nil {
		return nil, err
	}
	for _, record := range resp.RecordList {
		if *record.SubDomain == subDomain && *record.RecordType == recordType {
			return tea.Uint64(cast.ToUint64(record.RecordId)), nil
		}
	}
	return nil, fmt.Errorf("record not found")
}

// DeleteRecord
func (c *tencentClient) DeleteRecord(profile, region string, input model.DeleteRecordRequest) (model.CommonDnsResponse, error) {
	client, err := c.io.GetTencentDnsPodClient(profile, region)
	if err != nil {
		return model.CommonDnsResponse{}, err
	}
	record_id, err := c.getRecordIdBySubDomain(profile, region, *input.SubDomain, *input.Domain, *input.RecordType)
	if err != nil {
		return model.CommonDnsResponse{}, err
	}

	request := dnspod.NewDeleteRecordRequest()
	request.RecordId = record_id
	request.Domain = input.Domain

	resp, err := client.DeleteRecord(request)
	if err != nil {
		return model.CommonDnsResponse{}, err
	}
	return model.CommonDnsResponse{
		Meta: resp.Response,
	}, nil
}
