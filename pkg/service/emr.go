package service

import (
	"github.com/xops-infra/multi-cloud-sdk/config"
	"github.com/xops-infra/multi-cloud-sdk/pkg/model"
)

type EmrService struct {
	Config       config.ServerConfig
	Aws, Tencent model.CloudIO
}

func NewEmrService(cfg config.ServerConfig) model.EmrContact {
	return &EmrService{
		Config: cfg,
	}
}

func (s *EmrService) DescribeEmrCluster(profile, region string, ids []*string) ([]model.DescribeEmrCluster, error) {
	for _, cfgProfile := range s.Config.Profiles {
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
	for _, cfgProfile := range s.Config.Profiles {
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
