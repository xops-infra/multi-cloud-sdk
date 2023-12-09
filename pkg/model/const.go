package model

import "fmt"

type Cloud string

func (i Cloud) ToString() string {
	return string(i)
}

const (
	AWS     Cloud = "aws"
	TENCENT Cloud = "tencent"
)

type InstanceStatus string

func (i InstanceStatus) TString() *InstanceStatus {
	return &i
}

func (i InstanceStatus) ToString() string {
	return string(i)
}

const (
	InstanceStatusRunning InstanceStatus = "RUNNING"
	InstanceStatusStopped InstanceStatus = "STOPPED"
	InstanceStatusPending InstanceStatus = "PENDING"
	InstanceStatusUnknown InstanceStatus = "UNKNOWN"
	// STARTING
	InstanceStatusStarting InstanceStatus = "STARTING"
	// REBOOTING
	InstanceStatusRebooting InstanceStatus = "REBOOTING"
	// TERMINATED
	InstanceStatusTerminated InstanceStatus = "TERMINATED"
	// STOPPING
	InstanceStatusStopping InstanceStatus = "STOPPING"
)

func ToInstanceStatus(s string) InstanceStatus {
	switch s {
	case "RUNNING":
		return InstanceStatusRunning
	case "STOPPED", "SHUTDOWN":
		return InstanceStatusStopped
	case "PENDING":
		return InstanceStatusPending
	case "STARTING":
		return InstanceStatusStarting
	case "REBOOTING":
		return InstanceStatusRebooting
	case "TERMINATED", "TERMINATING":
		return InstanceStatusTerminated
	case "STOPPING", "SHUTTING-DOWN":
		return InstanceStatusStopping
	default:
		fmt.Printf("unknown instance status: [%s]\n", s)
		return InstanceStatusUnknown
	}
}
