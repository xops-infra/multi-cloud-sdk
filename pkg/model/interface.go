package model

type CloudIO interface {
	DescribeInstances(profile, region string, input DescribeInstancesInput) (InstanceResponse, error)

	QueryVPC(profile, region string, input CommonFilter) ([]VPC, error)
	QuerySubnet(profile, region string, input CommonFilter) ([]Subnet, error)
	QueryEIP(profile, region string, input CommonFilter) ([]EIP, error)
	QueryNAT(profile, region string, input CommonFilter) ([]NAT, error)

	// EMR
	QueryEmrCluster(profile, region string, input EmrFilter) (FilterEmrResponse, error)
	DescribeEmrCluster(profile, region string, ids []*string) ([]DescribeEmrCluster, error)

	// DNSDomain
	DescribeDomainList(profile, region string, input DescribeDomainListRequest) (DescribeDomainListResponse, error)

	// DNSRecord
	DescribeRecordList(profile, region string, input DescribeRecordListRequest) (DescribeRecordListResponse, error)
	DescribeRecord(profile, region string, input DescribeRecordRequest) (DescribeRecordResponse, error)
	CreateRecord(profile, region string, input CreateRecordRequest) (CreateRecordResponse, error)
	ModifyRecord(profile, region string, ignoreType bool, input ModifyRecordRequest) (ModifyRecordResponse, error)
	DeleteRecord(profile, region string, input DeleteRecordRequest) (CommonDnsResponse, error)

	// OCR
	CommonOCR(profile, region string, input OcrRequest) (OcrResponse, error)
	CreatePicture(profile, region string, input CreatePictureRequest) (CreatePictureResponse, error)
	GetPictureByName(profile, region string, input CommonPictureRequest) (GetPictureByNameResponse, error)
	QueryPicture(profile, region string, input QueryPictureRequest) (QueryPictureResponse, error)
	DeletePicture(profile, region string, input CommonPictureRequest) (CommonPictureResponse, error)
	UpdatePicture(profile, region string, input UpdatePictureRequest) (CommonPictureResponse, error)
	SearchPicture(profile, region string, input SearchPictureRequest) (SearchPictureResponse, error)
}
