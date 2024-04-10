package model

import (
	"time"
)

type CommonFilter struct {
	ID string `json:"id"`
}

// 按照`ISO8601`标准表示，并且使用`UTC`时间。格式为：`YYYY-MM-DDThh:mm:ssZ` to time.Time
func TimeParse(t string) (time.Time, error) {
	return time.Parse(time.RFC3339, t)
}

// 实现自动区分云的对象接口
type CommonContract interface {
	QueryOcr(profile, region string, input OcrRequest) (OcrResponse, error)
	// tiia
	CreatePicture(profile, region string, input CreatePictureRequest) (CreatePictureResponse, error)
	GetPictureByName(profile, region string, input CommonPictureRequest) (GetPictureByNameResponse, error)
	QueryPicture(profile, region string, input QueryPictureRequest) (QueryPictureResponse, error)
	DeletePicture(profile, region string, input CommonPictureRequest) (CommonPictureResponse, error)
	UpdatePicture(profile, region string, input UpdatePictureRequest) (CommonPictureResponse, error)
	SearchPicture(profile, region string, input SearchPictureRequest) (SearchPictureResponse, error)

	PrivateDomainList(profile string, req DescribeDomainListRequest) (DescribePrivateDomainListResponse, error)
	PrivateRecordList(profile string, req DescribePrivateRecordListRequest) (DescribePrivateRecordListResponse, error)
	PrivateRecordListWithPages(profile string, req DescribePrivateDnsRecordListWithPageRequest) (ListRecordsPageResponse, error)
	PrivateCreateRecord(profile string, req CreateRecordRequest) (CreateRecordResponse, error)
	PrivateModifyRecord(profile string, req ModifyRecordRequest) error
	PrivateDeleteRecord(profile string, req DeletePrivateRecordRequest) error

	DescribeDomainList(profile, region string, req DescribeDomainListRequest) (DescribeDomainListResponse, error)
	DescribeRecordList(profile, region string, req DescribeRecordListRequest) (DescribeRecordListResponse, error)
	DescribeRecordListWithPages(profile, region string, req DescribeRecordListWithPageRequest) (ListRecordsPageResponse, error)
	DescribeRecord(profile, region string, req DescribeRecordRequest) (Record, error)
	CreateRecord(profile, region string, req CreateRecordRequest) (CreateRecordResponse, error)
	ModifyRecord(profile, region string, ignoreType bool, req ModifyRecordRequest) error
	DeleteRecord(profile, region string, req DeleteRecordRequest) (CommonDnsResponse, error)

	DescribeEmrCluster(DescribeInput) ([]DescribeEmrCluster, error)
	QueryEmrCluster(EmrFilter) (FilterEmrResponse, error)

	DescribeInstances(profile, region string, input InstanceFilter) (InstanceResponse, error)
	CreateInstance(profile, region string, input CreateInstanceInput) (CreateInstanceResponse, error)
	ModifyInstance(profile, region string, input ModifyInstanceInput) (ModifyInstanceResponse, error)
	DeleteInstance(profile, region string, input DeleteInstanceInput) (DeleteInstanceResponse, error)

	QueryVPCs(profile, region string, input CommonFilter) ([]VPC, error)
	QuerySubnets(profile, region string, input CommonFilter) ([]Subnet, error)
	QueryEIPs(profile, region string, input CommonFilter) ([]EIP, error)
	QueryNATs(profile, region string, input CommonFilter) ([]NAT, error)
}
