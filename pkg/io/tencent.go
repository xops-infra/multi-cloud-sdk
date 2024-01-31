package io

import (
	"fmt"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/spf13/cast"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	dnspod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323"
	ocr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ocr/v20181119"
	tiia "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tiia/v20190529"
	tencentVpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	"github.com/xops-infra/multi-cloud-sdk/pkg/model"
)

type tencentClient struct {
	io model.ClientIo
}

func NewTencentClient(io model.ClientIo) model.CloudIO {
	return &tencentClient{
		io: io,
	}
}

func (c *tencentClient) QueryVPC(profile, region string, input model.CommonFilter) ([]model.VPC, error) {
	client, err := c.io.GetTencentVpcClient(profile, region)
	if err != nil {
		return nil, err
	}
	request := tencentVpc.NewDescribeVpcsRequest()
	if input.ID != "" {
		request.Filters = []*tencentVpc.Filter{
			{
				Name:   common.StringPtr("vpc-id"),
				Values: []*string{common.StringPtr(input.ID)},
			},
		}
	}
	response, err := client.DescribeVpcs(request)
	if err != nil {
		return nil, err
	}
	var vpcs []model.VPC
	for _, vpc := range response.Response.VpcSet {
		vpcs = append(vpcs, model.VPC{
			ID:            *vpc.VpcId,
			Region:        region,
			Account:       profile,
			CloudProvider: model.TENCENT,
			Tags:          model.TencentVpcTagsFmt(vpc.TagSet),
			IsDefault:     *vpc.IsDefault,
			CidrBlock:     *vpc.CidrBlock,
		})
	}

	return vpcs, nil
}

// QuerySubnet
func (c *tencentClient) QuerySubnet(profile, region string, input model.CommonFilter) ([]model.Subnet, error) {
	client, err := c.io.GetTencentVpcClient(profile, region)
	if err != nil {
		return nil, err
	}
	request := tencentVpc.NewDescribeSubnetsRequest()
	if input.ID != "" {
		request.Filters = []*tencentVpc.Filter{
			{
				Name:   common.StringPtr("subnet-id"),
				Values: []*string{common.StringPtr(input.ID)},
			},
		}
	}
	response, err := client.DescribeSubnets(request)
	if err != nil {
		return nil, err
	}
	var subnets []model.Subnet
	for _, subnet := range response.Response.SubnetSet {
		createTime, _ := model.TimeParse(*subnet.CreatedTime)
		subnets = append(subnets, model.Subnet{
			ID:                      subnet.SubnetId,
			Region:                  region,
			Account:                 profile,
			CloudProvider:           model.TENCENT,
			Tags:                    model.TencentVpcTagsFmt(subnet.TagSet),
			VpcID:                   subnet.VpcId,
			Name:                    subnet.SubnetName,
			CidrBlock:               subnet.CidrBlock,
			IsDefault:               subnet.IsDefault,
			Zone:                    subnet.Zone,
			RouteTableId:            subnet.RouteTableId,
			CreatedTime:             &createTime,
			AvailableIpAddressCount: cast.ToInt64(subnet.AvailableIpAddressCount),
			NetworkAclId:            subnet.NetworkAclId,
		})
	}
	return subnets, nil
}

// QueryEIP
func (c *tencentClient) QueryEIP(profile, region string, input model.CommonFilter) ([]model.EIP, error) {
	client, err := c.io.GetTencentVpcClient(profile, region)
	if err != nil {
		return nil, err
	}
	request := tencentVpc.NewDescribeAddressesRequest()
	if input.ID != "" {
		request.Filters = []*tencentVpc.Filter{
			{
				Name:   common.StringPtr("address-id"),
				Values: []*string{common.StringPtr(input.ID)},
			},
		}
	}
	response, err := client.DescribeAddresses(request)
	if err != nil {
		return nil, err
	}
	var eips []model.EIP
	for _, eip := range response.Response.AddressSet {
		createTime, _ := model.TimeParse(*eip.CreatedTime)
		eips = append(eips, model.EIP{
			ID:                 eip.AddressId,
			Region:             region,
			Account:            profile,
			CloudProvider:      model.TENCENT,
			Tags:               model.TencentVpcTagsFmt(eip.TagSet),
			Name:               eip.AddressName,
			Status:             eip.AddressStatus,
			AddressIp:          eip.AddressIp,
			InstanceId:         eip.InstanceId,
			CreatedTime:        &createTime,
			Bandwidth:          tea.Int64(cast.ToInt64(eip.Bandwidth)),
			InternetChargeType: eip.InternetChargeType,
		})
	}
	return eips, nil
}

