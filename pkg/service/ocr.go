package service

import "github.com/xops-infra/multi-cloud-sdk/pkg/model"

type OcrService struct {
	Profiles     []model.ProfileConfig
	Aws, Tencent model.CloudIO
}

func NewOcrService(profiles []model.ProfileConfig, aws, tencent model.CloudIO) model.OcrContract {
	return &OcrService{
		Profiles: profiles,
		Aws:      aws,
		Tencent:  tencent,
	}
}

func (s *OcrService) QueryOcr(profile, region string, request model.OcrRequest) (model.OcrResponse, error) {
	for _, p := range s.Profiles {
		if p.Name == profile {
			switch p.Cloud {
			case model.AWS:
				return s.Aws.CommonOCR(profile, region, request)
			case model.TENCENT:
				return s.Tencent.CommonOCR(profile, region, request)
			default:
				return model.OcrResponse{}, nil
			}
		}
	}
	return model.OcrResponse{}, nil
}

// tiia CreatePicture
func (s *OcrService) CreatePicture(profile, region string, request model.CreatePictureRequest) (model.CreatePictureResponse, error) {
	for _, p := range s.Profiles {
		if p.Name == profile {
			switch p.Cloud {
			case model.AWS:
				return s.Aws.CreatePicture(profile, region, request)
			case model.TENCENT:
				return s.Tencent.CreatePicture(profile, region, request)
			default:
				return model.CreatePictureResponse{}, nil
			}
		}
	}
	return model.CreatePictureResponse{}, nil
}

// tiia GetPictureByName
func (s *OcrService) GetPictureByName(profile, region string, input model.CommonPictureRequest) (model.GetPictureByNameResponse, error) {
	for _, p := range s.Profiles {
		if p.Name == profile {
			switch p.Cloud {
			case model.TENCENT:
				return s.Tencent.GetPictureByName(profile, region, input)
			default:
				return model.GetPictureByNameResponse{}, nil
			}
		}
	}
	return model.GetPictureByNameResponse{}, nil
}

// tiia QueryPicture
func (s *OcrService) QueryPicture(profile, region string, input model.QueryPictureRequest) (model.QueryPictureResponse, error) {
	for _, p := range s.Profiles {
		if p.Name == profile {
			switch p.Cloud {
			case model.TENCENT:
				return s.Tencent.QueryPicture(profile, region, input)
			default:
				return model.QueryPictureResponse{}, nil
			}
		}
	}
	return model.QueryPictureResponse{}, nil
}

// tiia DeletePicture
func (s *OcrService) DeletePicture(profile, region string, input model.CommonPictureRequest) (model.CommonPictureResponse, error) {
	for _, p := range s.Profiles {
		if p.Name == profile {
			switch p.Cloud {
			case model.TENCENT:
				return s.Tencent.DeletePicture(profile, region, input)
			default:
				return model.CommonPictureResponse{}, nil
			}
		}
	}
	return model.CommonPictureResponse{}, nil
}

// tiia UpdatePicture
func (s *OcrService) UpdatePicture(profile, region string, input model.UpdatePictureRequest) (model.CommonPictureResponse, error) {
	for _, p := range s.Profiles {
		if p.Name == profile {
			switch p.Cloud {
			case model.TENCENT:
				return s.Tencent.UpdatePicture(profile, region, input)
			default:
				return model.CommonPictureResponse{}, nil
			}
		}
	}
	return model.CommonPictureResponse{}, nil
}

// tiia SearchPicture
func (s *OcrService) SearchPicture(profile, region string, input model.SearchPictureRequest) (model.SearchPictureResponse, error) {
	for _, p := range s.Profiles {
		if p.Name == profile {
			switch p.Cloud {
			case model.TENCENT:
				return s.Tencent.SearchPicture(profile, region, input)
			default:
				return model.SearchPictureResponse{}, nil
			}
		}
	}
	return model.SearchPictureResponse{}, nil
}
