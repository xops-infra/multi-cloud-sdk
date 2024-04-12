package service

import (
	"fmt"
	"strings"

	"github.com/xops-infra/multi-cloud-sdk/pkg/model"
)

// PrivateDomainList
func (s *CommonService) PrivateDomainList(profile string, req model.DescribeDomainListRequest) (model.DescribePrivateDomainListResponse, error) {
	if p, ok := s.Profiles[profile]; ok {
		switch p.Cloud {
		case model.AWS:
			return s.Aws.DescribePrivateDomainList(profile, req)
		case model.TENCENT:
			return s.Tencent.DescribePrivateDomainList(profile, req)
		default:
			return model.DescribePrivateDomainListResponse{}, fmt.Errorf("%s %s", profile, model.ErrCloudNotSupported.Error())
		}
	}
	return model.DescribePrivateDomainListResponse{}, fmt.Errorf("%s %s", profile, model.ErrProfileNotFound.Error())
}

// PrivateRecordList
func (s *CommonService) PrivateRecordList(profile string, req model.DescribePrivateRecordListRequest) (model.DescribePrivateRecordListResponse, error) {
	if p, ok := s.Profiles[profile]; ok {
		switch p.Cloud {
		case model.AWS:
			return s.Aws.DescribePrivateRecordList(profile, req)
		case model.TENCENT:
			return s.Tencent.DescribePrivateRecordList(profile, req)
		default:
			return model.DescribePrivateRecordListResponse{}, fmt.Errorf("%s %s", profile, model.ErrCloudNotSupported.Error())
		}
	}
	return model.DescribePrivateRecordListResponse{}, fmt.Errorf("%s %s", profile, model.ErrProfileNotFound.Error())
}

// PrivateRecordListWithPages
func (s *CommonService) PrivateRecordListWithPages(profile string, req model.DescribePrivateDnsRecordListWithPageRequest) (model.ListRecordsPageResponse, error) {
	if p, ok := s.Profiles[profile]; ok {
		switch p.Cloud {
		case model.AWS:
			return s.Aws.DescribePrivateRecordListWithPages(profile, req)
		case model.TENCENT:
			return s.Tencent.DescribePrivateRecordListWithPages(profile, req)
		default:
			return model.ListRecordsPageResponse{}, fmt.Errorf("%s %s", profile, model.ErrCloudNotSupported.Error())
		}
	}
	return model.ListRecordsPageResponse{}, fmt.Errorf("%s %s", profile, model.ErrProfileNotFound.Error())

}

// PrivateCreateRecord
func (s *CommonService) PrivateCreateRecord(profile string, request model.CreateRecordRequest) (model.CreateRecordResponse, error) {
	if p, ok := s.Profiles[profile]; ok {
		switch p.Cloud {
		case model.AWS:
			return s.Aws.CreatePrivateRecord(profile, request)
		case model.TENCENT:
			return s.Tencent.CreatePrivateRecord(profile, request)
		default:
			return model.CreateRecordResponse{}, fmt.Errorf("%s %s", profile, model.ErrCloudNotSupported.Error())
		}
	}
	return model.CreateRecordResponse{}, fmt.Errorf("%s %s", profile, model.ErrProfileNotFound.Error())
}

// PrivateModifyRecord
func (s *CommonService) PrivateModifyRecord(profile string, request model.ModifyRecordRequest) error {
	if p, ok := s.Profiles[profile]; ok {
		switch p.Cloud {
		case model.AWS:
			return s.Aws.ModifyPrivateRecord(profile, request)
		case model.TENCENT:
			return s.Tencent.ModifyPrivateRecord(profile, request)
		default:
			return fmt.Errorf("%s %s", profile, model.ErrCloudNotSupported.Error())
		}
	}
	return fmt.Errorf("%s %s", profile, model.ErrProfileNotFound.Error())
}

// PrivateDeleteRecord
func (s *CommonService) PrivateDeleteRecord(profile string, request model.DeletePrivateRecordRequest) error {
	if p, ok := s.Profiles[profile]; ok {
		switch p.Cloud {
		case model.AWS:
			return s.Aws.DeletePrivateRecord(profile, request)
		case model.TENCENT:
			return s.Tencent.DeletePrivateRecord(profile, request)
		default:
			return fmt.Errorf("%s %s", profile, model.ErrCloudNotSupported.Error())
		}
	}
	return fmt.Errorf("%s %s", profile, model.ErrProfileNotFound.Error())
}

