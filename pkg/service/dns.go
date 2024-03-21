package service

import (
	"fmt"
	"strings"

	"github.com/xops-infra/multi-cloud-sdk/pkg/model"
)

type DnsService struct {
	Profiles     map[string]model.ProfileConfig
	Aws, Tencent model.CloudIO
}

func NewDnsService(profiles []model.ProfileConfig, aws, tencent model.CloudIO) model.DnsContract {
	_profiles := make(map[string]model.ProfileConfig)
	for _, p := range profiles {
		_profiles[p.Name] = p

	}
	return &DnsService{
		Profiles: _profiles,
		Aws:      aws,
		Tencent:  tencent,
	}
}

// PrivateDomainList
func (s *DnsService) PrivateDomainList(profile string, req model.DescribeDomainListRequest) (model.DescribePrivateDomainListResponse, error) {
	if p, ok := s.Profiles[profile]; ok {
		switch p.Cloud {
		case model.AWS:
			return s.Aws.DescribePrivateDomainList(profile, req)
		case model.TENCENT:
			return s.Tencent.DescribePrivateDomainList(profile, req)
		default:
			return model.DescribePrivateDomainListResponse{}, nil
		}
	}
	return model.DescribePrivateDomainListResponse{}, fmt.Errorf("profile not found")
}

// PrivateRecordList
func (s *DnsService) PrivateRecordList(profile string, req model.DescribePrivateRecordListRequest) (model.DescribePrivateRecordListResponse, error) {
	if p, ok := s.Profiles[profile]; ok {
		switch p.Cloud {
		case model.AWS:
			return s.Aws.DescribePrivateRecordList(profile, req)
		case model.TENCENT:
			return s.Tencent.DescribePrivateRecordList(profile, req)
		default:
			return model.DescribePrivateRecordListResponse{}, nil
		}
	}
	return model.DescribePrivateRecordListResponse{}, fmt.Errorf("profile not found")
}

// PrivateRecordListWithPages
func (s *DnsService) PrivateRecordListWithPages(profile string, req model.DescribePrivateDnsRecordListWithPageRequest) (model.ListRecordsPageResponse, error) {
	if p, ok := s.Profiles[profile]; ok {
		switch p.Cloud {
		case model.AWS:
			return s.Aws.DescribePrivateRecordListWithPages(profile, req)
		case model.TENCENT:
			return s.Tencent.DescribePrivateRecordListWithPages(profile, req)
		default:
			return model.ListRecordsPageResponse{}, nil
		}
	}
	return model.ListRecordsPageResponse{}, fmt.Errorf("profile not found")

}

// PrivateCreateRecord
func (s *DnsService) PrivateCreateRecord(profile string, request model.CreateRecordRequest) (model.CreateRecordResponse, error) {
	if p, ok := s.Profiles[profile]; ok {
		switch p.Cloud {
		case model.AWS:
			return s.Aws.CreatePrivateRecord(profile, request)
		case model.TENCENT:
			return s.Tencent.CreatePrivateRecord(profile, request)
		default:
			return model.CreateRecordResponse{}, nil
		}
	}
	return model.CreateRecordResponse{}, fmt.Errorf("profile not found")
}

// PrivateModifyRecord
func (s *DnsService) PrivateModifyRecord(profile string, request model.ModifyRecordRequest) error {
	if p, ok := s.Profiles[profile]; ok {
		switch p.Cloud {
		case model.AWS:
			return s.Aws.ModifyPrivateRecord(profile, request)
		case model.TENCENT:
			return s.Tencent.ModifyPrivateRecord(profile, request)
		default:
			return nil
		}
	}
	return fmt.Errorf("profile not found")
}

// PrivateDeleteRecord
func (s *DnsService) PrivateDeleteRecord(profile string, request model.DeletePrivateRecordRequest) error {
	if p, ok := s.Profiles[profile]; ok {
		switch p.Cloud {
		case model.AWS:
			return s.Aws.DeletePrivateRecord(profile, request)
		case model.TENCENT:
			return s.Tencent.DeletePrivateRecord(profile, request)
		default:
			return fmt.Errorf("not support cloud %s", p.Cloud)
		}
	}
	return fmt.Errorf("profile not found")
}

