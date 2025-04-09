package model

type CloudIO interface {
	DescribeInstances(profile, region string, input DescribeInstancesInput) (InstanceResponse, error)
	CreateInstance(profile, region string, input CreateInstanceInput) (CreateInstanceResponse, error)
	ModifyInstance(profile, region string, input ModifyInstanceInput) (ModifyInstanceResponse, error)
	DeleteInstance(profile, region string, input DeleteInstanceInput) (DeleteInstanceResponse, error)
	DescribeKeyPairs(profile, region string, input CommonFilter) ([]KeyPair, error)
	DescribeImages(profile, region string, input CommonFilter) ([]Image, error)
	DescribeInstanceTypes(profile, region string) ([]InstanceType, error)

	// Volume
	DescribeVolumes(profile, region string, input DescribeVolumesInput) ([]Volume, error)
	CreateVolume(profile, region string, input CreateVolumeInput) (string, error)
	ModifyVolume(profile, region string, input ModifyVolumeInput) error
	DeleteVolume(profile, region string, input DeleteVolumeInput) error
	AttachVolume(profile, region string, input AttachVolumeInput) error
	DetachVolume(profile, region string, input DetachVolumeInput) error

	// VPC
	QueryVPC(profile, region string, input CommonFilter) ([]VPC, error)
	QuerySubnet(profile, region string, input CommonFilter) ([]Subnet, error)
	QueryEIP(profile, region string, input CommonFilter) ([]EIP, error)
	QueryNAT(profile, region string, input CommonFilter) ([]NAT, error)
	QuerySecurityGroups(profile, region string, input CommonFilter) ([]SecurityGroup, error)
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
	CreateEmrCluster(profile, region string, input CreateEmrClusterInput) (CreateEmrClusterResponse, error)

	// tencent region is not required
	DescribeDomainList(profile, region string, input DescribeDomainListRequest) (DescribeDomainListResponse, error)
	DescribeRecordList(profile, region string, input DescribeRecordListRequest) (DescribeRecordListResponse, error)
	DescribeRecordListWithPages(profile, region string, input DescribeRecordListWithPageRequest) (ListRecordsPageResponse, error)
	DescribeRecord(profile, region string, input DescribeRecordRequest) (Record, error)
	CreateRecord(profile, region string, input CreateRecordRequest) (CreateRecordResponse, error)
	ModifyRecord(profile, region string, ignoreType bool, input ModifyRecordRequest) error
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

	// S3 COS
	CreateBucket(profile, region string, input CreateBucketRequest) error
	CreateBucketLifecycle(profile, region string, input CreateBucketLifecycleRequest) error
	GetBucketLifecycle(profile, region string, input GetBucketLifecycleRequest) (GetBucketLifecycleResponse, error)
	// DeleteBucketLifecycle(profile, region string, input DeleteBucketLifecycleRequest) error
	DeleteBucket(profile, region string, input DeleteBucketRequest) (DeleteBucketResponse, error)
	ListBucket(profile, region string, input ListBucketRequest) (ListBucketResponse, error) // 比官方多支持了 aws location 返回，并且都带上了tag返回。
	GetObjectPregisn(profile, region string, input ObjectPregisnRequest) (ObjectPregisnResponse, error)
	GetObjectPregisnWithAKSK(ak, sk, region string, input ObjectPregisnRequest) (ObjectPregisnResponse, error) // 支持AKSK的方式获取对象的预签名URL

	// SQS
	CreateSqs(profile, region string, input CreateSqsRequest) error
}
