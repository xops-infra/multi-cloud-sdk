package model

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
)

type InstanceContact interface {
	QueryInstances(input InstanceFilter) ([]Instance, error)
}

type Instance struct {
	Name       *string        `json:"name"`
	InstanceID *string        `json:"instance_id" gorm:"primarykey"`
	Profile    string         `json:"profile"`
	KeyName    []*string      `json:"key_name" gorm:"serializer:json"`
	Region     *string        `json:"region"`
	PrivateIP  []*string      `json:"private_ip" gorm:"serializer:json"`
	Platform   *string        `json:"platform"`
	PublicIP   []*string      `json:"public_ip" gorm:"serializer:json"`
	Status     InstanceStatus `json:"status"`
	Owner      *string        `json:"owner"`
	Tags       *Tags          `json:"tags" gorm:"serializer:json"`
}

type InstanceFilter struct {
	Name      *string         `json:"name"`       // 机器名称，使用字符串包含匹配方式
	IDs       []*string       `json:"ids"`        // 机器ID列表，使用字符串包含匹配方式
	Region    *string         `json:"region"`     // 区域
	Profile   *string         `json:"profile"`    // 账号
	PrivateIp *string         `json:"private_ip"` // 私有IP
	PublicIp  *string         `json:"public_ip"`  // 公有IP
	Status    *InstanceStatus `json:"status"`     // 机器状态
	Owner     *string         `json:"owner"`      // 机器所有者，tags的Owner
}

func (q *InstanceFilter) ToTxDescribeInstancesInput() DescribeInstancesInput {
	var instanceIds []*string
	if q.IDs != nil {
		instanceIds = append(instanceIds, q.IDs...)
	}
	var filters []*Filter
	if q.Name != nil {
		filters = append(filters, &Filter{
			Name:   aws.String("instance-name"),
			Values: []*string{q.Name},
		})
	}
	if q.PrivateIp != nil {
		filters = append(filters, &Filter{
			Name:   aws.String("private-ip-address"),
			Values: []*string{q.PrivateIp},
		})
	}
	if q.PublicIp != nil {
		filters = append(filters, &Filter{
			Name:   aws.String("public-ip-address"),
			Values: []*string{q.PublicIp},
		})
	}
	if q.Status != nil {
		filters = append(filters, &Filter{
			Name:   aws.String("instance-state"),
			Values: []*string{aws.String(string(*q.Status))},
		})
	}
	if q.Owner != nil {
		filters = append(filters, &Filter{
			Name:   aws.String("tag:Owner"),
			Values: []*string{q.Owner},
		})
	}
	return DescribeInstancesInput{
		InstanceIds: instanceIds,
		Filters:     filters,
	}
}

func (q *InstanceFilter) ToAwsDescribeInstancesInput() DescribeInstancesInput {
	var instanceIds []*string
	if q.IDs != nil {
		instanceIds = append(instanceIds, q.IDs...)
	}
	var filters []*Filter
	if q.Name != nil {
		filters = append(filters, &Filter{
			Name:   aws.String("tag:Name"),
			Values: []*string{q.Name},
		})
	}
	if q.PrivateIp != nil {
		filters = append(filters, &Filter{
			Name:   aws.String("network-interface.addresses.private-ip-address"),
			Values: []*string{q.PrivateIp},
		})
	}
	if q.PublicIp != nil {
		filters = append(filters, &Filter{
			Name:   aws.String("network-interface.addresses.association.public-ip"),
			Values: []*string{q.PublicIp},
		})
	}
	if q.Status != nil {
		filters = append(filters, &Filter{
			Name:   aws.String("instance-state-name"),
			Values: []*string{aws.String(strings.ToLower(string(*q.Status)))},
		})
	}
	if q.Owner != nil {
		filters = append(filters, &Filter{
			Name:   aws.String("tag:Owner"),
			Values: []*string{q.Owner},
		})
	}
	return DescribeInstancesInput{
		InstanceIds: instanceIds,
		Filters:     filters,
	}
}

// 每次请求的`Filters`的上限为10，`Filter.Values`的上限为5。参数不支持同时指定`InstanceIds`和`Filters`。
type DescribeInstancesInput struct {
	InstanceIds []*string
	Filters     []*Filter
	Size        *int64
	NextMarker  interface{}
}

type Filter struct {
	Name   *string   `type:"string"`
	Values []*string `locationName:"Value" locationNameList:"item" type:"list"`
}

type InstanceResponse struct {
	Instances  []Instance  `json:"instances"`
	NextMarker interface{} `json:"next_marker"` // 如果没有下一页，返回nil 腾讯云直接返回所有数据，不需要分页
}