// QueryNAT
func (c *tencentClient) QueryNAT(profile, region string, input model.CommonFilter) ([]model.NAT, error) {
	client, err := c.io.GetTencentVpcClient(profile, region)
	if err != nil {
		return nil, err
	}
	request := tencentVpc.NewDescribeNatGatewaysRequest()
	if input.ID != "" {
		request.Filters = []*tencentVpc.Filter{
			{
				Name:   common.StringPtr("nat-gateway-id"),
				Values: []*string{common.StringPtr(input.ID)},
			},
		}
	}
	response, err := client.DescribeNatGateways(request)
	if err != nil {
		return nil, err
	}
	var nats []model.NAT
	for _, nat := range response.Response.NatGatewaySet {
		createTime, _ := model.TimeParse(*nat.CreatedTime)
		nats = append(nats, model.NAT{
			ID:            *nat.NatGatewayId,
			Region:        region,
			Account:       profile,
			CloudProvider: model.TENCENT,
			Tags:          model.TencentVpcTagsFmt(nat.TagSet),
			Name:          *nat.NatGatewayName,
			Status:        *nat.State,
			VpcID:         *nat.VpcId,
			Zone:          nat.Zone,
			SubnetID:      *nat.SubnetId,
			CreatedTime:   createTime,
		})
	}
	return nats, nil
}

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

// DescribeDomainList
func (c *tencentClient) DescribeDomainList(profile, region string, input model.DescribeDomainListRequest) (model.DescribeDomainListResponse, error) {
	client, err := c.io.GetTencentDnsPodClient(profile, region)
	if err != nil {
		return model.DescribeDomainListResponse{}, err
	}
	request := dnspod.NewDescribeDomainListRequest()
	request.Type = tea.String("ALL")
	request.Keyword = input.DomainKeyword

	response, err := client.DescribeDomainList(request)
	if err != nil {
		return model.DescribeDomainListResponse{}, err
	}
	var domains []model.Domain
	for _, domain := range response.Response.DomainList {
		domains = append(domains, model.Domain{
			DomainId: tea.String(cast.ToString(domain.DomainId)),
			Name:     domain.Name,
			Meta:     domain,
		})

	}
	return model.DescribeDomainListResponse{
		RequestId:  response.Response.RequestId,
		DomainList: domains,
		DomainCountInfo: &model.DomainCountInfo{
			Total: tea.Int64(cast.ToInt64(response.Response.DomainCountInfo.AllTotal)),
		},
	}, nil
}

// DescribeRecordList
func (c *tencentClient) DescribeRecordList(profile, region string, input model.DescribeRecordListRequest) (model.DescribeRecordListResponse, error) {
	client, err := c.io.GetTencentDnsPodClient(profile, region)
	if err != nil {
		return model.DescribeRecordListResponse{}, err
	}
	request := dnspod.NewDescribeRecordListRequest()
	request.Domain = input.Domain
	// request.Subdomain = input.Keyword
	request.RecordType = input.RecordType
	request.Keyword = input.Keyword

	response, err := client.DescribeRecordList(request)
	if err != nil {
		return model.DescribeRecordListResponse{}, err
	}
	var records []model.Record
	for _, record := range response.Response.RecordList {
		records = append(records, model.Record{
			RecordId:   record.RecordId,
			SubDomain:  record.Name,
			RecordType: record.Type,
			Value:      record.Value,
			Status:     record.Status,
			UpdatedOn:  record.UpdatedOn,
			TTL:        record.TTL,
			RecordLine: record.Line,
			Remark:     record.Remark,
			Weight:     record.Weight,
			Meta:       record,
		})
	}
	return model.DescribeRecordListResponse{
		RequestId:  response.Response.RequestId,
		RecordList: records,
		RecordCountInfo: &model.RecordCountInfo{
			Total: tea.Int64(cast.ToInt64(response.Response.RecordCountInfo.TotalCount)),
		},
	}, nil
}

// DescribeRecord
func (c *tencentClient) DescribeRecord(profile, region string, input model.DescribeRecordRequest) (model.DescribeRecordResponse, error) {
	resp, err := c.DescribeRecordList(profile, region, model.DescribeRecordListRequest{
		Domain:  input.Domain,
		Keyword: input.SubDomain,
	})
	if err != nil {
		return model.DescribeRecordResponse{}, err
	}
	for _, record := range resp.RecordList {
		if *record.SubDomain == *input.SubDomain {
			if input.RecordType != nil && *input.RecordType != "" && *input.RecordType != *record.RecordType {
				continue
			}
			return model.DescribeRecordResponse{
				Record: record,
			}, nil
		}
	}
	return model.DescribeRecordResponse{}, fmt.Errorf("record not found")
}

// CreateRecord
func (c *tencentClient) CreateRecord(profile, region string, input model.CreateRecordRequest) (model.CreateRecordResponse, error) {
	client, err := c.io.GetTencentDnsPodClient(profile, region)
	if err != nil {
		return model.CreateRecordResponse{}, err
	}
	// 实例化一个请求对象,每个接口都会对应一个request对象
	request := dnspod.NewCreateRecordRequest()
	request.Domain = input.Domain
	request.SubDomain = input.SubDomain
	request.RecordType = input.RecordType
	request.Value = input.Value
	request.RecordLine = tea.String("默认")
	request.Remark = input.Info
	if input.TTL != nil {
		request.TTL = input.TTL
	} else {
		request.TTL = tea.Uint64(600)
	}
	// 返回的resp是一个CreatePrivateZoneRecordResponse的实例，与请求对象对应
	response, err := client.CreateRecord(request)
	if err != nil {
		return model.CreateRecordResponse{}, err
	}
	return model.CreateRecordResponse{
		RecordId: tea.String(cast.ToString(response.Response.RecordId)),
		Meta:     response.Response,
	}, nil
}

