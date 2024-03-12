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

type DescribePrivateRecordListResponse struct {
	TotalCount *int64   `json:"total_count"`
	RecordList []Record `json:"record_list"`
	NextMarker *string  `json:"next_marker"`
}

type DeletePrivateRecordRequest struct {
	Domain    *string   `json:"domain" binding:"required"`
	RecordId  *string   `json:"sub_domain"`
	RecordIds []*string `json:"record_type"`
}
