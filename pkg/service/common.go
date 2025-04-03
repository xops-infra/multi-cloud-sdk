package service

import (
	"fmt"

	"github.com/xops-infra/multi-cloud-sdk/pkg/model"
)

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

func (s *CommonService) CreateTags(profile, region string, input model.CreateTagsInput) error {
	if p, ok := s.Profiles[profile]; ok {
		switch p.Cloud {
		case model.AWS:
			return s.Aws.CreateTags(profile, region, input)
		case model.TENCENT:
			return s.Tencent.CreateTags(profile, region, input)
		}
	}
	return fmt.Errorf("profile %s not found", profile)
}

func (s *CommonService) AddTagsToResource(profile, region string, input model.AddTagsInput) error {
	if p, ok := s.Profiles[profile]; ok {
		switch p.Cloud {
		case model.AWS:
			return s.Aws.AddTagsToResource(profile, region, input)
		case model.TENCENT:
			return s.Tencent.AddTagsToResource(profile, region, input)
		}
	}
	return fmt.Errorf("profile %s not found", profile)
}

func (s *CommonService) RemoveTagsFromResource(profile, region string, input model.RemoveTagsInput) error {
	if p, ok := s.Profiles[profile]; ok {
		switch p.Cloud {
		case model.AWS:
			return s.Aws.RemoveTagsFromResource(profile, region, input)
		case model.TENCENT:
			return s.Tencent.RemoveTagsFromResource(profile, region, input)
		}
	}
	return fmt.Errorf("profile %s not found", profile)
}

func (s *CommonService) ModifyTagsForResource(profile, region string, input model.ModifyTagsInput) error {
	if p, ok := s.Profiles[profile]; ok {
		switch p.Cloud {
		case model.AWS:
			return s.Aws.ModifyTagsForResource(profile, region, input)
		case model.TENCENT:
			return s.Tencent.ModifyTagsForResource(profile, region, input)
		}
	}
	return fmt.Errorf("profile %s not found", profile)
}
