package io

import (
	"github.com/alibabacloud-go/tea/tea"
	"github.com/spf13/cast"
	ocr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ocr/v20181119"
	tiia "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tiia/v20190529"

	"github.com/xops-infra/multi-cloud-sdk/pkg/model"
)

// CommonOCR
func (c *tencentClient) CommonOCR(profile, region string, input model.OcrRequest) (model.OcrResponse, error) {
	client, err := c.io.GetTencentOcrClient(profile, region)
	if err != nil {
		return model.OcrResponse{}, err
	}
	request := ocr.NewGeneralBasicOCRRequest()
	if input.ImageUrl != nil {
		request.ImageUrl = input.ImageUrl
	}
	if input.ImageBase64 != nil {
		request.ImageBase64 = input.ImageBase64
	}
	if input.LanguageType != nil {
		request.LanguageType = input.LanguageType
	}

	response, err := client.GeneralBasicOCR(request)
	if err != nil {
		return model.OcrResponse{}, err
	}
	var textDetections []model.TextDetection
	for _, textDetection := range response.Response.TextDetections {
		textDetections = append(textDetections, model.TextDetection{
			DetectedText: textDetection.DetectedText,
			Confidence:   textDetection.Confidence,
			Polygon: []*model.Coord{
				{
					X: textDetection.Polygon[0].X,
					Y: textDetection.Polygon[0].Y,
				},
				{
					X: textDetection.Polygon[1].X,
					Y: textDetection.Polygon[1].Y,
				},
				{
					X: textDetection.Polygon[2].X,
					Y: textDetection.Polygon[2].Y,
				},
				{
					X: textDetection.Polygon[3].X,
					Y: textDetection.Polygon[3].Y,
				},
			},
		})
	}
	return model.OcrResponse{
		TextDetections: textDetections,
	}, nil
}

// CreatePicture
func (c *tencentClient) CreatePicture(profile, region string, input model.CreatePictureRequest) (model.CreatePictureResponse, error) {
	client, err := c.io.GetTencentOcrTiiaClient(profile, region)
	if err != nil {
		return model.CreatePictureResponse{}, err
	}
	// 实例化一个请求对象,每个接口都会对应一个request对象
	request := tiia.NewCreateImageRequest()
	request.PicName = input.PicName
	request.EntityId = input.EntityId
	request.GroupId = input.GroupId
	request.ImageUrl = input.ImageUrl
	request.ImageBase64 = input.ImageBase64
	request.Tags = input.Tags

	response, err := client.CreateImage(request)
	if err != nil {
		return model.CreatePictureResponse{}, err
	}
	var object model.Object
	if response.Response.Object != nil {
		object = model.Object{
			Box: &model.Box{
				Rect: &model.ImageRect{
					X:      tea.Int64(cast.ToInt64(response.Response.Object.Box.Rect.X)),
					Y:      tea.Int64(cast.ToInt64(response.Response.Object.Box.Rect.Y)),
					Width:  tea.Int64(cast.ToInt64(response.Response.Object.Box.Rect.Width)),
					Height: tea.Int64(cast.ToInt64(response.Response.Object.Box.Rect.Height)),
				},
			},
			CategoryId: tea.Int64(cast.ToInt64(response.Response.Object.CategoryId)),
			Colors:     tencentColorsToModelColors(response.Response.Object.Colors),
			Attributes: tencentAttrsToModelAttrs(response.Response.Object.Attributes),
			AllBox:     tencentAllBoxToModelAllBox(response.Response.Object.AllBox),
		}
	}
	return model.CreatePictureResponse{
		RequestId: response.Response.RequestId,
		Object:    object,
	}, nil
}

func tencentAllBoxToModelAllBox(allBox []*tiia.Box) []*model.Box {
	var modelAllBox []*model.Box
	for _, box := range allBox {
		modelAllBox = append(modelAllBox, &model.Box{
			Rect: &model.ImageRect{
				X:      tea.Int64(cast.ToInt64(box.Rect.X)),
				Y:      tea.Int64(cast.ToInt64(box.Rect.Y)),
				Width:  tea.Int64(cast.ToInt64(box.Rect.Width)),
				Height: tea.Int64(cast.ToInt64(box.Rect.Height)),
			},
		})
	}
	return modelAllBox
}

