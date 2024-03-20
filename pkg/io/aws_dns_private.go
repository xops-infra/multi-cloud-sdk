package io

import "github.com/xops-infra/multi-cloud-sdk/pkg/model"

func (c *awsClient) DescribePrivateDomainList(profile string, input model.DescribeDomainListRequest) (model.DescribePrivateDomainListResponse, error) {
	panic("implement me")
}

func (c *awsClient) DescribePrivateRecordList(profile string, input model.DescribeRecordListRequest) (model.DescribePrivateRecordListResponse, error) {
	panic("implement me")
}

func (c *awsClient) CreatePrivateRecord(profile string, input model.CreateRecordRequest) (model.CreateRecordResponse, error) {
	panic("implement me")
}

func (c *awsClient) ModifyPrivateRecord(profile string, input model.ModifyRecordRequest) error {
	panic("implement me")
}

func (c *awsClient) DeletePrivateRecord(profile string, input model.DeletePrivateRecordRequest) error {
	panic("implement me")
}

func (c *awsClient) DescribePrivateRecordListWithPages(profile string, input model.DescribeRecordListWithPageRequest) (model.ListRecordsPageResponse, error) {
	panic("implement me")
}
