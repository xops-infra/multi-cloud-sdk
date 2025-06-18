package io

import (
	"github.com/xops-infra/multi-cloud-sdk/pkg/model"
)

func (c *tencentClient) GetMonitorMetricData(profile, region string, input model.GetMonitorMetricDataRequest) (*model.GetMonitorMetricDataResponse, error) {
	client, err := c.io.GetTencentMonitorClient(profile, region)
	if err != nil {
		return nil, err
	}
	// 实例化一个请求对象,每个接口都会对应一个request对象
	request, err := input.ToTencentRequest()
	if err != nil {
		return nil, err
	}
	response, err := client.GetMonitorData(request)
	if err != nil {
		return nil, err
	}

	// fmt.Println(tea.Prettify(response))
	metricDatas := make([]model.MetricData, 0)
	for _, data := range response.Response.DataPoints {
		// 处理单个数据
		metricData := model.MetricData{
			Dimensions: make([]model.Dimension, 0),
			DataPoints: make([]model.DataPoint, 0),
		}
		for _, dimension := range data.Dimensions {
			metricData.Dimensions = append(metricData.Dimensions, model.Dimension{
				Name:  *dimension.Name,
				Value: *dimension.Value,
			})
		}
		metricData.DataPoints = append(metricData.DataPoints, model.DataPoint{
			Timestamps: data.Timestamps,
			Values:     data.Values,
		})
		metricDatas = append(metricDatas, metricData)
	}
	return &model.GetMonitorMetricDataResponse{
		MetricName:  input.MetricsType.String(),
		MetricDatas: metricDatas,
		Period:      int64(*response.Response.Period),
	}, nil
}