// DescribeDomainList
func (s *DnsService) DescribeDomainList(profile, region string, req model.DescribeDomainListRequest) (model.DescribeDomainListResponse, error) {
	if p, ok := s.Profiles[profile]; ok {
		switch p.Cloud {
		case model.AWS:
			return s.Aws.DescribeDomainList(profile, region, req)
		case model.TENCENT:
			return s.Tencent.DescribeDomainList(profile, region, req)
		default:
			return model.DescribeDomainListResponse{}, nil
		}
	}
	return model.DescribeDomainListResponse{}, fmt.Errorf("profile not found")
}

// DescribeRecordList
func (s *DnsService) DescribeRecordList(profile, region string, req model.DescribeRecordListRequest) (model.DescribeRecordListResponse, error) {
	if p, ok := s.Profiles[profile]; ok {
		switch p.Cloud {
		case model.AWS:
			return s.Aws.DescribeRecordList(profile, region, req)
		case model.TENCENT:
			return s.Tencent.DescribeRecordList(profile, region, req)
		default:
			return model.DescribeRecordListResponse{}, nil
		}
	}
	return model.DescribeRecordListResponse{}, fmt.Errorf("profile not found")
}

// DescribeRecordListWithPages
func (s *DnsService) DescribeRecordListWithPages(profile, region string, req model.DescribeRecordListWithPageRequest) (model.ListRecordsPageResponse, error) {
	if p, ok := s.Profiles[profile]; ok {
		switch p.Cloud {
		case model.AWS:
			return s.Aws.DescribeRecordListWithPages(profile, region, req)
		case model.TENCENT:
			return s.Tencent.DescribeRecordListWithPages(profile, region, req)
		default:
			return model.ListRecordsPageResponse{}, nil
		}
	}
	return model.ListRecordsPageResponse{}, fmt.Errorf("profile not found")

}

// DescribeRecord
func (s *DnsService) DescribeRecord(profile, region string, req model.DescribeRecordRequest) (model.Record, error) {
	if p, ok := s.Profiles[profile]; ok {
		switch p.Cloud {
		case model.AWS:
			return s.Aws.DescribeRecord(profile, region, req)
		case model.TENCENT:
			resp, err := s.Tencent.DescribeRecord(profile, region, req)
			if err != nil {
				if strings.Contains(err.Error(), "ResourceNotFound") {
					return model.Record{}, nil
				}
			}
			return resp, err
		}
	}
	return model.Record{}, fmt.Errorf("profile not found")
}

// CreateRecord
func (s *DnsService) CreateRecord(profile, region string, request model.CreateRecordRequest) (model.CreateRecordResponse, error) {
	if p, ok := s.Profiles[profile]; ok {
		switch p.Cloud {
		case model.AWS:
			return s.Aws.CreateRecord(profile, region, request)
		case model.TENCENT:
			return s.Tencent.CreateRecord(profile, region, request)
		default:
			return model.CreateRecordResponse{}, nil
		}
	}
	return model.CreateRecordResponse{}, fmt.Errorf("profile not found")
}

func (s *DnsService) ModifyRecord(profile, region string, ignoreType bool, request model.ModifyRecordRequest) error {
	if p, ok := s.Profiles[profile]; ok {
		switch p.Cloud {
		case model.AWS:
			return s.Aws.ModifyRecord(profile, region, ignoreType, request)
		case model.TENCENT:
			return s.Tencent.ModifyRecord(profile, region, ignoreType, request)
		default:
			return nil
		}
	}
	return fmt.Errorf("profile not found")
}

func (s *DnsService) DeleteRecord(profile, region string, request model.DeleteRecordRequest) (model.CommonDnsResponse, error) {
	if request.Domain == nil || request.SubDomain == nil || request.RecordType == nil {
		return model.CommonDnsResponse{}, fmt.Errorf("domain, subDomain, recordType is required")
	}
	if p, ok := s.Profiles[profile]; ok {
		switch p.Cloud {
		case model.AWS:
			return s.Aws.DeleteRecord(profile, region, request)
		case model.TENCENT:
			return s.Tencent.DeleteRecord(profile, region, request)
		default:
			return model.CommonDnsResponse{}, fmt.Errorf("not support cloud %s", p.Cloud)
		}
	}
	return model.CommonDnsResponse{}, fmt.Errorf("profile not found")
}
