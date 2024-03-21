package model

type CloudIO interface {
	DescribeInstances(profile, region string, input DescribeInstancesInput) (InstanceResponse, error)
	CreateInstance(profile, region string, input CreateInstanceInput) (CreateInstanceResponse, error)
	ModifyInstance(profile, region string, input ModifyInstanceInput) (ModifyInstanceResponse, error)
	DeleteInstance(profile, region string, input DeleteInstanceInput) (DeleteInstanceResponse, error)

	// VPC
	QueryVPC(profile, region string, input CommonFilter) ([]VPC, error)
	QuerySubnet(profile, region string, input CommonFilter) ([]Subnet, error)
	QueryEIP(profile, region string, input CommonFilter) ([]EIP, error)
	QueryNAT(profile, region string, input CommonFilter) ([]NAT, error)
	CreateSecurityGroupWithPolicies(profile, region string, input CreateSecurityGroupWithPoliciesInput) (CreateSecurityGroupWithPoliciesResponse, error) // 创建安全组并添加策略
	CreateSecurityGroupPolicies(profile, region string, input CreateSecurityGroupPoliciesInput) (CreateSecurityGroupPoliciesResponse, error)             // 创建安全组策略

	// Tags
	CreateTags(profile, region string, input CreateTagsInput) error
	AddTagsToResource(profile, region string, input AddTagsInput) error
	RemoveTagsFromResource(profile, region string, input RemoveTagsInput) error
	ModifyTagsForResource(profile, region string, input ModifyTagsInput) error

	// EMR
	QueryEmrCluster(EmrFilter) (FilterEmrResponse, error) // 方便 Post使用，将Profile和Region放入filter
	DescribeEmrCluster(DescribeInput) ([]DescribeEmrCluster, error)

	// tencent region is not required
	DescribeDomainList(profile, region string, input DescribeDomainListRequest) (DescribeDomainListResponse, error)
	// tencent region is not required
	DescribeRecordList(profile, region string, input DescribeRecordListRequest) (DescribeRecordListResponse, error)
	// tencent region is not required
	DescribeRecordListWithPages(profile, region string, input DescribeRecordListWithPageRequest) (ListRecordsPageResponse, error)
	// tencent region is not required
	DescribeRecord(profile, region string, input DescribeRecordRequest) (Record, error)
	// tencent region is not required
	CreateRecord(profile, region string, input CreateRecordRequest) (CreateRecordResponse, error)
	// tencent region is not required
	ModifyRecord(profile, region string, ignoreType bool, input ModifyRecordRequest) error
	// tencent region is not required
	DeleteRecord(profile, region string, input DeleteRecordRequest) (CommonDnsResponse, error)

	// Private_Dns
	DescribePrivateDomainList(profile string, input DescribeDomainListRequest) (DescribePrivateDomainListResponse, error)
	CreatePrivateRecord(profile string, input CreateRecordRequest) (CreateRecordResponse, error)
	DeletePrivateRecord(profile string, input DeletePrivateRecordRequest) error
	ModifyPrivateRecord(profile string, input ModifyRecordRequest) error
	DescribePrivateRecordList(profile string, input DescribePrivateRecordListRequest) (DescribePrivateRecordListResponse, error)
	DescribePrivateRecordListWithPages(profile string, input DescribePrivateDnsRecordListWithPageRequest) (ListRecordsPageResponse, error)

	// OCR
	CommonOCR(profile, region string, input OcrRequest) (OcrResponse, error)
	CreatePicture(profile, region string, input CreatePictureRequest) (CreatePictureResponse, error)
	GetPictureByName(profile, region string, input CommonPictureRequest) (GetPictureByNameResponse, error)
	QueryPicture(profile, region string, input QueryPictureRequest) (QueryPictureResponse, error)
	DeletePicture(profile, region string, input CommonPictureRequest) (CommonPictureResponse, error)
	UpdatePicture(profile, region string, input UpdatePictureRequest) (CommonPictureResponse, error)
	SearchPicture(profile, region string, input SearchPictureRequest) (SearchPictureResponse, error)
}
