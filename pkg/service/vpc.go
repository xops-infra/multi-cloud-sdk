package service

import (
	"github.com/xops-infra/multi-cloud-sdk/pkg/model"
)

type VpcService struct {
	Profiles     []model.ProfileConfig
	Aws, Tencent model.CloudIO
}

func NewVpcService(profiles []model.ProfileConfig, aws, tencent model.CloudIO) model.VpcContract {
	return &VpcService{
		Profiles: profiles,
		Aws:      aws,
		Tencent:  tencent,
	}
}

func (s *VpcService) QueryVPCs(profile, region string, input model.CommonFilter) ([]model.VPC, error) {
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

func (s *VpcService) QueryEIPs(profile, region string, input model.CommonFilter) ([]model.EIP, error) {
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

func (s *VpcService) QueryNATs(profile, region string, input model.CommonFilter) (nats []model.NAT, err error) {
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

func (s *VpcService) QuerySubnets(profile, region string, input model.CommonFilter) (subnets []model.Subnet, err error) {
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
