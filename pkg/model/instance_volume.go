package model

import (
	"time"

	"github.com/alibabacloud-go/tea/tea"
)

// Volume represents a cloud disk volume
type Volume struct {
	VolumeID           *string    `json:"volume_id"`
	InstanceID         *string    `json:"instance_id"`
	Name               *string    `json:"name"`
	Size               *int64     `json:"size"`                 // GB
	Type               *string    `json:"type"`                 // volume type
	Status             *string    `json:"status"`               // volume status
	Zone               *string    `json:"zone"`                 // availability zone
	Profile            string     `json:"profile"`              // cloud profile
	Tags               *Tags      `json:"tags"`                 // volume tags
	DeleteWithInstance *bool      `json:"delete_with_instance"` // whether delete with instance
	CreatedTime        *time.Time `json:"created_time"`         // creation time
}

// VolumeFilter represents filters for listing volumes
type VolumeFilter struct {
	VolumeIDs  []*string `json:"volume_ids"`
	InstanceID *string   `json:"instance_id"`
	Name       *string   `json:"name"`
	Size       *int64    `json:"size"`
	Type       *string   `json:"type"`
	Zone       *string   `json:"zone"`
	NextMarker *string   `json:"next_marker"`
}

// ToDescribeVolumesInput converts VolumeFilter to DescribeVolumesInput
func (f *VolumeFilter) ToDescribeVolumesInput() DescribeVolumesInput {
	var filters []*Filter

	if f.Name != nil {
		filters = append(filters, &Filter{
			Name:   tea.String("tag:Name"),
			Values: []*string{f.Name},
		})
	}

	if f.Type != nil {
		filters = append(filters, &Filter{
			Name:   tea.String("volume-type"),
			Values: []*string{f.Type},
		})
	}

	if f.Zone != nil {
		filters = append(filters, &Filter{
			Name:   tea.String("availability-zone"),
			Values: []*string{f.Zone},
		})
	}

	return DescribeVolumesInput{
		VolumeIDs:  f.VolumeIDs,
		InstanceID: f.InstanceID,
		Filters:    filters,
		Size:       f.Size,
		NextMarker: f.NextMarker,
	}
}

// DescribeVolumesInput represents the input for describing volumes
type DescribeVolumesInput struct {
	VolumeIDs  []*string `json:"volume_ids"`
	InstanceID *string   `json:"instance_id"`
	Filters    []*Filter `json:"filters"`
	Size       *int64    `json:"size"`
	NextMarker *string   `json:"next_marker"`
}

// DescribeVolumesResponse represents the response for describing volumes
type DescribeVolumesResponse struct {
	Volumes    []Volume `json:"volumes"`
	NextMarker *string  `json:"next_marker"`
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

// CreateVolumeResponse represents the response for creating a volume
type CreateVolumeResponse struct {
	VolumeID  *string `json:"volume_id"`
	RequestID *string `json:"request_id"`
}

// DeleteVolumeInput represents the input for deleting a volume
type DeleteVolumeInput struct {
	VolumeID *string `json:"volume_id"`
	Force    *bool   `json:"force"` // force delete if volume is in-use
}

// DeleteVolumeResponse represents the response for deleting a volume
type DeleteVolumeResponse struct {
	RequestID *string `json:"request_id"`
}

// ModifyVolumeInput represents the input for modifying a volume
type ModifyVolumeInput struct {
	VolumeID *string `json:"volume_id"`
	Size     *int64  `json:"size"` // new size in GB
	Type     *string `json:"type"` // new volume type
}

// ModifyVolumeResponse represents the response for modifying a volume
type ModifyVolumeResponse struct {
	RequestID *string `json:"request_id"`
}

// AttachVolumeInput represents the input for attaching a volume to an instance
type AttachVolumeInput struct {
	VolumeID           *string `json:"volume_id"`
	InstanceID         *string `json:"instance_id"`
	DeleteWithInstance *bool   `json:"delete_with_instance"`
}

// AttachVolumeResponse represents the response for attaching a volume
type AttachVolumeResponse struct {
	RequestID *string `json:"request_id"`
}

// DetachVolumeInput represents the input for detaching a volume from an instance
type DetachVolumeInput struct {
	VolumeID   *string `json:"volume_id"`
	InstanceID *string `json:"instance_id"`
	Force      *bool   `json:"force"` // force detach
}

// DetachVolumeResponse represents the response for detaching a volume
type DetachVolumeResponse struct {
	RequestID *string `json:"request_id"`
}
