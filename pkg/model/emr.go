package model

import (
	"fmt"
	"time"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/aws/aws-sdk-go/service/emr"
	txemr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/emr/v20190103"
)

// 更详细的信息用Describe接口查询
type EmrCluster struct {
	ID, Name *string
	Status   EMRClusterStatus
	AddTime  time.Time
}

type EmrFilter struct {
	Profile       *string            `json:"profile" binding:"required"`
	Region        *string            `json:"region" binding:"required"` // 为空则取默认
	ClusterStates []EMRClusterStatus `json:"cluster_states"`
	Period        *time.Duration     `json:"period"` // aws支持的，腾讯云没有用到
	NextMarker    *string            `json:"next_marker"`
}

type DescribeInput struct {
	Profile *string   `json:"profile" binding:"required"`
	Region  *string   `json:"region" binding:"required"` // 为空则取默认
	IDS     []*string `json:"ids"`                       // 为空则取所有
}

type FilterEmrResponse struct {
	NextMarker *string
	Clusters   []EmrCluster
}

type DescribeEmrCluster struct {
	Meta       any // 原始数据
	ID, Name   *string
	Status     EMRClusterStatus
	CreateTime *time.Time
	Tags       Tags
}

type InstanceChargeType string

// *uint64 实例计费模式。取值范围： <li>0：表示按量计费。</li> <li>1：表示包年包月。</li>
func (t InstanceChargeType) ToTencentEmrChargeType() *uint64 {
	if t == PREPAID {
		return tea.Uint64(1)
	}
	return tea.Uint64(0)
}

func (t InstanceChargeType) String() *InstanceChargeType {
	return &t
}

type InstancePolicy string

const (
	POSTPAID_BY_HOUR InstanceChargeType = "POSTPAID_BY_HOUR"
	ON_DEMAND        InstanceChargeType = "ON_DEMAND"
	SPOT             InstanceChargeType = "SPOT"
	PREPAID          InstanceChargeType = "PREPAID"

	Terminate InstancePolicy = "Terminate"
	Reserve   InstancePolicy = "Reserve"
)

type EmrConfigDetail struct {
	KeyId           string `mapstructure:"key_id"`            // 云账号的key_id
	VpcId           string `mapstructure:"vpc_id"`            // vpc_id
	SubnetId        string `mapstructure:"subnet_id"`         // 子网id
	SgId            string `mapstructure:"sg_id" `            // 安全组id
	SgSlave         string `mapstructure:"sg_slave"`          // 安全组id
	SgServiceAccess string `mapstructure:"sg_service_access"` // 安全组id aws 特有私有子网需要
	Zone            string `mapstructure:"zone"`              // 可用区
	LogUri          string `mapstructure:"log_uri"`           // 日志路径 cosn://xxxxx/xxx
	Role            Role   `mapstructure:"role"`              // 角色
}

type Role struct {
	JobFlowRole string `mapstructure:"job_flow_role"` // ec2的role
	ServiceRole string `mapstructure:"service_role"`  // emr服务的role
}

type EmrConfig struct {
	Private        EmrConfigDetail `mapstructure:"private"`
	Public         EmrConfigDetail `mapstructure:"public"`
	EmrTidbPrivate EmrConfigDetail `mapstructure:"tidb-private"`
}

type EmrSubnet string

const (
	EmrSubnetPublic  EmrSubnet = "public"
	EmrSubnetPrivate EmrSubnet = "private"
	EmrTidbPrivate   EmrSubnet = "tidb-private"
)

type EMRClusterStatus string

func EMRClusterStatusList() []EMRClusterStatus {
	return []EMRClusterStatus{
		EMRClusterWaiting,
		EMRClusterStarting,
		EMRClusterBootstrapping,
		EMRClusterRunning,
		EMRClusterTerminated,
		EMRClusterTerminating,
		EMRClusterTerminatedWithErrors,
		EMRClusterUnknown,
	}
}

const (
	EMRClusterWaiting              EMRClusterStatus = "WAITING"
	EMRClusterStarting             EMRClusterStatus = "STARTING"
	EMRClusterBootstrapping        EMRClusterStatus = "BOOTSTRAPPING"
	EMRClusterRunning              EMRClusterStatus = "RUNNING"
	EMRClusterTerminated           EMRClusterStatus = "TERMINATED"
	EMRClusterTerminating          EMRClusterStatus = "TERMINATING"
	EMRClusterTerminatedWithErrors EMRClusterStatus = "TERMINATED_WITH_ERRORS"
	EMRClusterUnknown              EMRClusterStatus = "UNKNOWN"
)

