package model

type CloudProvider string

func (i CloudProvider) ToString() string {
	return string(i)
}

const (
	AWS     CloudProvider = "aws"
	TENCENT CloudProvider = "tencent"
)

type InstanceStatus string

const (
	InstanceStatusRunning InstanceStatus = "RUNNING"
	InstanceStatusStopped InstanceStatus = "STOPPED"
	InstanceStatusPending InstanceStatus = "PENDING"
	InstanceStatusUnknown InstanceStatus = "UNKNOWN"
)
