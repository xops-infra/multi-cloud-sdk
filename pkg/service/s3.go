package service

import (
	"fmt"

	"github.com/xops-infra/multi-cloud-sdk/pkg/model"
)

func (s *CommonService) CreateBucketLifecycle(profile, region string, input model.CreateBucketLifecycleRequest) error {
	if p, ok := s.Profiles[profile]; ok {
		switch p.Cloud {
		case model.AWS:
			return s.Aws.CreateBucketLifecycle(profile, region, input)
		case model.TENCENT:
			return s.Tencent.CreateBucketLifecycle(profile, region, input)
		default:
			return fmt.Errorf("%s %s", profile, model.ErrCloudNotSupported.Error())
		}
	}
	return fmt.Errorf("%s %s", profile, model.ErrProfileNotFound.Error())
}

func (s *CommonService) GetBucketLifecycle(profile, region string, input model.GetBucketLifecycleRequest) (model.GetBucketLifecycleResponse, error) {
	if p, ok := s.Profiles[profile]; ok {
		switch p.Cloud {
		case model.AWS:
			return s.Aws.GetBucketLifecycle(profile, region, input)
		case model.TENCENT:
			return s.Tencent.GetBucketLifecycle(profile, region, input)
		default:
			return model.GetBucketLifecycleResponse{}, fmt.Errorf("%s %s", profile, model.ErrCloudNotSupported.Error())
		}
	}
	return model.GetBucketLifecycleResponse{}, fmt.Errorf("%s %s", profile, model.ErrProfileNotFound.Error())
}

func (s *CommonService) ListBuckets(profile, region string, input model.ListBucketRequest) (model.ListBucketResponse, error) {
	if p, ok := s.Profiles[profile]; ok {
		switch p.Cloud {
		case model.AWS:
			return s.Aws.ListBucket(profile, region, input)
		case model.TENCENT:
			return s.Tencent.ListBucket(profile, region, input)
		default:
			return model.ListBucketResponse{}, fmt.Errorf("%s %s", profile, model.ErrCloudNotSupported.Error())
		}
	}
	return model.ListBucketResponse{}, fmt.Errorf("%s %s", profile, model.ErrProfileNotFound.Error())
}

func (s *CommonService) CreateBucket(profile, region string, input model.CreateBucketRequest) error {
	if p, ok := s.Profiles[profile]; ok {
		switch p.Cloud {
		case model.AWS:
			return s.Aws.CreateBucket(profile, region, input)
		case model.TENCENT:
			return s.Tencent.CreateBucket(profile, region, input)
		default:
			return fmt.Errorf("%s %s", profile, model.ErrCloudNotSupported.Error())
		}
	}
	return fmt.Errorf("%s %s", profile, model.ErrProfileNotFound.Error())
}

func (s *CommonService) DeleteBucket(profile, region string, input model.DeleteBucketRequest) (model.DeleteBucketResponse, error) {
	if p, ok := s.Profiles[profile]; ok {
		switch p.Cloud {
		case model.AWS:
			return s.Aws.DeleteBucket(profile, region, input)
		case model.TENCENT:
			return s.Tencent.DeleteBucket(profile, region, input)
		default:
			return model.DeleteBucketResponse{}, fmt.Errorf("%s %s", profile, model.ErrCloudNotSupported.Error())
		}
	}
	return model.DeleteBucketResponse{}, fmt.Errorf("%s %s", profile, model.ErrProfileNotFound.Error())
}

func (s *CommonService) GetObjectPregisn(profile, region string, input model.ObjectPregisnRequest) (model.ObjectPregisnResponse, error) {
	if p, ok := s.Profiles[profile]; ok {
		switch p.Cloud {
		case model.AWS:
			return s.Aws.GetObjectPregisn(profile, region, input)
		case model.TENCENT:
			return s.Tencent.GetObjectPregisn(profile, region, input)
		default:
			return model.ObjectPregisnResponse{}, model.ErrCloudNotSupported
		}
	}
	return model.ObjectPregisnResponse{}, fmt.Errorf("%s %s", profile, model.ErrProfileNotFound.Error())
}

func (s *CommonService) GetObjectPregisnWithAKSK(cloud model.Cloud, ak, sk, region string, input model.ObjectPregisnRequest) (model.ObjectPregisnResponse, error) {
	switch cloud {
	case model.AWS:
		return s.Aws.GetObjectPregisnWithAKSK(ak, sk, region, input)
	case model.TENCENT:
		return s.Tencent.GetObjectPregisnWithAKSK(ak, sk, region, input)
	default:
		return model.ObjectPregisnResponse{}, model.ErrCloudNotSupported
	}
}
