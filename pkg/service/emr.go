package service

import (
	"fmt"

	"github.com/xops-infra/multi-cloud-sdk/pkg/model"
)

func (s *CommonService) DescribeEmrCluster(input model.DescribeInput) ([]model.DescribeEmrCluster, error) {
	if p, ok := s.Profiles[*input.Profile]; ok {
		switch p.Cloud {
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
			return nil, fmt.Errorf("%s %s", *input.Profile, model.ErrCloudNotSupported.Error())
		}
	}
	return nil, fmt.Errorf("%s %s", *input.Profile, model.ErrProfileNotFound.Error())
}

func (s *CommonService) QueryEmrCluster(input model.EmrFilter) (model.FilterEmrResponse, error) {
	if p, ok := s.Profiles[*input.Profile]; ok {
		switch p.Cloud {
		case model.AWS:
			return s.Aws.QueryEmrCluster(input)
		case model.TENCENT:
			return s.Tencent.QueryEmrCluster(input)
		default:
			return model.FilterEmrResponse{}, fmt.Errorf("%s %s", *input.Profile, model.ErrCloudNotSupported.Error())
		}
	}
	return model.FilterEmrResponse{}, fmt.Errorf("%s %s", *input.Profile, model.ErrProfileNotFound.Error())
}