// DescribeDomainList
func (s *CommonService) DescribeDomainList(profile, region string, req model.DescribeDomainListRequest) (model.DescribeDomainListResponse, error) {
	if p, ok := s.Profiles[profile]; ok {
		switch p.Cloud {
		case model.AWS:
			return s.Aws.DescribeDomainList(profile, region, req)
		case model.TENCENT:
			return s.Tencent.DescribeDomainList(profile, region, req)
		default:
			return model.DescribeDomainListResponse{}, fmt.Errorf("%s %s", profile, model.ErrCloudNotSupported.Error())
		}
	}
	return model.DescribeDomainListResponse{}, fmt.Errorf("%s %s", profile, model.ErrProfileNotFound.Error())
}

// DescribeRecordList
func (s *CommonService) DescribeRecordList(profile, region string, req model.DescribeRecordListRequest) (model.DescribeRecordListResponse, error) {
	if p, ok := s.Profiles[profile]; ok {
		switch p.Cloud {
		case model.AWS:
			return s.Aws.DescribeRecordList(profile, region, req)
		case model.TENCENT:
			return s.Tencent.DescribeRecordList(profile, region, req)
		default:
			return model.DescribeRecordListResponse{}, fmt.Errorf("%s %s", profile, model.ErrCloudNotSupported.Error())
		}
	}
	return model.DescribeRecordListResponse{}, fmt.Errorf("%s %s", profile, model.ErrProfileNotFound.Error())
}

// DescribeRecordListWithPages
func (s *CommonService) DescribeRecordListWithPages(profile, region string, req model.DescribeRecordListWithPageRequest) (model.ListRecordsPageResponse, error) {
	if p, ok := s.Profiles[profile]; ok {
		switch p.Cloud {
		case model.AWS:
			return s.Aws.DescribeRecordListWithPages(profile, region, req)
		case model.TENCENT:
			return s.Tencent.DescribeRecordListWithPages(profile, region, req)
		default:
			return model.ListRecordsPageResponse{}, fmt.Errorf("%s %s", profile, model.ErrCloudNotSupported.Error())
		}
	}

	return model.ListRecordsPageResponse{}, fmt.Errorf("%s %s", profile, model.ErrProfileNotFound.Error())
}

// DescribeRecord
func (s *CommonService) DescribeRecord(profile, region string, req model.DescribeRecordRequest) (model.Record, error) {
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
		default:
			return model.Record{}, fmt.Errorf("%s %s", profile, model.ErrCloudNotSupported.Error())
		}
	}
	return model.Record{}, fmt.Errorf("%s %s", profile, model.ErrProfileNotFound.Error())
}

// CreateRecord
func (s *CommonService) CreateRecord(profile, region string, request model.CreateRecordRequest) (model.CreateRecordResponse, error) {
	if p, ok := s.Profiles[profile]; ok {
		switch p.Cloud {
		case model.AWS:
			return s.Aws.CreateRecord(profile, region, request)
		case model.TENCENT:
			return s.Tencent.CreateRecord(profile, region, request)
		default:
			return model.CreateRecordResponse{}, fmt.Errorf("%s %s", profile, model.ErrCloudNotSupported.Error())
		}
	}
	return model.CreateRecordResponse{}, fmt.Errorf("%s %s", profile, model.ErrProfileNotFound.Error())
}

func (s *CommonService) ModifyRecord(profile, region string, ignoreType bool, request model.ModifyRecordRequest) error {
	if p, ok := s.Profiles[profile]; ok {
		switch p.Cloud {
		case model.AWS:
			return s.Aws.ModifyRecord(profile, region, ignoreType, request)
		case model.TENCENT:
			return s.Tencent.ModifyRecord(profile, region, ignoreType, request)
		default:
			return fmt.Errorf("%s %s", profile, model.ErrCloudNotSupported.Error())
		}
	}
	return fmt.Errorf("%s %s", profile, model.ErrProfileNotFound.Error())
}

func (s *CommonService) DeleteRecord(profile, region string, request model.DeleteRecordRequest) (model.CommonDnsResponse, error) {
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
			return model.CommonDnsResponse{}, fmt.Errorf("%s %s", profile, model.ErrCloudNotSupported.Error())
		}
	}
	return model.CommonDnsResponse{}, fmt.Errorf("%s %s", profile, model.ErrProfileNotFound.Error())
}
