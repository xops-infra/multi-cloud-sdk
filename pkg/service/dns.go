package service

import (
	"fmt"
	"strings"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/xops-infra/multi-cloud-sdk/pkg/model"
)

type DnsService struct {
	Profiles     []model.ProfileConfig
	Aws, Tencent model.CloudIo
}

func NewDnsService(profiles []model.ProfileConfig, aws, tencent model.CloudIo) model.DnsContract {
	return &DnsService{
		Profiles: profiles,
		Aws:      aws,
		Tencent:  tencent,
	}
}

// DescribeDomainList
func (s *DnsService) DescribeDomainList(profile, region string, req model.DescribeDomainListRequest) (model.DescribeDomainListResponse, error) {
	for _, p := range s.Profiles {
		if p.Name == profile {
			switch p.Cloud {
			case model.AWS:
				return s.Aws.DescribeDomainList(profile, region, req)
			case model.TENCENT:
				return s.Tencent.DescribeDomainList(profile, region, req)
			default:
				return model.DescribeDomainListResponse{}, nil
			}
		}
	}
	return model.DescribeDomainListResponse{}, fmt.Errorf("profile not found")
}

// DescribeRecordList
func (s *DnsService) DescribeRecordList(profile, region string, req model.DescribeRecordListRequest) (model.DescribeRecordListResponse, error) {
	for _, p := range s.Profiles {
		if p.Name == profile {
			switch p.Cloud {
			case model.AWS:
				return s.Aws.DescribeRecordList(profile, region, req)
			case model.TENCENT:
				resp, err := s.Tencent.DescribeRecordList(profile, region, req)
				if err != nil {
					if strings.Contains(err.Error(), "ResourceNotFound") {
						return model.DescribeRecordListResponse{
							RecordCountInfo: &model.RecordCountInfo{
								Total: tea.Int64(0),
							},
						}, nil
					}
				}
				return resp, err
			default:
				return model.DescribeRecordListResponse{}, nil
			}
		}
	}
	return model.DescribeRecordListResponse{}, fmt.Errorf("profile not found")
}

// DescribeRecord
func (s *DnsService) DescribeRecord(profile, region string, req model.DescribeRecordRequest) (model.DescribeRecordResponse, error) {
	for _, p := range s.Profiles {
		if p.Name == profile {
			switch p.Cloud {
			case model.AWS:
				return s.Aws.DescribeRecord(profile, region, req)
			case model.TENCENT:
				resp, err := s.Tencent.DescribeRecord(profile, region, req)
				if err != nil {
					if strings.Contains(err.Error(), "ResourceNotFound") {
						return model.DescribeRecordResponse{}, nil
					}
				}
				return resp, err
			default:
				return model.DescribeRecordResponse{}, nil
			}
		}
	}
	return model.DescribeRecordResponse{}, fmt.Errorf("profile not found")
}

// CreateRecord
func (s *DnsService) CreateRecord(profile, region string, request model.CreateRecordRequest) (model.CreateRecordResponse, error) {
	for _, p := range s.Profiles {
		if p.Name == profile {
			switch p.Cloud {
			case model.AWS:
				return s.Aws.CreateRecord(profile, region, request)
			case model.TENCENT:
				return s.Tencent.CreateRecord(profile, region, request)
			default:
				return model.CreateRecordResponse{}, nil
			}
		}
	}
	return model.CreateRecordResponse{}, fmt.Errorf("profile not found")
}

func (s *DnsService) ModifyRecord(profile, region string, request model.ModifyRecordRequest) (model.ModifyRecordResponse, error) {
	for _, p := range s.Profiles {
		if p.Name == profile {
			switch p.Cloud {
			case model.AWS:
				return s.Aws.ModifyRecord(profile, region, request)
			case model.TENCENT:
				return s.Tencent.ModifyRecord(profile, region, request)
			default:
				return model.ModifyRecordResponse{}, nil
			}
		}
	}
	return model.ModifyRecordResponse{}, fmt.Errorf("profile not found")
}

func (s *DnsService) DeleteRecord(profile, region string, request model.DeleteRecordRequest) (model.CommonDnsResponse, error) {
	if request.Domain == nil || request.SubDomain == nil || request.RecordType == nil {
		return model.CommonDnsResponse{}, fmt.Errorf("domain, subDomain, recordType is required")
	}
	for _, p := range s.Profiles {
		if p.Name == profile {
			switch p.Cloud {
			case model.AWS:
				return s.Aws.DeleteRecord(profile, region, request)
			case model.TENCENT:
				return s.Tencent.DeleteRecord(profile, region, request)
			default:
				return model.CommonDnsResponse{}, fmt.Errorf("not support cloud %s", p.Cloud)
			}
		}
	}
	return model.CommonDnsResponse{}, fmt.Errorf("profile not found")
}
