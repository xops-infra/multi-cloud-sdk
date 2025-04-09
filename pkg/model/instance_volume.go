package model

import (
	"time"
)

// Volume represents a cloud disk volume
type Volume struct {
	VolumeID    *string             `json:"volume_id"`
	Attachments []*VolumeAttachment `json:"attachments"`
	Name        *string             `json:"name"`
	Size        *int64              `json:"size"`         // GB
	Type        *string             `json:"type"`         // volume type
	Status      *string             `json:"status"`       // volume status
	Zone        *string             `json:"zone"`         // availability zone
	Profile     string              `json:"profile"`      // cloud profile
	Tags        *Tags               `json:"tags"`         // volume tags
	CreatedTime *time.Time          `json:"created_time"` // creation time
}

type VolumeAttachment struct {
	InstanceID         *string `json:"instance_id"`
	DeleteWithInstance *bool   `json:"delete_with_instance"`
}

// DescribeVolumesInput represents the input for describing volumes
type DescribeVolumesInput struct {
	VolumeIDs  []string `json:"volume_ids"`
	InstanceID *string  `json:"instance_id"`
}

// CreateVolumeInput represents the input for creating a volume
type CreateVolumeInput struct {
	Name               *string `json:"name"`
	Size               *int64  `json:"size"`        // GB
	Type               *string `json:"type"`        // volume type
	Zone               *string `json:"zone"`        // availability zone
	InstanceID         *string `json:"instance_id"` // attach to instance
	Tags               Tags    `json:"tags"`        // volume tags
	DeleteWithInstance *bool   `json:"delete_with_instance"`
}

// DeleteVolumeInput represents the input for deleting a volume
type DeleteVolumeInput struct {
	VolumeID *string `json:"volume_id"`
	Force    *bool   `json:"force"` // force delete if volume is in-use
}

// ModifyVolumeInput represents the input for modifying a volume
type ModifyVolumeInput struct {
	VolumeID *string `json:"volume_id"`
	Size     *int64  `json:"size"` // new size in GB
	Type     *string `json:"type"` // new volume type
}

// AttachVolumeInput represents the input for attaching a volume to an instance
type AttachVolumeInput struct {
	VolumeID           *string `json:"volume_id"`
	InstanceID         *string `json:"instance_id"`
	DeleteWithInstance *bool   `json:"delete_with_instance"`
}

// DetachVolumeInput represents the input for detaching a volume from an instance
type DetachVolumeInput struct {
	VolumeID   *string `json:"volume_id"`
	InstanceID *string `json:"instance_id"`
	Force      *bool   `json:"force"` // force detach
}
