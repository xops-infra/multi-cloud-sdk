package model

type DescribePrivateDomainListResponse struct {
	DomainList []PrivateDomain `json:"domain_list"`
	TotalCount *int64          `json:"total_count"`
}

type PrivateDomain struct {
	DomainId    *string `json:"domain_id"`
	Name        *string `json:"name"`
	RecordCount *int64  `json:"record_count"`
	VpcSet      any     `json:"vpc_set"`
	Status      *string `json:"status"` // 私有域绑定VPC状态，未关联vpc：SUSPEND，已关联VPC：ENABLED，关联VPC失败：FAILED
	Tags        any     `json:"tags"`
}
type DescribePrivateDnsRecordListWithPageRequest struct {
	Domain *string `json:"domain" binding:"required"` // 支持使用域名或者ID
	Limit  *int64  `json:"limit"`                     // 分页 默认100
	Page   *int64  `json:"page"`                      // 页码
}

type DescribePrivateRecordListRequest struct {
	Domain  *string `json:"domain" binding:"required"`
	Keyword *string `json:"keyword"` // 只支持二级域名的模糊搜索
}

type DescribePrivateRecordListResponse struct {
	TotalCount *int64   `json:"total_count"`
	RecordList []Record `json:"record_list"`
}

type DeletePrivateRecordRequest struct {
	Domain    *string   `json:"domain" binding:"required"`
	RecordId  *string   `json:"record_id" `
	RecordIds []*string `json:"record_ids"`
}
