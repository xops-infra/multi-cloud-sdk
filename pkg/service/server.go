package service

import (
	"sync"

	"github.com/rogpeppe/go-internal/cache"

	"github.com/xops-infra/multi-cloud-sdk/pkg/model"
)

// 汇总所有账号的机器信息
type ServerService struct {
	Profiles     []model.ProfileConfig
	AWS, Tencent model.CloudIO
	Cache        *cache.Cache
}

func NewServer(profiles []model.ProfileConfig, aws, tencent model.CloudIO) model.InstanceContact {
	return &ServerService{
		Profiles: profiles,
		AWS:      aws,
		Tencent:  tencent,
	}
}

func (s *ServerService) QueryInstances(input model.InstanceFilter) ([]model.Instance, error) {
	instances := make([]model.Instance, 0)
	wg := sync.WaitGroup{}
	for _, profile := range s.Profiles {
		if input.Profile != nil && *input.Profile != profile.Name {
			// 加速，如果有指定账号，且不是当前账号，直接跳过
			continue
		}
		if profile.Cloud == model.AWS {
			for _, region := range profile.Regions {
				wg.Add(1)
				go func(profile model.ProfileConfig, region string) {
					defer wg.Done()
					req := input.ToAwsDescribeInstancesInput()
					for {
						_instances, err := s.AWS.DescribeInstances(profile.Name, region, req)
						if err != nil {
							return
						}
						instances = append(instances, _instances.Instances...)
						if _instances.NextMarker.(*string) == nil {
							break
						} else {
							req.NextMarker = _instances.NextMarker
						}
					}
				}(profile, region)
			}
		} else if profile.Cloud == model.TENCENT {
			for _, region := range profile.Regions {
				wg.Add(1)
				go func(profile model.ProfileConfig, region string) {
					// 腾讯因为分页使用偏移量，所以不需要循环，直接获取所有数据即可
					defer wg.Done()
					_instances, err := s.Tencent.DescribeInstances(profile.Name, region, input.ToTxDescribeInstancesInput())
					if err != nil {
						return
					}
					instances = append(instances, _instances.Instances...)
				}(profile, region)
			}
		}
		wg.Wait()
	}
	return instances, nil
}
