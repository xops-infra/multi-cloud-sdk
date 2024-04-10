package service

import "github.com/xops-infra/multi-cloud-sdk/pkg/model"

type CommonService struct {
	Profiles     map[string]model.ProfileConfig
	Aws, Tencent model.CloudIO
}

func NewCommonService(profiles []model.ProfileConfig, aws, tencent model.CloudIO) model.CommonContract {
	_profiles := make(map[string]model.ProfileConfig)
	for _, p := range profiles {
		_profiles[p.Name] = p

	}
	return &CommonService{
		Profiles: _profiles,
		Aws:      aws,
		Tencent:  tencent,
	}
}
