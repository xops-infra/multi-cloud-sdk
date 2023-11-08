package service

import (
	"fmt"
	"sync"

	"github.com/rogpeppe/go-internal/cache"
	"github.com/xops-infra/multi-cloud-sdk/pkg/model"
)

type ServerService struct {
	Profiles     []model.ProfileConfig
	AWS, Tencent model.CloudIo
	Cache        *cache.Cache
}

func NewServer(profiles []model.ProfileConfig, aws, tencent model.CloudIo) *ServerService {
	return &ServerService{
		Profiles: profiles,
		AWS:      aws,
		Tencent:  tencent,
	}
}

func (s *ServerService) QueryInstances(input model.InstanceQueryInput) []*model.Instance {
	instances := make([]*model.Instance, 0)
	var wg = sync.WaitGroup{}
	for _, profile := range s.Profiles {
		if input.Account != "" && input.Account != profile.Name {
			// 加速，如果有指定账号，且不是当前账号，直接跳过
			continue
		}
		if profile.Cloud == model.AWS {
			for _, region := range profile.Regions {
				wg.Add(1)
				go func(profile model.ProfileConfig, region string) {
					defer wg.Done()
					_instances, err := s.AWS.QueryInstances(profile.Name, region)
					if err != nil {
						return
					}
					instances = append(instances, input.Filter(_instances)...)
				}(profile, region)

			}
		} else if profile.Cloud == model.TENCENT {
			for _, region := range profile.Regions {
				wg.Add(1)
				go func(profile model.ProfileConfig, region string) {
					defer wg.Done()
					_instances, err := s.Tencent.QueryInstances(profile.Name, region)
					if err != nil {
						return
					}
					instances = append(instances, input.Filter(_instances)...)
				}(profile, region)
			}
		}
	}
	wg.Wait()
	return instances
}

func (s *ServerService) GetInstance(instance_id string) (*model.Instance, error) {
	instances := s.QueryInstances(model.InstanceQueryInput{
		Name: instance_id,
	})
	if len(instances) == 0 {
		return nil, fmt.Errorf("instance %s not found", instance_id)
	}
	return instances[0], nil
}
