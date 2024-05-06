package io

import (
	"fmt"
	"time"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/spf13/cast"
	"github.com/xops-infra/multi-cloud-sdk/pkg/model"

	emr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/emr/v20190103"
)

// EMR 腾讯云因为数据量少所以递归查询所有结果返回
func (c *tencentClient) QueryEmrCluster(input model.EmrFilter) (model.FilterEmrResponse, error) {
	if input.Region == nil && input.Profile == nil {
		return model.FilterEmrResponse{}, fmt.Errorf("region or profile is empty")
	}
	client, err := c.io.GetTencentEmrClient(*input.Profile, *input.Region)
	if err != nil {
		return model.FilterEmrResponse{}, err
	}
	request := emr.NewDescribeInstancesListRequest()
	request.DisplayStrategy = tea.String("clusterList")
	request.Limit = tea.Uint64(100) // TODO: 处理分页，目前不会超过 100 个。
	var clusters []model.EmrCluster
	response, err := client.DescribeInstancesList(request)
	if err != nil {
		return model.FilterEmrResponse{}, err
	}
	for _, cluster := range response.Response.InstancesList {
		state := model.FmtTencentState(tea.Int64(cast.ToInt64(cluster.Status)))
		if len(input.ClusterStates) != 0 {
			if !model.Contains(input.ClusterStates, state) {
				continue
			}
		}
		if cluster.AddTime == nil {
			return model.FilterEmrResponse{}, fmt.Errorf("cluster.AddTime is nil")
		}
		addTime, _ := time.Parse("2006-01-02 15:04:05", *cluster.AddTime)
		clusters = append(clusters, model.EmrCluster{
			ID:      cluster.ClusterId,
			Name:    cluster.ClusterName,
			AddTime: addTime,
			Status:  state,
		})
	}
	return model.FilterEmrResponse{
		Clusters:   clusters,
		NextMarker: nil,
	}, nil
}

func (c *tencentClient) DescribeEmrCluster(input model.DescribeInput) ([]model.DescribeEmrCluster, error) {
	if input.Region == nil && input.Profile == nil {
		return nil, fmt.Errorf("region or profile is empty")
	}
	client, err := c.io.GetTencentEmrClient(*input.Profile, *input.Region)
	if err != nil {
		return nil, err
	}
	request := emr.NewDescribeInstancesRequest()
	request.DisplayStrategy = tea.String("clusterList")
	request.ProjectId = tea.Int64(-1) // 默认-1 查询所有
	request.InstanceIds = input.IDS
	response, err := client.DescribeInstances(request)
	if err != nil {
		return nil, err
	}
	var clusters []model.DescribeEmrCluster
	for _, cluster := range response.Response.ClusterList {
		layout := "2006-01-02 15:04:05"
		createTime, _ := time.Parse(layout, *cluster.AddTime)
		clusters = append(clusters, model.DescribeEmrCluster{
			ID:         cluster.ClusterId,
			Name:       cluster.ClusterName,
			Status:     model.FmtTencentState(cluster.Status),
			CreateTime: &createTime,
			Meta:       cluster,
			Tags:       model.NewTagsFromTencentEmrTags(cluster.Tags),
		})
	}
	return clusters, nil
}

func (c *tencentClient) CreateEmrCluster(profile, region string, input model.CreateEmrClusterInput) (model.CreateEmrClusterResponse, error) {
	client, err := c.io.GetTencentEmrClient(profile, region)
	if err != nil {
		return model.CreateEmrClusterResponse{}, err
	}
	req, err := input.ToTencentEmrInstanceRequest()
	if err != nil {
		return model.CreateEmrClusterResponse{}, err
	}
	fmt.Println(tea.Prettify(req))
	response, err := client.CreateInstance(req)
	if err != nil {
		return model.CreateEmrClusterResponse{}, fmt.Errorf("create emr cluster error: %s", err)
	}
	return model.CreateEmrClusterResponse{ID: *response.Response.InstanceId}, nil
}
