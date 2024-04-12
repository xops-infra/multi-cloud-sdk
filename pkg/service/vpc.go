package service

import (
	"fmt"

	"github.com/xops-infra/multi-cloud-sdk/pkg/model"
)

func (s *CommonService) QueryVPCs(profile, region string, input model.CommonFilter) ([]model.VPC, error) {
	if p, ok := s.Profiles[profile]; ok {
		switch p.Cloud {
		case model.AWS:
			return s.Aws.QueryVPC(profile, region, input)
		case model.TENCENT:
			return s.Tencent.QueryVPC(profile, region, input)
		default:
			return nil, fmt.Errorf("%s %s", profile, model.ErrCloudNotSupported.Error())
		}
	}
	return nil, fmt.Errorf("%s %s", profile, model.ErrProfileNotFound.Error())
}

func (s *CommonService) QueryEIPs(profile, region string, input model.CommonFilter) ([]model.EIP, error) {
	if p, ok := s.Profiles[profile]; ok {
		switch p.Cloud {
		case model.AWS:
			return s.Aws.QueryEIP(profile, region, input)
		case model.TENCENT:
			return s.Tencent.QueryEIP(profile, region, input)
		default:
			return nil, fmt.Errorf("%s %s", profile, model.ErrCloudNotSupported.Error())
		}
	}
	return nil, fmt.Errorf("%s %s", profile, model.ErrProfileNotFound.Error())
}

func (s *CommonService) QueryNATs(profile, region string, input model.CommonFilter) (nats []model.NAT, err error) {
	if p, ok := s.Profiles[profile]; ok {
		switch p.Cloud {
		case model.AWS:
			return s.Aws.QueryNAT(profile, region, input)
		case model.TENCENT:
			return s.Tencent.QueryNAT(profile, region, input)
		default:
			return nil, fmt.Errorf("%s %s", profile, model.ErrCloudNotSupported.Error())
		}
	}
	return nil, fmt.Errorf("%s %s", profile, model.ErrProfileNotFound.Error())
}

func (s *CommonService) QuerySubnets(profile, region string, input model.CommonFilter) (subnets []model.Subnet, err error) {
	if p, ok := s.Profiles[profile]; ok {
		switch p.Cloud {
		case model.AWS:
			return s.Aws.QuerySubnet(profile, region, input)
		case model.TENCENT:
			return s.Tencent.QuerySubnet(profile, region, input)
		default:
			return nil, fmt.Errorf("%s %s", profile, model.ErrCloudNotSupported.Error())
		}
	}
	return nil, fmt.Errorf("%s %s", profile, model.ErrProfileNotFound.Error())
}
