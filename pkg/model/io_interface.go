package model

type CloudIo interface {
	QueryInstances(profile, region string) ([]*Instance, error)

	QueryVPC(profile, region string, input CommonQueryInput) ([]*VPC, error)
	QuerySubnet(profile, region string, input CommonQueryInput) ([]*Subnet, error)
	QueryEIP(profile, region string, input CommonQueryInput) ([]*EIP, error)
	QueryNAT(profile, region string, input CommonQueryInput) ([]*NAT, error)

	// domain
	DescribeDomainList(profile, region string, input DescribeDomainListRequest) (DescribeDomainListResponse, error)

	// record
	DescribeRecordList(profile, region string, input DescribeRecordListRequest) (DescribeRecordListResponse, error)
	DescribeRecord(profile, region string, input DescribeRecordRequest) (DescribeRecordResponse, error)
	CreateRecord(profile, region string, input CreateRecordRequest) (CreateRecordResponse, error)
	ModifyRecord(profile, region string, ignoreType bool, input ModifyRecordRequest) (ModifyRecordResponse, error)
	DeleteRecord(profile, region string, input DeleteRecordRequest) (CommonDnsResponse, error)

	CommonOCR(profile, region string, input OcrRequest) (OcrResponse, error)
	CreatePicture(profile, region string, input CreatePictureRequest) (CreatePictureResponse, error)
	GetPictureByName(profile, region string, input CommonPictureRequest) (GetPictureByNameResponse, error)
	QueryPicture(profile, region string, input QueryPictureRequest) (QueryPictureResponse, error)
	DeletePicture(profile, region string, input CommonPictureRequest) (CommonPictureResponse, error)
	UpdatePicture(profile, region string, input UpdatePictureRequest) (CommonPictureResponse, error)
	SearchPicture(profile, region string, input SearchPictureRequest) (SearchPictureResponse, error)
}