// ModifyRecord
// ignoreType 是否开启忽略 recordType,
// true 注意这里会删除所有相同 subDomain 的记录，然后创建新的记录
// false 如果 recordType 不同，会报没找到记录
func (c *tencentClient) ModifyRecord(profile, region string, ignoreType bool, input model.ModifyRecordRequest) (model.ModifyRecordResponse, error) {
	client, err := c.io.GetTencentDnsPodClient(profile, region)
	if err != nil {
		return model.ModifyRecordResponse{}, err
	}
	if ignoreType {
		resp, err := c.DescribeRecordList(profile, region, model.DescribeRecordListRequest{
			Domain:  input.Domain,
			Keyword: input.SubDomain,
		})
		if err != nil {
			return model.ModifyRecordResponse{}, err
		}
		var delDomain []map[string]interface{}
		for _, record := range resp.RecordList {
			if *record.SubDomain == *input.SubDomain {
				_, err := c.DeleteRecord(profile, region, model.DeleteRecordRequest{
					Domain:     input.Domain,
					SubDomain:  input.SubDomain,
					RecordType: record.RecordType,
				})
				if err != nil {
					return model.ModifyRecordResponse{}, fmt.Errorf("delete record error: %v", err)
				}
				delDomain = append(delDomain, map[string]interface{}{
					"recordId":   record.RecordId,
					"recordType": record.RecordType,
					"subDomain":  record.SubDomain,
					"ttl":        record.TTL,
					"value":      record.Value,
				})
			}
		}

		if delDomain == nil {
			return model.ModifyRecordResponse{}, fmt.Errorf("record not found")
		}

		createInput := model.CreateRecordRequest{
			Domain:     input.Domain,
			SubDomain:  input.SubDomain,
			RecordType: input.RecordType,
			Value:      input.Value,
			TTL:        tea.Uint64(600),
			Info:       input.Info,
		}
		if input.TTL != nil {
			createInput.TTL = input.TTL
		}
		createResp, err := c.CreateRecord(profile, region, createInput)
		if err != nil {
			return model.ModifyRecordResponse{}, fmt.Errorf("create record error: %v", err)
		}
		return model.ModifyRecordResponse{
			RecordId: createResp.RecordId,
			Meta: map[string]interface{}{
				"info":      "enable ignoreType",
				"delDomain": delDomain,
			},
		}, nil
	} else {
		recordId, err := c.getRecordIdBySubDomain(profile, region, *input.SubDomain, *input.Domain, *input.RecordType)
		if err != nil {
			return model.ModifyRecordResponse{}, err
		}

		request := dnspod.NewModifyRecordRequest()
		request.RecordId = recordId
		request.Domain = input.Domain
		request.SubDomain = input.SubDomain
		request.RecordType = input.RecordType
		request.Value = input.Value
		request.RecordLine = tea.String("默认")
		request.TTL = input.TTL
		request.Weight = input.Weight

		if input.Status != nil {
			if *input.Status {
				request.Status = tea.String("ENABLE")
			} else {
				request.Status = tea.String("DISABLE")
			}
		}

		response, err := client.ModifyRecord(request)
		if err != nil {
			return model.ModifyRecordResponse{}, err
		}
		return model.ModifyRecordResponse{
			RecordId: tea.String(cast.ToString(response.Response.RecordId)),
			Meta:     response.Response,
		}, nil
	}
}

// getRecordIdBySubDomain
func (c *tencentClient) getRecordIdBySubDomain(profile, region, subDomain, domain, recordType string) (*uint64, error) {
	resp, err := c.DescribeRecordList(profile, region, model.DescribeRecordListRequest{
		Domain:  &domain,
		Keyword: &subDomain,
	})
	if err != nil {
		return nil, err
	}
	for _, record := range resp.RecordList {
		if *record.SubDomain == subDomain && *record.RecordType == recordType {
			return record.RecordId, nil
		}
	}
	return nil, fmt.Errorf("record not found")
}

// DeleteRecord
func (c *tencentClient) DeleteRecord(profile, region string, input model.DeleteRecordRequest) (model.CommonDnsResponse, error) {
	client, err := c.io.GetTencentDnsPodClient(profile, region)
	if err != nil {
		return model.CommonDnsResponse{}, err
	}
	record_id, err := c.getRecordIdBySubDomain(profile, region, *input.SubDomain, *input.Domain, *input.RecordType)
	if err != nil {
		return model.CommonDnsResponse{}, err
	}

	request := dnspod.NewDeleteRecordRequest()
	request.RecordId = record_id
	request.Domain = input.Domain

	resp, err := client.DeleteRecord(request)
	if err != nil {
		return model.CommonDnsResponse{}, err
	}
	return model.CommonDnsResponse{
		Meta: resp.Response,
	}, nil
}
