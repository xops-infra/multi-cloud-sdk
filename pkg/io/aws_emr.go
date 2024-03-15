package io

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/emr"
	"github.com/xops-infra/multi-cloud-sdk/pkg/model"
)

// EMR
func (c *awsClient) QueryEmrCluster(filter model.EmrFilter) (model.FilterEmrResponse, error) {
	if filter.Region == nil && filter.Profile == nil {
		return model.FilterEmrResponse{}, fmt.Errorf("region or profile is empty")
	}
	svc, err := c.io.GetAWSEmrClient(*filter.Profile, *filter.Region)
	if err != nil {
		return model.FilterEmrResponse{}, err
	}
	input := &emr.ListClustersInput{}
	if filter.ClusterStates != nil {
		for _, state := range filter.ClusterStates {
			input.ClusterStates = append(input.ClusterStates, aws.String(string(state)))
		}
	}
	if filter.NextMarker != nil {
		input.Marker = filter.NextMarker
	}
	if filter.Period != nil {
		input.CreatedAfter = aws.Time(time.Now().Add(-*filter.Period))
	}
	result, err := svc.ListClusters(input)
	if err != nil {
		return model.FilterEmrResponse{}, err
	}
	var clusters []model.EmrCluster
	for _, cluster := range result.Clusters {
		clusters = append(clusters, model.EmrCluster{
			ID:     cluster.Id,
			Name:   cluster.Name,
			Status: model.EMRClusterStatus(*cluster.Status.State),
		})
	}
	return model.FilterEmrResponse{
		Clusters:   clusters,
		NextMarker: result.Marker,
	}, nil
}

func (c *awsClient) DescribeEmrCluster(input model.DescribeInput) ([]model.DescribeEmrCluster, error) {
	if input.Region == nil && input.Profile == nil {
		return nil, fmt.Errorf("region or profile is empty")
	}
	svc, err := c.io.GetAWSEmrClient(*input.Profile, *input.Region)
	if err != nil {
		return nil, err
	}
	var clusters []model.DescribeEmrCluster
	for _, id := range input.IDS {
		out, err := svc.DescribeCluster(&emr.DescribeClusterInput{
			ClusterId: id,
		})
		if err != nil {
			return nil, err
		}
		clusters = append(clusters, model.DescribeEmrCluster{
			ID:         out.Cluster.Id,
			Name:       out.Cluster.Name,
			Status:     model.EMRClusterStatus(*out.Cluster.Status.State),
			CreateTime: out.Cluster.Status.Timeline.CreationDateTime,
			Meta:       out.Cluster,
		})
	}
	return clusters, nil
}