// fmtState 转换 EMR 集群状态
// 接入腾讯云
// 实例的状态码。取值范围： <li>2：表示集群运行中。</li> <li>3：表示集群创建中。</li> <li>4：表示集群扩容中。
// </li> <li>5：表示集群增加router节点中。</li> <li>6：表示集群安装组件中。</li> <li>7：表示集群执行命令中。
// </li> <li>8：表示重启服务中。</li> <li>9：表示进入维护中。</li> <li>10：表示服务暂停中。</li> <li>11：表示退出维护中。
// </li> <li>12：表示退出暂停中。</li> <li>13：表示配置下发中。</li> <li>14：表示销毁集群中。</li> <li>15：表示销毁core节点中。
// </li> <li>16：销毁task节点中。</li> <li>17：表示销毁router节点中。</li> <li>18：表示更改webproxy密码中。</li> <li>19：表示集群隔离中。
// </li> <li>20：表示集群冲正中。</li> <li>21：表示集群回收中。</li> <li>22：表示变配等待中。</li> <li>23：表示集群已隔离。
// </li> <li>24：表示缩容节点中。</li> <li>33：表示集群等待退费中。</li> <li>34：表示集群已退费。</li> <li>301：表示创建失败。
// </li> <li>302：表示扩容失败。</li> 注意：此字段可能返回 null，表示取不到有效值。
// 201 集群隔离中
func FmtTencentState(state *int64) EMRClusterStatus {
	switch *state {
	case 3, 5, 6, 7, 8, 9, 10, 11, 12:
		return EMRClusterStarting
	case 4, 13:
		return EMRClusterBootstrapping
	case 2:
		return EMRClusterRunning
	case 15, 16, 17, 18, 19, 20, 22, 23, 24, 33:
		return EMRClusterRunning
	case 14, 21, 34:
		return EMRClusterTerminated
	case 301, 302:
		return EMRClusterTerminatedWithErrors
	default:
		return EMRClusterUnknown
	}
}

func Contains(states []EMRClusterStatus, state EMRClusterStatus) bool {
	for _, s := range states {
		if s == state {
			return true
		}
	}
	return false
}

type CreateEmrClusterInput struct {
	Name               *string             `json:"name"`
	Tags               Tags                `json:"tags"`
	APPs               []*string           `json:"apps"` // hive、flink、spark
	EMRVersion         *string             `json:"emr_version"`
	InstanceChargeType *InstanceChargeType `json:"instance_charge_type"`
	ResourceSpec       *ResourceSpec       `json:"resource_spec"`
}

func (c *CreateEmrClusterInput) ToAwsRequest() (*emr.RunJobFlowInput, error) {
	return &emr.RunJobFlowInput{}, nil
}

