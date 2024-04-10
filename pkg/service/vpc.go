package service

import (
	"github.com/xops-infra/multi-cloud-sdk/pkg/model"
)

func (s *CommonService) QueryVPCs(profile, region string, input model.CommonFilter) ([]model.VPC, error) {
	for _, cfgProfile := range s.Profiles {
		if cfgProfile.Name == profile {
			switch cfgProfile.Cloud {
			case model.AWS:
				return s.Aws.QueryVPC(profile, region, input)
			case model.TENCENT:
				return s.Tencent.QueryVPC(profile, region, input)
			default:
				return nil, model.ErrCloudNotSupported
			}
		}
	}
	return nil, model.ErrProfileNotFound
}

func (s *CommonService) QueryEIPs(profile, region string, input model.CommonFilter) ([]model.EIP, error) {
	for _, cfgProfile := range s.Profiles {
		if cfgProfile.Name == profile {
			switch cfgProfile.Cloud {
			case model.AWS:
				return s.Aws.QueryEIP(profile, region, input)
			case model.TENCENT:
				return s.Tencent.QueryEIP(profile, region, input)
			default:
				return nil, model.ErrCloudNotSupported
			}
		}
	}
	return nil, model.ErrProfileNotFound
}

func (s *CommonService) QueryNATs(profile, region string, input model.CommonFilter) (nats []model.NAT, err error) {
	for _, cfgProfile := range s.Profiles {
		if cfgProfile.Name == profile {
			switch cfgProfile.Cloud {
			case model.AWS:
				return s.Aws.QueryNAT(profile, region, input)
			case model.TENCENT:
				return s.Tencent.QueryNAT(profile, region, input)
			default:
				return nil, model.ErrCloudNotSupported
			}
		}
	}
	return nil, model.ErrProfileNotFound
}

func (s *CommonService) QuerySubnets(profile, region string, input model.CommonFilter) (subnets []model.Subnet, err error) {
	for _, cfgProfile := range s.Profiles {
		if cfgProfile.Name == profile {
			switch cfgProfile.Cloud {
			case model.AWS:
				return s.Aws.QuerySubnet(profile, region, input)
			case model.TENCENT:
				return s.Tencent.QuerySubnet(profile, region, input)
			default:
				return nil, model.ErrCloudNotSupported
			}
		}
	}
	return nil, model.ErrProfileNotFound
}
