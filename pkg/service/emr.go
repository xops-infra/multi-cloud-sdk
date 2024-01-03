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

func (s *EmrService) DescribeEmrCluster(input model.DescribeInput) ([]model.DescribeEmrCluster, error) {
	for _, cfgProfile := range s.Profiles {
		if cfgProfile.Name == *input.Profile {
			if input.Region == nil {
				if len(cfgProfile.Regions) == 0 {
					return nil, model.ErrRegionNotConfigured
				}
				input.Region = &cfgProfile.Regions[0] // 默认取第一个
			}

			switch cfgProfile.Cloud {
			case model.AWS:
				return s.Aws.DescribeEmrCluster(model.DescribeInput{
					Profile: input.Profile,
					Region:  input.Region,
					IDS:     input.IDS,
				})
			case model.TENCENT:
				return s.Tencent.DescribeEmrCluster(model.DescribeInput{
					Profile: input.Profile,
					Region:  input.Region,
					IDS:     input.IDS,
				})
			default:
				return nil, model.ErrCloudNotSupported
			}
		}
	}
	return nil, model.ErrProfileNotFound
}

func (s *EmrService) QueryEmrCluster(input model.EmrFilter) (model.FilterEmrResponse, error) {
	for _, cfgProfile := range s.Profiles {
		if cfgProfile.Name == *input.Profile {
			if input.Region == nil {
				if len(cfgProfile.Regions) == 0 {
					return model.FilterEmrResponse{}, model.ErrRegionNotConfigured
				}
				input.Region = &cfgProfile.Regions[0] // 默认取第一个
			}
			switch cfgProfile.Cloud {
			case model.AWS:
				return s.Aws.QueryEmrCluster(input)
			case model.TENCENT:
				return s.Tencent.QueryEmrCluster(input)
			default:
				return model.FilterEmrResponse{}, model.ErrCloudNotSupported
			}
		}
	}
	return model.FilterEmrResponse{}, model.ErrProfileNotFound
}
