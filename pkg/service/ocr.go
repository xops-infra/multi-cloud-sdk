package service

import (
	"fmt"

	"github.com/xops-infra/multi-cloud-sdk/pkg/model"
)

func (s *CommonService) QueryOcr(profile, region string, request model.OcrRequest) (model.OcrResponse, error) {
	if p, ok := s.Profiles[profile]; ok {
		switch p.Cloud {
		case model.AWS:
			return s.Aws.CommonOCR(profile, region, request)
		case model.TENCENT:
			return s.Tencent.CommonOCR(profile, region, request)
		default:
			return model.OcrResponse{}, fmt.Errorf("%s %s", profile, model.ErrCloudNotSupported.Error())
		}
	}
	return model.OcrResponse{}, fmt.Errorf("%s %s", profile, model.ErrProfileNotFound.Error())
}

// tiia CreatePicture
func (s *CommonService) CreatePicture(profile, region string, request model.CreatePictureRequest) (model.CreatePictureResponse, error) {
	if p, ok := s.Profiles[profile]; ok {
		switch p.Cloud {
		case model.AWS:
			return s.Aws.CreatePicture(profile, region, request)
		case model.TENCENT:
			return s.Tencent.CreatePicture(profile, region, request)
		default:
			return model.CreatePictureResponse{}, fmt.Errorf("%s %s", profile, model.ErrCloudNotSupported.Error())
		}
	}
	return model.CreatePictureResponse{}, fmt.Errorf("%s %s", profile, model.ErrProfileNotFound.Error())
}

// tiia GetPictureByName
func (s *CommonService) GetPictureByName(profile, region string, input model.CommonPictureRequest) (model.GetPictureByNameResponse, error) {
	if p, ok := s.Profiles[profile]; ok {
		switch p.Cloud {
		case model.TENCENT:
			return s.Tencent.GetPictureByName(profile, region, input)
		default:
			return model.GetPictureByNameResponse{}, fmt.Errorf("%s %s", profile, model.ErrCloudNotSupported.Error())
		}
	}
	return model.GetPictureByNameResponse{}, fmt.Errorf("%s %s", profile, model.ErrProfileNotFound.Error())
}

// tiia QueryPicture
func (s *CommonService) QueryPicture(profile, region string, input model.QueryPictureRequest) (model.QueryPictureResponse, error) {
	if p, ok := s.Profiles[profile]; ok {
		switch p.Cloud {
		case model.TENCENT:
			return s.Tencent.QueryPicture(profile, region, input)
		default:
			return model.QueryPictureResponse{}, fmt.Errorf("%s %s", profile, model.ErrCloudNotSupported.Error())
		}
	}
	return model.QueryPictureResponse{}, fmt.Errorf("%s %s", profile, model.ErrProfileNotFound.Error())
}

// tiia DeletePicture
func (s *CommonService) DeletePicture(profile, region string, input model.CommonPictureRequest) (model.CommonPictureResponse, error) {
	if p, ok := s.Profiles[profile]; ok {
		switch p.Cloud {
		case model.TENCENT:
			return s.Tencent.DeletePicture(profile, region, input)
		default:
			return model.CommonPictureResponse{}, fmt.Errorf("%s %s", profile, model.ErrCloudNotSupported.Error())
		}
	}
	return model.CommonPictureResponse{}, fmt.Errorf("%s %s", profile, model.ErrProfileNotFound.Error())
}

// tiia UpdatePicture
func (s *CommonService) UpdatePicture(profile, region string, input model.UpdatePictureRequest) (model.CommonPictureResponse, error) {
	if p, ok := s.Profiles[profile]; ok {
		switch p.Cloud {
		case model.TENCENT:
			return s.Tencent.UpdatePicture(profile, region, input)
		default:
			return model.CommonPictureResponse{}, fmt.Errorf("%s %s", profile, model.ErrCloudNotSupported.Error())
		}
	}
	return model.CommonPictureResponse{}, fmt.Errorf("%s %s", profile, model.ErrProfileNotFound.Error())
}

// tiia SearchPicture
func (s *CommonService) SearchPicture(profile, region string, input model.SearchPictureRequest) (model.SearchPictureResponse, error) {
	if p, ok := s.Profiles[profile]; ok {
		switch p.Cloud {
		case model.TENCENT:
			return s.Tencent.SearchPicture(profile, region, input)
		default:
			return model.SearchPictureResponse{}, fmt.Errorf("%s %s", profile, model.ErrCloudNotSupported.Error())
		}
	}
	return model.SearchPictureResponse{}, fmt.Errorf("%s %s", profile, model.ErrProfileNotFound.Error())
}
