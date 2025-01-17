package model

import (
	"strings"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
)

type Instance struct {
	Name       *string        `json:"name"`
	InstanceID *string        `json:"instance_id" gorm:"primarykey"`
	Profile    string         `json:"profile"`
	KeyIDs     []*string      `json:"key_ids" gorm:"serializer:json"`
	Region     *string        `json:"region"`
	PrivateIP  []*string      `json:"private_ip" gorm:"serializer:json"`
	Platform   *string        `json:"platform"`
	PublicIP   []*string      `json:"public_ip" gorm:"serializer:json"`
	Status     InstanceStatus `json:"status"`
	Owner      *string        `json:"owner"`
	Tags       *Tags          `json:"tags" gorm:"serializer:json"`
}

type InstanceFilter struct {
	Name       *string         `json:"name"`        // 机器名称，使用字符串包含匹配方式
	IDs        []*string       `json:"ids"`         // 机器ID列表，使用字符串包含匹配方式
	PrivateIp  *string         `json:"private_ip"`  // 私有IP
	PublicIp   *string         `json:"public_ip"`   // 公有IP
	Status     *InstanceStatus `json:"status"`      // 机器状态
	Owner      *string         `json:"owner"`       // 机器所有者，tags的Owner
	Size       *int64          `json:"size"`        // 分页大小
	NextMarker *string         `json:"next_marker"` // 如果没有下一页，返回nil 腾讯云直接返回所有数据，不需要分页
}

func (q *InstanceFilter) ToTxDescribeInstancesInput() DescribeInstancesInput {
	var instanceIds []*string
	if q.IDs != nil {
		instanceIds = append(instanceIds, q.IDs...)
	}
	var filters []*Filter
	if q.Name != nil {
		filters = append(filters, &Filter{
			Name:   tea.String("instance-name"),
			Values: []*string{q.Name},
		})
	}
	if q.PrivateIp != nil {
		filters = append(filters, &Filter{
			Name:   tea.String("private-ip-address"),
			Values: []*string{q.PrivateIp},
		})
	}
	if q.PublicIp != nil {
		filters = append(filters, &Filter{
			Name:   tea.String("public-ip-address"),
			Values: []*string{q.PublicIp},
		})
	}
	if q.Status != nil {
		filters = append(filters, &Filter{
			Name:   tea.String("instance-state"),
			Values: []*string{tea.String(string(*q.Status))},
		})
	}
	if q.Owner != nil {
		filters = append(filters, &Filter{
			Name:   tea.String("tag:Owner"),
			Values: []*string{q.Owner},
		})
	}
	return DescribeInstancesInput{
		InstanceIds: instanceIds,
		Filters:     filters,
		Size:        q.Size,
		NextMarker:  q.NextMarker,
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
			Name:   tea.String("tag:Name"),
			Values: []*string{q.Name},
		})
	}
	if q.PrivateIp != nil {
		filters = append(filters, &Filter{
			Name:   tea.String("network-interface.addresses.private-ip-address"),
			Values: []*string{q.PrivateIp},
		})
	}
	if q.PublicIp != nil {
		filters = append(filters, &Filter{
			Name:   tea.String("network-interface.addresses.association.public-ip"),
			Values: []*string{q.PublicIp},
		})
	}
	if q.Status != nil {
		filters = append(filters, &Filter{
			Name:   tea.String("instance-state-name"),
			Values: []*string{tea.String(strings.ToLower(string(*q.Status)))},
		})
	}
	if q.Owner != nil {
		filters = append(filters, &Filter{
			Name:   tea.String("tag:Owner"),
			Values: []*string{q.Owner},
		})
	}
	return DescribeInstancesInput{
		InstanceIds: instanceIds,
		Filters:     filters,
		NextMarker:  q.NextMarker,
		Size:        q.Size,
	}
}

// 每次请求的`Filters`的上限为10，`Filter.Values`的上限为5。参数不支持同时指定`InstanceIds`和`Filters`。
type DescribeInstancesInput struct {
	InstanceIds []*string
	Filters     []*Filter
	Size        *int64
	NextMarker  *string
}

type Filter struct {
	Name   *string   `type:"string"`
	Values []*string `locationName:"Value" locationNameList:"item" type:"list"`
}

type InstanceResponse struct {
	Instances  []Instance `json:"instances"`
	NextMarker *string    `json:"next_marker"` // 如果没有下一页，返回nil 腾讯云直接返回所有数据，不需要分页
}

// Create
type CreateInstanceInput struct {
	Name               *string   `json:"name"`
	Count              *int64    `json:"count" default:"1"` // 默认 1 个
	ImageID            *string   `json:"image_id"`
	InstanceType       *string   `json:"instance_type"`
	InstanceChargeType *string   `json:"instance_charge_type"` // 默认按需
	Zone               *string   `json:"zone"`                 // 这里写可用区 ID后台转换
	SystemDisk         *Disk     `json:"system_disk"`
	DataDisks          []Disk    `json:"data_disks"`
	RoleName           *string   `json:"role_name"`
	VpcID              *string   `json:"vpc_id"`
	SecurityGroupIDs   []*string `json:"security_group_ids"`
	SubnetID           *string   `json:"subnet_id"`
	UserData           *string   `json:"user_data"` // base64
	Password           *string   `json:"password"`
	KeyIds             []*string `json:"key_ids"`
	Tags               Tags      `json:"tags"`
}

type Disk struct {
	Size *int64  `json:"size"`
	Type *string `json:"type"` // 硬盘介质类型，腾讯云支持的类型如下：CLOUD_BASIC：表示普通云硬盘 CLOUD_PREMIUM：表示高性能云硬盘 CLOUD_SSD：表示SSD云硬盘 CLOUD_HSSD：表示增强型SSD云硬盘 CLOUD_TSSD：表示极速型SSD云硬盘 CLOUD_BSSD：表示通用型SSD云硬盘
}