func tencentColorsToModelColors(colors []*tiia.ColorInfo) []*model.Color {
	var modelColors []*model.Color
	for _, color := range colors {
		modelColors = append(modelColors, &model.Color{
			Color:      color.Color,
			Percentage: tea.Float64(cast.ToFloat64(color.Percentage)),
			Label:      color.Label,
		})
	}
	return modelColors
}

func tencentAttrsToModelAttrs(attrs []*tiia.Attribute) []*model.Attr {
	var modelAttrs []*model.Attr
	for _, attr := range attrs {
		modelAttrs = append(modelAttrs, &model.Attr{
			Type:    attr.Type,
			Details: attr.Details,
		})
	}
	return modelAttrs
}

// GetPictureByName
func (c *tencentClient) GetPictureByName(profile, region string, input model.CommonPictureRequest) (model.GetPictureByNameResponse, error) {
	client, err := c.io.GetTencentOcrTiiaClient(profile, region)
	if err != nil {
		return model.GetPictureByNameResponse{}, err
	}
	// 实例化一个请求对象,每个接口都会对应一个request对象
	request := tiia.NewDescribeImagesRequest()
	request.EntityId = input.EntityId
	request.GroupId = input.GroupId
	request.PicName = input.PicName
	// 返回的resp是一个DescribeImagesResponse的实例，与请求对象对应
	response, err := client.DescribeImages(request)
	if err != nil {
		return model.GetPictureByNameResponse{}, err
	}
	var imageInfos []model.ImageInfo
	for _, imageInfo := range response.Response.ImageInfos {
		imageInfos = append(imageInfos, model.ImageInfo{
			CustomContent: imageInfo.CustomContent,
			EntityId:      imageInfo.EntityId,
			PicName:       imageInfo.PicName,
			Score:         imageInfo.Score,
			Tags:          imageInfo.Tags,
		})
	}
	return model.GetPictureByNameResponse{
		RequestId:  response.Response.RequestId,
		GroupId:    response.Response.GroupId,
		EntityId:   response.Response.EntityId,
		ImageInfos: imageInfos,
	}, nil
}

// QueryPicture
func (c *tencentClient) QueryPicture(profile, region string, input model.QueryPictureRequest) (model.QueryPictureResponse, error) {
	client, err := c.io.GetTencentOcrTiiaClient(profile, region)
	if err != nil {
		return model.QueryPictureResponse{}, err
	}
	// 实例化一个请求对象,每个接口都会对应一个request对象
	request := tiia.NewDescribeGroupsRequest()

	// 返回的resp是一个DescribeGroupsResponse的实例，与请求对象对应
	response, err := client.DescribeGroups(request)
	if err != nil {
		return model.QueryPictureResponse{}, err
	}
	var groupInfos []model.Group
	for _, groupInfo := range response.Response.Groups {
		groupInfos = append(groupInfos, model.Group{
			GroupId:    groupInfo.GroupId,
			GroupName:  groupInfo.GroupName,
			GroupType:  tea.Int64(cast.ToInt64(groupInfo.GroupType)),
			MaxQps:     tea.Int64(cast.ToInt64(groupInfo.MaxQps)),
			PicCount:   tea.Int64(cast.ToInt64(groupInfo.PicCount)),
			Brief:      groupInfo.Brief,
			CreateTime: groupInfo.CreateTime,
			UpdateTime: groupInfo.UpdateTime,
		})
	}
	return model.QueryPictureResponse{
		RequestId: response.Response.RequestId,
		Groups:    groupInfos,
	}, nil
}