func (c *CreateEmrClusterInput) ToTencentEmrInstanceRequest() (*txemr.CreateInstanceRequest, error) {
	request := txemr.NewCreateInstanceRequest()
	request.ApplicationRole = tea.String("EMR_QCSLinkedRoleInApplicationDataAccess") // EMR_QCSLinkedRoleInApplicationDataAccess
	request.Tags = c.Tags.ToTencentEmrTags()
	request.Software = c.APPs
	request.InstanceName = c.Name
	if c.InstanceChargeType == nil {
		return nil, fmt.Errorf("instance charge type is nil")
	}
	request.PayMode = c.InstanceChargeType.ToTencentEmrChargeType()

	request.SupportHA = tea.Uint64(0)
	if *c.ResourceSpec.HA {
		request.SupportHA = tea.Uint64(1)
	}
	fmt.Println("request.SupportHA", *request.SupportHA)
	request.VPCSettings = &txemr.VPCSettings{
		VpcId:    c.ResourceSpec.VPC,
		SubnetId: c.ResourceSpec.Subnet,
	}
	request.SgId = c.ResourceSpec.SgId
	request.LoginSettings = &txemr.LoginSettings{
		Password:    tea.String("loodai0le!Gh"),
		PublicKeyId: c.ResourceSpec.KeyID,
	}
	request.NeedMasterWan = tea.String("NEED_MASTER_WAN")

	// TODO
	request.ProductId = tea.Uint64(33)
	if c.EMRVersion != nil {
		request.ProductId = tea.Uint64(33)
	}
	request.TimeSpan = tea.Uint64(3600)
	request.TimeUnit = tea.String("s")
	request.MultiZone = tea.Bool(false)
	request.Placement = &txemr.Placement{
		Zone: tea.String("ap-shanghai-5"),
	}

	request.ResourceSpec = &txemr.NewResourceSpec{
		MasterCount: c.ResourceSpec.MasterResourceSpec.InstanceCount,
		MasterResourceSpec: &txemr.Resource{
			InstanceType: c.ResourceSpec.MasterResourceSpec.InstanceType,
			DiskType:     c.ResourceSpec.MasterResourceSpec.DiskType,
			DiskSize:     c.ResourceSpec.MasterResourceSpec.DiskSize,
			DiskNum:      tea.Uint64(1),
			RootSize:     c.ResourceSpec.MasterResourceSpec.RootSize,
			Tags:         c.Tags.ToTencentEmrTags(),
		},
	}
	if request.ResourceSpec.MasterResourceSpec.DiskNum != nil {
		request.ResourceSpec.MasterResourceSpec.DiskNum = tea.Uint64(uint64(tea.Int64Value(c.ResourceSpec.MasterResourceSpec.DiskNum)))
	}

	if c.ResourceSpec.CoreResourceSpec != nil {
		request.ResourceSpec.CoreCount = c.ResourceSpec.CoreResourceSpec.InstanceCount
		request.ResourceSpec.CoreResourceSpec = &txemr.Resource{
			InstanceType: c.ResourceSpec.CoreResourceSpec.InstanceType,
			DiskType:     c.ResourceSpec.CoreResourceSpec.DiskType,
			DiskSize:     c.ResourceSpec.CoreResourceSpec.DiskSize,
			DiskNum:      tea.Uint64(1),
			RootSize:     c.ResourceSpec.CoreResourceSpec.RootSize,
			Tags:         c.Tags.ToTencentEmrTags(),
		}
		if request.ResourceSpec.CoreResourceSpec.DiskNum != nil {
			request.ResourceSpec.CoreResourceSpec.DiskNum = tea.Uint64(uint64(tea.Int64Value(c.ResourceSpec.CoreResourceSpec.DiskNum)))
		}
	}

	if c.ResourceSpec.TaskResourceSpec != nil {
		request.ResourceSpec.TaskCount = c.ResourceSpec.TaskResourceSpec.InstanceCount
		request.ResourceSpec.TaskResourceSpec = &txemr.Resource{
			InstanceType: c.ResourceSpec.TaskResourceSpec.InstanceType,
			DiskType:     c.ResourceSpec.TaskResourceSpec.DiskType,
			DiskSize:     c.ResourceSpec.TaskResourceSpec.DiskSize,
			DiskNum:      tea.Uint64(1),
			RootSize:     c.ResourceSpec.TaskResourceSpec.RootSize,
			Tags:         c.Tags.ToTencentEmrTags(),
		}
		if request.ResourceSpec.TaskResourceSpec.DiskNum != nil {
			request.ResourceSpec.TaskResourceSpec.DiskNum = tea.Uint64(uint64(tea.Int64Value(c.ResourceSpec.TaskResourceSpec.DiskNum)))
		}
	}

	return request, nil
}

type ResourceSpec struct {
	HA                 *bool           `json:"ha"`
	VPC                *string         `json:"vpc"`
	Subnet             *string         `json:"subnet"`
	SgId               *string         `json:"sg_id"`
	Passwd             *string         `json:"passwd"`
	KeyID              *string         `json:"key_id"`
	MasterResourceSpec *EMRInstaceSpec `json:"master_resource_spec"`
	CoreResourceSpec   *EMRInstaceSpec `json:"core_resource_spec"`
	TaskResourceSpec   *EMRInstaceSpec `json:"task_resource_spec"`
}

// tencent https://cloud.tencent.com/document/api/589/33981#Resource
type EMRInstaceSpec struct {
	InstanceCount *int64  `json:"instance_count"`
	InstanceType  *string `json:"instance_type"`
	DiskType      *string `json:"disk_type"`
	DiskSize      *int64  `json:"disk_size"`
	DiskNum       *int64  `json:"disk_num"`
	RootSize      *int64  `json:"root_size"`
}

type CreateEmrClusterResponse struct {
	ID string `json:"id"`
}
