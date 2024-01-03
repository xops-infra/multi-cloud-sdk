package service

import (
	"github.com/xops-infra/multi-cloud-sdk/pkg/model"
)

type EmrService struct {
	Profiles     []model.ProfileConfig
	Aws, Tencent model.CloudIO
}

func NewEmrService(profiles []model.ProfileConfig, aws, tencent model.CloudIO) model.EmrContact {
	return &EmrService{
		Profiles: profiles,
		Aws:      aws,
		Tencent:  tencent,
	}
}

func (s *EmrService) DescribeEmrCluster(profile, region string, ids []*string) ([]model.DescribeEmrCluster, error) {
	for _, cfgProfile := range s.Profiles {
		if cfgProfile.Name == profile {
			switch cfgProfile.Cloud {
			case model.AWS:
				return s.Aws.DescribeEmrCluster(profile, region, ids)
			case model.TENCENT:
				return s.Tencent.DescribeEmrCluster(profile, region, ids)
			default:
				return nil, model.ErrCloudNotSupported
			}
		}
	}
	return nil, model.ErrProfileNotFound
}

func (s *EmrService) QueryEmrCluster(profile, region string, input model.EmrFilter) (model.FilterEmrResponse, error) {
	for _, cfgProfile := range s.Profiles {
		if cfgProfile.Name == profile {
			switch cfgProfile.Cloud {
			case model.AWS:
				return s.Aws.QueryEmrCluster(profile, region, input)
			case model.TENCENT:
				return s.Tencent.QueryEmrCluster(profile, region, input)
			default:
				return model.FilterEmrResponse{}, model.ErrCloudNotSupported
			}
		}
	}
	return model.FilterEmrResponse{}, model.ErrProfileNotFound
}