// 硬盘介质类型。取值范围：
// CLOUD_BASIC：表示普通云硬盘
// CLOUD_PREMIUM：表示高性能云硬盘
// CLOUD_BSSD：表示通用型SSD云硬盘
// CLOUD_SSD：表示SSD云硬盘
// CLOUD_HSSD：表示增强型SSD云硬盘
// CLOUD_TSSD：表示极速型SSD云硬盘。
type TencenteDiskType string

const (
	TencenteDiskTypeCLOUD_BASIC   TencenteDiskType = "CLOUD_BASIC"
	TencenteDiskTypeCLOUD_PREMIUM TencenteDiskType = "CLOUD_PREMIUM"
	TencenteDiskTypeCLOUD_BSSD    TencenteDiskType = "CLOUD_BSSD"
	TencenteDiskTypeCLOUD_SSD     TencenteDiskType = "CLOUD_SSD"
	TencenteDiskTypeCLOUD_HSSD    TencenteDiskType = "CLOUD_HSSD"
	TencenteDiskTypeCLOUD_TSSD    TencenteDiskType = "CLOUD_TSSD"
)

// set default
func (i *CreateInstanceInput) ToTencentRunInstancesRequest() *cvm.RunInstancesRequest {
	request := cvm.NewRunInstancesRequest()
	request.InstanceChargeType = common.StringPtr("POSTPAID_BY_HOUR")
	if i.InstanceChargeType != nil {
		request.InstanceChargeType = i.InstanceChargeType
	}
	request.InstanceCount = common.Int64Ptr(1)
	if i.Count != nil {
		request.InstanceCount = i.Count
	}
	request.ImageId = i.ImageID
	request.InstanceType = i.InstanceType
	request.Placement = &cvm.Placement{}
	if i.Zone != nil {
		request.Placement.Zone = i.Zone
	}
	request.SystemDisk = &cvm.SystemDisk{
		DiskSize: common.Int64Ptr(40),
	}
	if i.SystemDisk != nil {
		request.SystemDisk.DiskSize = i.SystemDisk.Size
		request.SystemDisk.DiskType = i.SystemDisk.Type
	}
	request.DataDisks = make([]*cvm.DataDisk, 0)
	for _, disk := range i.DataDisks {
		request.DataDisks = append(request.DataDisks, &cvm.DataDisk{
			DiskSize: disk.Size,
		})
	}
	request.CamRoleName = i.RoleName
	request.DisableApiTermination = common.BoolPtr(false) // 默认关闭实例保护
	request.VirtualPrivateCloud = &cvm.VirtualPrivateCloud{
		SubnetId: i.SubnetID,
		VpcId:    i.VpcID,
	}
	request.SecurityGroupIds = i.SecurityGroupIDs
	if i.Name != nil {
		request.InstanceName = i.Name
	}
	request.UserData = common.StringPtr("IyEvYmluL2Jhc2gKZWNobyAiSGVsbG8gTXVsdGlDbG91ZFNkayIK")
	if i.UserData != nil {
		request.UserData = i.UserData
	}
	request.LoginSettings = &cvm.LoginSettings{
		KeyIds:   make([]*string, 0),
		Password: common.StringPtr("MultiCloud@2024"),
	}
	if i.Password != nil {
		request.LoginSettings.Password = i.Password
	}
	if i.KeyIds != nil {
		request.LoginSettings.KeyIds = i.KeyIds
		request.LoginSettings.Password = nil
	}
	if i.Tags != nil {
		request.TagSpecification = i.Tags.ToRunInstanceTags()
	}

	return request
}

type CreateInstanceResponse struct {
	Meta        any       `json:"meta"`
	InstanceIds []*string `json:"instance_ids"`
}

type ModifyInstanceInput struct {
	Action       ModifyAction
	InstanceIDs  []*string `json:"instance_ids"`  // ["ins-r8hr2upy","ins-5d8a23rs"]
	InstanceType *string   `json:"instance_type"` // Action="change_instance_type" 时必填
}

type ModifyAction string

const (
	StartInstance      ModifyAction = "start_instance"
	StopInstance       ModifyAction = "stop_instance"
	RebootInstance     ModifyAction = "reboot_instance"
	ResetInstance      ModifyAction = "reset_instance"
	ChangeInstanceType ModifyAction = "change_instance_type"
	ChangeInstanceTags ModifyAction = "change_instance_tags"
)

type ModifyInstanceResponse struct {
	Meta any `json:"meta"`
}

type DeleteInstanceInput struct {
	InstanceIds []*string `json:"instance_ids"`
	ReleaseDisk *bool     `json:"release_disk"`
}

// 默认不释放
func (i *DeleteInstanceInput) ToTencentTerminateInstancesRequest() *cvm.TerminateInstancesRequest {
	request := cvm.NewTerminateInstancesRequest()
	request.InstanceIds = i.InstanceIds
	request.ReleasePrepaidDataDisks = common.BoolPtr(false)
	if i.ReleaseDisk != nil {
		request.ReleasePrepaidDataDisks = i.ReleaseDisk
	}
	return request
}

type DeleteInstanceResponse struct {
	Meta any `json:"meta"`
}

type SecurityGroup struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type KeyPair struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	PublicKey string `json:"public_ley"`
}

type Image struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Arch     string `json:"arch"`
	Platform string `json:"platform"`
}

type InstanceType struct {
	Type string `json:"type"`
}