// DeletePicture
func (c *tencentClient) DeletePicture(profile, region string, input model.CommonPictureRequest) (model.CommonPictureResponse, error) {
	client, err := c.io.GetTencentOcrTiiaClient(profile, region)
	if err != nil {
		return model.CommonPictureResponse{}, err
	}
	// 实例化一个请求对象,每个接口都会对应一个request对象
	request := tiia.NewDeleteImagesRequest()
	request.EntityId = input.EntityId
	request.GroupId = input.GroupId
	request.PicName = input.PicName
	// 返回的resp是一个DeleteImageResponse的实例，与请求对象对应
	response, err := client.DeleteImages(request)
	if err != nil {
		return model.CommonPictureResponse{}, err
	}
	return model.CommonPictureResponse{
		RequestId: response.Response.RequestId,
	}, nil
}

// UpdatePicture
func (c *tencentClient) UpdatePicture(profile, region string, input model.UpdatePictureRequest) (model.CommonPictureResponse, error) {
	client, err := c.io.GetTencentOcrTiiaClient(profile, region)
	if err != nil {
		return model.CommonPictureResponse{}, err
	}
	// 实例化一个请求对象,每个接口都会对应一个request对象
	request := tiia.NewUpdateImageRequest()
	request.EntityId = input.EntityId
	request.GroupId = input.GroupId
	request.PicName = input.PicName
	request.Tags = input.Tags
	// 返回的resp是一个ModifyImageResponse的实例，与请求对象对应
	response, err := client.UpdateImage(request)
	if err != nil {
		return model.CommonPictureResponse{}, err
	}
	return model.CommonPictureResponse{
		RequestId: response.Response.RequestId,
	}, nil
}

// SearchPicture
func (c *tencentClient) SearchPicture(profile, region string, input model.SearchPictureRequest) (model.SearchPictureResponse, error) {
	client, err := c.io.GetTencentOcrTiiaClient(profile, region)
	if err != nil {
		return model.SearchPictureResponse{}, err
	}
	// 实例化一个请求对象,每个接口都会对应一个request对象
	request := tiia.NewSearchImageRequest()
	request.GroupId = input.GroupId
	request.ImageUrl = input.ImageUrl
	request.ImageBase64 = input.ImageBase64
	request.Limit = input.Limit
	request.Offset = input.Offset
	request.MatchThreshold = input.MatchThreshold
	request.Filter = input.Filter
	if input.ImageRect != nil {
		request.ImageRect = &tiia.ImageRect{
			X:      input.ImageRect.X,
			Y:      input.ImageRect.Y,
			Width:  input.ImageRect.Width,
			Height: input.ImageRect.Height,
		}
	}
	request.EnableDetect = input.EnableDetect
	request.CategoryId = input.CategoryId
	// 返回的resp是一个SearchImageResponse的实例，与请求对象对应
	response, err := client.SearchImage(request)
	if err != nil {
		return model.SearchPictureResponse{}, err
	}
	var imageInfos []model.ImageInfo
	for _, imageInfo := range response.Response.ImageInfos {
		imageInfos = append(imageInfos, model.ImageInfo{
			CustomContent: imageInfo.CustomContent,
			EntityId:      imageInfo.EntityId,
			PicName:       imageInfo.PicName,
			Score:         imageInfo.Score,
			Tags:          imageInfo.Tags,
		})
	}
	var object model.Object
	if response.Response.Object != nil {
		object = model.Object{
			Box: &model.Box{
				Rect: &model.ImageRect{
					X:      tea.Int64(cast.ToInt64(response.Response.Object.Box.Rect.X)),
					Y:      tea.Int64(cast.ToInt64(response.Response.Object.Box.Rect.Y)),
					Width:  tea.Int64(cast.ToInt64(response.Response.Object.Box.Rect.Width)),
					Height: tea.Int64(cast.ToInt64(response.Response.Object.Box.Rect.Height)),
				},
			},
			CategoryId: tea.Int64(cast.ToInt64(response.Response.Object.CategoryId)),
			Colors:     tencentColorsToModelColors(response.Response.Object.Colors),
			Attributes: tencentAttrsToModelAttrs(response.Response.Object.Attributes),
			AllBox:     tencentAllBoxToModelAllBox(response.Response.Object.AllBox),
		}
	}
	return model.SearchPictureResponse{
		RequestId:  response.Response.RequestId,
		Count:      response.Response.Count,
		ImageInfos: imageInfos,
		Object:     object,
	}, nil
}
