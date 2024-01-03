package service

import (
	"fmt"
	"sync"

	"github.com/xops-infra/multi-cloud-sdk/pkg/model"
)

type VpcService struct {
	Profiles     []model.ProfileConfig
	Aws, Tencent model.CloudIO
}

func NewVpcService(profiles []model.ProfileConfig, aws, tencent model.CloudIO) model.VpcContract {
	return &VpcService{
		Profiles: profiles,
		Aws:      aws,
		Tencent:  tencent,
	}
}

func (s *VpcService) QueryVPCs(input model.CommonFilter) (vpcs []model.VPC, err error) {
	var wg = sync.WaitGroup{}
	for _, profile := range s.Profiles {
		if profile.Cloud == model.AWS {
			for _, region := range profile.Regions {
				wg.Add(1)
				go func(profile model.ProfileConfig, region string) {
					defer wg.Done()
					vpc, err := s.Aws.QueryVPC(profile.Name, region, input)
					if err != nil {
						err = fmt.Errorf("aws query vpc error: %v", err)
						return
					}
					vpcs = append(vpcs, vpc...)
				}(profile, region)
			}
		} else if profile.Cloud == model.TENCENT {
			for _, region := range profile.Regions {
				wg.Add(1)
				go func(profile model.ProfileConfig, region string) {
					defer wg.Done()
					vpc, err := s.Tencent.QueryVPC(profile.Name, region, input)
					if err != nil {
						err = fmt.Errorf("tencent query vpc error: %v", err)
						return
					}
					vpcs = append(vpcs, vpc...)
				}(profile, region)
			}
		}
	}
	wg.Wait()
	return
}

func (s *VpcService) GetVPC(vpc_id string) (model.VPC, error) {
	vpcs, err := s.QueryVPCs(model.CommonFilter{
		ID: vpc_id,
	})
	if err != nil {
		return model.VPC{}, err
	}
	if len(vpcs) == 1 {
		return vpcs[0], nil
	}
	return model.VPC{}, fmt.Errorf("vpc not found,or multiple vpcs found")
}

func (s *VpcService) QueryEIPs(input model.CommonFilter) ([]model.EIP, error) {
	eips := []model.EIP{}
	var wg = sync.WaitGroup{}
	var err error
	for _, profile := range s.Profiles {
		if profile.Cloud == model.AWS {
			for _, region := range profile.Regions {
				wg.Add(1)
				go func(profile model.ProfileConfig, region string, err error) {
					defer wg.Done()
					eip, err := s.Aws.QueryEIP(profile.Name, region, input)
					if err != nil {
						return
					}
					eips = append(eips, eip...)
				}(profile, region, err)
			}
		} else if profile.Cloud == model.TENCENT {
			for _, region := range profile.Regions {
				wg.Add(1)
				go func(profile model.ProfileConfig, region string) {
					defer wg.Done()
					eip, err := s.Tencent.QueryEIP(profile.Name, region, input)
					if err != nil {
						err = fmt.Errorf("tencent query eip error: %v", err)
						return
					}
					eips = append(eips, eip...)
				}(profile, region)
			}
		}
	}
	wg.Wait()
	return eips, err
}

func (s *VpcService) QueryNATs(input model.CommonFilter) (nats []model.NAT, err error) {
	var wg = sync.WaitGroup{}
	for _, profile := range s.Profiles {
		if profile.Cloud == model.AWS {
			for _, region := range profile.Regions {
				wg.Add(1)
				go func(profile model.ProfileConfig, region string) {
					defer wg.Done()
					nat, err := s.Aws.QueryNAT(profile.Name, region, input)
					if err != nil {
						err = fmt.Errorf("aws query nat error: %v", err)
						return
					}
					nats = append(nats, nat...)
				}(profile, region)
			}
		} else if profile.Cloud == model.TENCENT {
			for _, region := range profile.Regions {
				wg.Add(1)
				go func(profile model.ProfileConfig, region string) {
					defer wg.Done()
					nat, err := s.Tencent.QueryNAT(profile.Name, region, input)
					if err != nil {
						err = fmt.Errorf("tencent query nat error: %v", err)
						return
					}
					nats = append(nats, nat...)
				}(profile, region)
			}
		}
	}
	wg.Wait()
	return
}

func (s *VpcService) QuerySubnets(input model.CommonFilter) (subnets []model.Subnet, err error) {
	var wg = sync.WaitGroup{}
	for _, profile := range s.Profiles {
		if profile.Cloud == model.AWS {
			for _, region := range profile.Regions {
				wg.Add(1)
				go func(profile model.ProfileConfig, region string) {
					defer wg.Done()
					subnet, err := s.Aws.QuerySubnet(profile.Name, region, input)
					if err != nil {
						err = fmt.Errorf("aws query subnet error: %v", err)
						return
					}
					subnets = append(subnets, subnet...)
				}(profile, region)
			}
		} else if profile.Cloud == model.TENCENT {
			for _, region := range profile.Regions {
				wg.Add(1)
				go func(profile model.ProfileConfig, region string) {
					defer wg.Done()
					subnet, err := s.Tencent.QuerySubnet(profile.Name, region, input)
					if err != nil {
						err = fmt.Errorf("tencent query subnet error: %v", err)
						return
					}
					subnets = append(subnets, subnet...)
				}(profile, region)
			}
		}
	}
	wg.Wait()
	return
}
