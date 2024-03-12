package model

import (
	"github.com/alibabacloud-go/tea/tea"
	privatedns "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/privatedns/v20201028"
)

type DnsContract interface {
	PrivateDomainList(profile string, req DescribeDomainListRequest) (DescribePrivateDomainListResponse, error)
	PrivateRecordList(profile string, req DescribeRecordListRequest) (DescribePrivateRecordListResponse, error)
	PrivateRecord(profile string, req DescribeRecordRequest) (Record, error)
	PrivateCreateRecord(profile string, req CreateRecordRequest) (CreateRecordResponse, error)
	PrivateModifyRecord(profile string, ignoreType bool, req ModifyRecordRequest) error
	PrivateDeleteRecord(profile string, req DeletePrivateRecordRequest) error

	DescribeDomainList(profile string, req DescribeDomainListRequest) (DescribeDomainListResponse, error)
	DescribeRecordList(profile string, req DescribeRecordListRequest) (DescribeRecordListResponse, error)
	DescribeRecord(profile string, req DescribeRecordRequest) (Record, error)
	CreateRecord(profile string, req CreateRecordRequest) (CreateRecordResponse, error)
	ModifyRecord(profile string, ignoreType bool, req ModifyRecordRequest) error
	DeleteRecord(profile string, req DeleteRecordRequest) (CommonDnsResponse, error)
}

type DescribeDomainListRequest struct {
	DomainKeyword *string `json:"keyword"`
}

type DescribeDomainListResponse struct {
	RequestId       *string          `json:"request_id"`
	DomainList      []Domain         `json:"domain_list"`
	DomainCountInfo *DomainCountInfo `json:"domain_count_info"`
}

type Domain struct {
	DomainId *string `json:"domain_id"`
	Name     *string `json:"name"`
	Meta     any     `json:"meta"`
}

type DomainCountInfo struct {
	Total *int64 `json:"total"`
}

type DescribeRecordListRequest struct {
	Domain     *string `json:"domain" binding:"required"`
	RecordType *string `json:"record_type"`
	Keyword    *string `json:"keyword"` // 当前支持搜索主机头和记录值
	Limit      *int64  `json:"limit"`
	NextMarker *string `json:"next_marker"`
}

func (r DescribeRecordListRequest) ToTencentFilter() []*privatedns.Filter {
	var filters []*privatedns.Filter
	if r.RecordType != nil {
		filters = append(filters, &privatedns.Filter{
			Name: tea.String("RecordType"),
			Values: []*string{
				r.RecordType,
			},
		})
	}
	if r.Keyword != nil {
		filters = append(filters, &privatedns.Filter{
			Name: tea.String("SubDomain"),
			Values: []*string{
				r.Keyword,
			},
		})
	}
	return filters
}

type DescribeRecordListResponse struct {
	RecordList []Record `json:"record_list"`
	NextMarker *string  `json:"next_marker"`
}

type DescribeRecordRequest struct {
	Domain     *string `json:"domain" binding:"required"`
	SubDomain  *string `json:"sub_domain" binding:"required"`
	RecordType *string `json:"record_type"`
}

type CreateRecordRequest struct {
	Domain     *string `json:"domain" binding:"required"`
	SubDomain  *string `json:"sub_domain" binding:"required"`  //主机记录，如 www，可选，如果不传，默认为 @。
	RecordType *string `json:"record_type" binding:"required"` //记录类型，通过 API 记录类型获得，大写英文，比如：A 。
	Value      *string `json:"value" binding:"required"`       //记录值，如 IP。
	TTL        *uint64 `json:"ttl"`                            //记录生效时间，默认（aws 300）（腾讯 600），最大值604800秒。
	Info       *string `json:"info"`                           //备注，主要描述创建原因用途（aws不支持，tencent支持，但是private dns 不支持）
}

type CreateRecordResponse struct {
	RecordId *string     `json:"record_id"`
	Meta     interface{} `json:"meta"`
}

type ModifyRecordRequest struct {
	Domain     *string `json:"domain" binding:"required"`
	RecordType *string `json:"record_type" binding:"required"` //记录类型，通过 API 记录类型获得，大写英文，比如：A 。
	Value      *string `json:"value" binding:"required"`       //记录值，如 IP。
	RecordId   *uint64 `json:"record_id"`                      //记录ID。
	SubDomain  *string `json:"sub_domain" binding:"required"`  //主机记录，如 www，可选，如果不传，默认为 @。
	TTL        *uint64 `json:"ttl"`                            //记录生效时间，默认600秒（10分钟），最大值604800秒。
	Weight     *uint64 `json:"weight"`                         //记录权重，值为1-100。
	Status     *bool   `json:"status"`                         //AWS该参数无效。腾讯该参数为是否启用，true 启用，false 禁用。
	Info       *string `json:"info"`                           //备注，主要描述修改原因用途（aws不支持，tencent支持）
}

type Record struct {
	RecordId   *string `json:"record_id"`
	Value      *string `json:"value"` // aws []string 腾讯 string，aws取 1 个可能会有 bug.
	SubDomain  *string `json:"sub_domain"`
	RecordLine *string `json:"record_line"`
	RecordType *string `json:"record_type"`
	TTL        *uint64 `json:"ttl"`
	Status     *string `json:"status"` // ENABLE 和 DISABLE
	UpdatedOn  *string `json:"updated_on"`
	Weight     *uint64 `json:"weight"`
	DomainId   *uint64 `json:"domain_id"`
	Remark     *string `json:"remark"`
	// Meta       interface{} `json:"meta"`
}

type DeleteRecordRequest struct {
	Domain     *string `json:"domain" binding:"required"`
	SubDomain  *string `json:"sub_domain" binding:"required"`
	RecordType *string `json:"record_type" binding:"required"`
}

type CommonDnsResponse struct {
	Meta interface{} `json:"meta"`
}
