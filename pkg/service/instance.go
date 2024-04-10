package service

import (
	"github.com/xops-infra/multi-cloud-sdk/pkg/model"
)

func (s *CommonService) DescribeInstances(profile, region string, input model.InstanceFilter) (model.InstanceResponse, error) {
	for _, cfgProfile := range s.Profiles {
		if cfgProfile.Name == profile {
			switch cfgProfile.Cloud {
			case model.AWS:
				return s.Aws.DescribeInstances(profile, region, input.ToAwsDescribeInstancesInput())
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
func (s *CommonService) CreateInstance(profile, region string, input model.CreateInstanceInput) (model.CreateInstanceResponse, error) {
	for _, cfgProfile := range s.Profiles {
		if cfgProfile.Name == profile {
			switch cfgProfile.Cloud {
			case model.AWS:
				return s.Aws.CreateInstance(profile, region, input)
			case model.TENCENT:
				return s.Tencent.CreateInstance(profile, region, input)
			default:
				return model.CreateInstanceResponse{}, model.ErrCloudNotSupported
			}
		}
	}
	return model.CreateInstanceResponse{}, model.ErrProfileNotFound
}

func (s *CommonService) ModifyInstance(profile, region string, input model.ModifyInstanceInput) (model.ModifyInstanceResponse, error) {
	for _, cfgProfile := range s.Profiles {
		if cfgProfile.Name == profile {
			switch cfgProfile.Cloud {
			case model.AWS:
				return s.Aws.ModifyInstance(profile, region, input)
			case model.TENCENT:
				return s.Tencent.ModifyInstance(profile, region, input)
			default:
				return model.ModifyInstanceResponse{}, model.ErrCloudNotSupported
			}
		}
	}
	return model.ModifyInstanceResponse{}, model.ErrProfileNotFound
}

func (s *CommonService) DeleteInstance(profile, region string, input model.DeleteInstanceInput) (model.DeleteInstanceResponse, error) {
	for _, cfgProfile := range s.Profiles {
		if cfgProfile.Name == profile {
			switch cfgProfile.Cloud {
			case model.AWS:
				return s.Aws.DeleteInstance(profile, region, input)
			case model.TENCENT:
				return s.Tencent.DeleteInstance(profile, region, input)
			default:
				return model.DeleteInstanceResponse{}, model.ErrCloudNotSupported
			}
		}
	}
	return model.DeleteInstanceResponse{}, model.ErrProfileNotFound
}
