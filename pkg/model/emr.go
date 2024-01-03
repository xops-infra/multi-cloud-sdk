package model

import (
	"time"
)

// 更详细的信息用Describe接口查询
type EmrCluster struct {
	ID, Name *string
	Status   EMRClusterStatus
}

type EmrFilter struct {
	ClusterStates []EMRClusterStatus
	Period        *time.Duration // 有效期,eg: 1h, 1d, 1w, 1m, 1y
	NextMarker    *string        // 下一页
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
}

type InstanceChargeType string
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
