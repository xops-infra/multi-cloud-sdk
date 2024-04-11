package service

import (
	"github.com/xops-infra/multi-cloud-sdk/pkg/model"
)

func (s *CommonService) ListBuckets(profile, region string, input model.ListBucketRequest) (model.ListBucketResponse, error) {
	for _, cfgProfile := range s.Profiles {
		if cfgProfile.Name == profile {
			switch cfgProfile.Cloud {
			case model.AWS:
				return s.Aws.ListBucket(profile, region, input)
			case model.TENCENT:
				return s.Tencent.ListBucket(profile, region, input)
			default:
				return model.ListBucketResponse{}, model.ErrCloudNotSupported
			}
		}
	}
	return model.ListBucketResponse{}, model.ErrProfileNotFound
}

func (s *CommonService) CreateBucket(profile, region string, input model.CreateBucketRequest) error {
	for _, cfgProfile := range s.Profiles {
		if cfgProfile.Name == profile {
			switch cfgProfile.Cloud {
			case model.AWS:
				return s.Aws.CreateBucket(profile, region, input)
			case model.TENCENT:
				return s.Tencent.CreateBucket(profile, region, input)
			default:
				return model.ErrCloudNotSupported
			}
		}
	}
	return model.ErrProfileNotFound
}

func (s *CommonService) DeleteBucket(profile, region string, input model.DeleteBucketRequest) (model.DeleteBucketResponse, error) {
	for _, cfgProfile := range s.Profiles {
		if cfgProfile.Name == profile {
			switch cfgProfile.Cloud {
			case model.AWS:
				return s.Aws.DeleteBucket(profile, region, input)
			case model.TENCENT:
				return s.Tencent.DeleteBucket(profile, region, input)
			default:
				return model.DeleteBucketResponse{}, model.ErrCloudNotSupported
			}
		}
	}
	return model.DeleteBucketResponse{}, model.ErrProfileNotFound
}

func (s *CommonService) GetObjectPregisn(profile, region string, input model.ObjectPregisnRequest) (model.ObjectPregisnResponse, error) {
	for _, cfgProfile := range s.Profiles {
		if cfgProfile.Name == profile {
			switch cfgProfile.Cloud {
			case model.AWS:
				return s.Aws.GetObjectPregisn(profile, region, input)
			case model.TENCENT:
				return s.Tencent.GetObjectPregisn(profile, region, input)
			default:
				return model.ObjectPregisnResponse{}, model.ErrCloudNotSupported
			}
		}
	}
	return model.ObjectPregisnResponse{}, model.ErrProfileNotFound
}
