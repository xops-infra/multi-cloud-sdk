package model

type EmrContact interface {
	DescribeEmrCluster(profile, region string, ids []*string) ([]DescribeEmrCluster, error)
	QueryEmrCluster(profile, region string, input EmrFilter) (FilterEmrResponse, error)
}
