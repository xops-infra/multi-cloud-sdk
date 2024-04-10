package service

import (
	"github.com/xops-infra/multi-cloud-sdk/pkg/model"
)

func (s *CommonService) DescribeEmrCluster(input model.DescribeInput) ([]model.DescribeEmrCluster, error) {
	for _, cfgProfile := range s.Profiles {
		if cfgProfile.Name == *input.Profile {
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

func (s *CommonService) QueryEmrCluster(input model.EmrFilter) (model.FilterEmrResponse, error) {
	for _, cfgProfile := range s.Profiles {
		if cfgProfile.Name == *input.Profile {
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
