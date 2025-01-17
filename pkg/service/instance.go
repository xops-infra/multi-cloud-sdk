package service

import (
	"fmt"

	"github.com/xops-infra/multi-cloud-sdk/pkg/model"
)

func (s *CommonService) DescribeInstances(profile, region string, input model.InstanceFilter) (model.InstanceResponse, error) {
	if p, ok := s.Profiles[profile]; ok {
		switch p.Cloud {
		case model.AWS:
			return s.Aws.DescribeInstances(profile, region, input.ToAwsDescribeInstancesInput())
		case model.TENCENT:
			return s.Tencent.DescribeInstances(profile, region, input.ToTxDescribeInstancesInput())
		default:
			return model.InstanceResponse{}, fmt.Errorf("%s %s", profile, model.ErrCloudNotSupported.Error())
		}
	}
	return model.InstanceResponse{}, fmt.Errorf("%s %s", profile, model.ErrProfileNotFound.Error())
}

// CreateInstance
func (s *CommonService) CreateInstance(profile, region string, input model.CreateInstanceInput) (model.CreateInstanceResponse, error) {
	if p, ok := s.Profiles[profile]; ok {
		switch p.Cloud {
		case model.AWS:
			return s.Aws.CreateInstance(profile, region, input)
		case model.TENCENT:
			return s.Tencent.CreateInstance(profile, region, input)
		default:
			return model.CreateInstanceResponse{}, fmt.Errorf("%s %s", profile, model.ErrCloudNotSupported.Error())
		}
	}
	return model.CreateInstanceResponse{}, fmt.Errorf("%s %s", profile, model.ErrProfileNotFound.Error())
}

func (s *CommonService) ModifyInstance(profile, region string, input model.ModifyInstanceInput) (model.ModifyInstanceResponse, error) {
	if p, ok := s.Profiles[profile]; ok {
		switch p.Cloud {
		case model.AWS:
			return s.Aws.ModifyInstance(profile, region, input)
		case model.TENCENT:
			return s.Tencent.ModifyInstance(profile, region, input)
		default:
			return model.ModifyInstanceResponse{}, fmt.Errorf("%s %s", profile, model.ErrCloudNotSupported.Error())
		}
	}
	return model.ModifyInstanceResponse{}, fmt.Errorf("%s %s", profile, model.ErrProfileNotFound.Error())
}

func (s *CommonService) DeleteInstance(profile, region string, input model.DeleteInstanceInput) (model.DeleteInstanceResponse, error) {
	if p, ok := s.Profiles[profile]; ok {
		switch p.Cloud {
		case model.AWS:
			return s.Aws.DeleteInstance(profile, region, input)
		case model.TENCENT:
			return s.Tencent.DeleteInstance(profile, region, input)
		default:
			return model.DeleteInstanceResponse{}, fmt.Errorf("%s %s", profile, model.ErrCloudNotSupported.Error())
		}
	}
	return model.DeleteInstanceResponse{}, fmt.Errorf("%s %s", profile, model.ErrProfileNotFound.Error())
}

func (s *CommonService) DescribeKeyPairs(profile, region string, input model.CommonFilter) ([]model.KeyPair, error) {
	if p, ok := s.Profiles[profile]; ok {
		switch p.Cloud {
		case model.AWS:
			return s.Aws.DescribeKeyPairs(profile, region, input)
		case model.TENCENT:
			return s.Tencent.DescribeKeyPairs(profile, region, input)
		default:
			return nil, fmt.Errorf("%s %s", profile, model.ErrCloudNotSupported.Error())
		}
	}
	return nil, fmt.Errorf("%s %s", profile, model.ErrProfileNotFound.Error())
}

func (s *CommonService) DescribeImages(profile, region string, input model.CommonFilter) ([]model.Image, error) {
	if p, ok := s.Profiles[profile]; ok {
		switch p.Cloud {
		case model.AWS:
			return s.Aws.DescribeImages(profile, region, input)
		case model.TENCENT:
			return s.Tencent.DescribeImages(profile, region, input)
		default:
			return nil, fmt.Errorf("%s %s", profile, model.ErrCloudNotSupported.Error())
		}
	}
	return nil, fmt.Errorf("%s %s", profile, model.ErrProfileNotFound.Error())
}
