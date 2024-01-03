package model

type EmrContact interface {
	DescribeEmrCluster(DescribeInput) ([]DescribeEmrCluster, error)
	QueryEmrCluster(EmrFilter) (FilterEmrResponse, error)
}
