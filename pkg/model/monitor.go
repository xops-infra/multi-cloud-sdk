package model

import (
	"fmt"

	"github.com/alibabacloud-go/tea/tea"
	monitor "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20180724"
)

// MetricInstanceType 监控指标实例类型
type MetricInstanceType int

const (
	// 云服务器实例
	MetricInstanceTypeCVM MetricInstanceType = iota
	// 云数据库实例
	MetricInstanceTypeCDB
	// 负载均衡实例
	MetricInstanceTypeLB
	// 对象存储实例
	MetricInstanceTypeCOS
	// 内容分发网络实例
	MetricInstanceTypeCDN
)

// String 转换为字符串
func (t MetricInstanceType) String() string {
	switch t {
	case MetricInstanceTypeCVM:
		return "CVM"
	case MetricInstanceTypeCDB:
		return "CDB"
	case MetricInstanceTypeLB:
		return "LB"
	case MetricInstanceTypeCOS:
		return "COS"
	case MetricInstanceTypeCDN:
		return "CDN"
	default:
		return "UNKNOWN"
	}
}

// MetricsType 存储类型
type MetricsType int

const (
	// 标准存储
	MetricsTypeStandard MetricsType = iota
	// 低频存储
	MetricsTypeInfrequent
	// 归档存储
	MetricsTypeArchive
	// PUT类请求 QPS
	MetricsTypePutQps
	// GET类请求 QPS
	MetricsTypeGetQps
	// 删除对象请求 QPS DeleteObjectRequestsPs
	MetricsTypeDeleteQps
	// 批量删除对象请求 QPS
	MetricsTypeDeleteMultiObjQps

	// 磁盘使用率
	MetricsTypeCvmDiskUsage
)

// String 转换为字符串
func (t MetricsType) String() string {
	switch t {
	case MetricsTypeStandard:
		return "标准存储用量(MB)"
	case MetricsTypeInfrequent:
		return "低频存储用量(MB)"
	case MetricsTypeArchive:
		return "归档存储用量(MB)"
	case MetricsTypePutQps:
		return "PUT类请求 QPS"
	case MetricsTypeGetQps:
		return "GET类请求 QPS"
	case MetricsTypeDeleteQps:
		return "删除对象请求 QPS"
	case MetricsTypeDeleteMultiObjQps:
		return "批量删除对象请求 QPS"
	case MetricsTypeCvmDiskUsage:
		return "机器磁盘使用率（所有磁盘最大的那个）"
	default:
		return "UNKNOWN"
	}
}

// MetricName 监控指标名称
type MetricName string

const (
	// 对象存储总量
	MetricNameCOSStorage MetricName = "cos_storage"
	// 机器磁盘使用率
	MetricNameDiskUsage MetricName = "disk_usage"
)

/*
只封装几个常用的数据获取，其他需要的可以自己扩展
*/
type GetMonitorMetricDataRequest struct {
	InstanceType MetricInstanceType `json:"instance_type"`
	// 存储类型，仅当 InstanceType 为 MetricInstanceTypeCOS 时有效
	MetricsType *MetricsType `json:"storage_type,omitempty"`

	Instances []string `json:"instances,omitempty"` // cos 的 bucket 列表 一次最多 10个对象
}

// toTencentRequest
func (r *GetMonitorMetricDataRequest) ToTencentRequest() (*monitor.GetMonitorDataRequest, error) {
	request := monitor.NewGetMonitorDataRequest()
	switch r.InstanceType {
	case MetricInstanceTypeCOS:
		request.Namespace = tea.String("QCE/COS")
		if r.MetricsType == nil {
			return nil, fmt.Errorf("storage_type is required")
		}
		/*
			实例对象的维度组合，格式为key-value键值对形式的集合。不同类型的实例字段完全不同，如CVM为[{"Name":"InstanceId","Value":"ins-j0hk02zo"}]，Ckafka为[{"Name":"instanceId","Value":"ckafka-l49k54dd"}]，COS为[{"Name":"appid","Value":"1258344699"},{"Name":"bucket","Value":"rig-1258344699"}]。各个云产品的维度请参阅各个产品监控指标文档，对应的维度列即为维度组合的key，value为key对应的值。单请求最多支持批量拉取10个实例的监控数据。
		*/
		if len(r.Instances) > 0 {
			instances := make([]*monitor.Instance, 0)
			for _, instance := range r.Instances {
				instances = append(instances, &monitor.Instance{
					Dimensions: []*monitor.Dimension{
						{
							Name:  tea.String("bucket"),
							Value: tea.String(instance),
						},
					},
				})
			}
			request.Instances = instances
		} else {
			return nil, fmt.Errorf("instances is required")
		}

		// request.SpecifyStatistics = tea.Int64(2) // avg, max, min (1,2,4)可以自由组合

		switch *r.MetricsType {
		case MetricsTypeStandard:
			request.MetricName = tea.String("StdStorage") // 标准存储 单位MB
		case MetricsTypeInfrequent:
			request.MetricName = tea.String("InfrequentStorage") // 低频存储 单位MB
		case MetricsTypeArchive:
			request.MetricName = tea.String("ArcStorage") // 归档存储 单位MB
		case MetricsTypePutQps:
			request.MetricName = tea.String("PutRequestsPs") // PUT类请求 QPS
		case MetricsTypeGetQps:
			request.MetricName = tea.String("GetRequestsPs") // GET类请求 QPS
		case MetricsTypeDeleteQps:
			request.MetricName = tea.String("DeleteObjectRequestsPs") // 删除对象请求 QPS
		case MetricsTypeDeleteMultiObjQps:
			request.MetricName = tea.String("DeleteMultiObjectRequestsPs") // 批量删除对象请求 QPS
		default:
			return nil, fmt.Errorf("metrics_type is not supported")
		}
	case MetricInstanceTypeCVM:
		// https://cloud.tencent.com/document/product/248/6843#.E7.A3.81.E7.9B.98.E7.9B.91.E6.8E.A7
		request.Namespace = tea.String("QCE/CVM")
		if r.MetricsType == nil {
			return nil, fmt.Errorf("metrics_type is required")
		}
		request.Instances = make([]*monitor.Instance, 0)
		for _, instance := range r.Instances {
			request.Instances = append(request.Instances, &monitor.Instance{
				Dimensions: []*monitor.Dimension{
					{
						Name:  tea.String("InstanceId"),
						Value: tea.String(instance),
					},
				},
			})
		}
		switch *r.MetricsType {
		case MetricsTypeCvmDiskUsage:
			// 磁盘已使用容量占总容量的百分比（所有磁盘中最大值）
			request.MetricName = tea.String("CvmDiskUsage")
		default:
			return nil, fmt.Errorf("metrics_type is not supported")
		}
	default:
		return nil, fmt.Errorf("instance_type is not supported")
	}
	return request, nil
}

type GetMonitorMetricDataResponse struct {
	MetricName  string       `json:"metric_name"`
	MetricDatas []MetricData `json:"metric_datas"`
	Period      int64        `json:"period"`
}

type MetricData struct {
	Dimensions []Dimension `json:"dimensions"`
	DataPoints []DataPoint `json:"data_points"`
}

type Dimension struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type DataPoint struct {
	// 时间戳
	Timestamps []*float64 `json:"timestamp"`
	// 指标值
	Values []*float64 `json:"value"`
}
