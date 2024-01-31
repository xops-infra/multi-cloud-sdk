package service

import (
	"github.com/rogpeppe/go-internal/cache"

	"github.com/xops-infra/multi-cloud-sdk/pkg/model"
)

// 汇总所有账号的机器信息
type ServerService struct {
	Profiles     []model.ProfileConfig
	AWS, Tencent model.CloudIO
	Cache        *cache.Cache
}

func NewServer(profiles []model.ProfileConfig, aws, tencent model.CloudIO) model.InstanceContact {
	return &ServerService{
		Profiles: profiles,
		AWS:      aws,
		Tencent:  tencent,
	}
}

func (s *ServerService) DescribeInstances(profile, region string, input model.InstanceFilter) (model.InstanceResponse, error) {
	for _, cfgProfile := range s.Profiles {
		if cfgProfile.Name == profile {
			switch cfgProfile.Cloud {
			case model.AWS:
				return s.AWS.DescribeInstances(profile, region, input.ToAwsDescribeInstancesInput())
			case model.TENCENT:
				return s.Tencent.DescribeInstances(profile, region, input.ToTxDescribeInstancesInput())
			default:
				return model.InstanceResponse{}, model.ErrCloudNotSupported
			}
		}
	}
	return model.InstanceResponse{}, model.ErrProfileNotFound
}

// CreateInstance
func (s *ServerService) CreateInstance(profile, region string, input model.CreateInstanceInput) (model.CreateInstanceResponse, error) {
	for _, cfgProfile := range s.Profiles {
		if cfgProfile.Name == profile {
			switch cfgProfile.Cloud {
			case model.AWS:
				return s.AWS.CreateInstance(profile, region, input)
			case model.TENCENT:
				return s.Tencent.CreateInstance(profile, region, input)
			default:
				return model.CreateInstanceResponse{}, model.ErrCloudNotSupported
			}
		}
	}
	return model.CreateInstanceResponse{}, model.ErrProfileNotFound
}

func (s *ServerService) ModifyInstance(profile, region string, input model.ModifyInstanceInput) (model.ModifyInstanceResponse, error) {
	for _, cfgProfile := range s.Profiles {
		if cfgProfile.Name == profile {
			switch cfgProfile.Cloud {
			case model.AWS:
				return s.AWS.ModifyInstance(profile, region, input)
			case model.TENCENT:
				return s.Tencent.ModifyInstance(profile, region, input)
			default:
				return model.ModifyInstanceResponse{}, model.ErrCloudNotSupported
			}
		}
	}
	return model.ModifyInstanceResponse{}, model.ErrProfileNotFound
}

func (s *ServerService) DeleteInstance(profile, region string, input model.DeleteInstanceInput) (model.DeleteInstanceResponse, error) {
	for _, cfgProfile := range s.Profiles {
		if cfgProfile.Name == profile {
			switch cfgProfile.Cloud {
			case model.AWS:
				return s.AWS.DeleteInstance(profile, region, input)
			case model.TENCENT:
				return s.Tencent.DeleteInstance(profile, region, input)
			default:
				return model.DeleteInstanceResponse{}, model.ErrCloudNotSupported
			}
		}
	}
	return model.DeleteInstanceResponse{}, model.ErrProfileNotFound
}
