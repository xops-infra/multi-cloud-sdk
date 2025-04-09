package io

import (
	"time"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/spf13/cast"
	cbs "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cbs/v20170312"
	"github.com/xops-infra/multi-cloud-sdk/pkg/model"
)

func (i *tencentClient) DescribeVolumes(profile, region string, input model.DescribeVolumesInput) ([]model.Volume, error) {
	svc, err := i.io.GetTencentCbsClient(profile, region)
	if err != nil {
		return nil, err
	}
	request := cbs.NewDescribeDisksRequest()
	request.Limit = tea.Uint64(50)
	if len(input.VolumeIDs) > 0 {
		var diskIds []*string
		for _, volumeID := range input.VolumeIDs {
			diskIds = append(diskIds, tea.String(volumeID))
		}
		request.DiskIds = diskIds
	}

	volumes := make([]model.Volume, 0)
	for {
		response, err := svc.DescribeDisks(request)
		if err != nil {
			return nil, err
		}
		for _, disk := range response.Response.DiskSet {
			size := cast.ToInt64(disk.DiskSize)
			tags := make(model.Tags, 0)
			for _, tag := range disk.Tags {
				if tag.Key != nil && tag.Value != nil {
					tags = append(tags, model.Tag{
						Key:   *tag.Key,
						Value: *tag.Value,
					})
				}
			}
			// "2025-04-09 11:00:20"
			createdTime, _ := time.Parse("2006-01-02 15:04:05", *disk.CreateTime)
			var attachments []*model.VolumeAttachment
			for _, instanceID := range disk.InstanceIdList {
				attachments = append(attachments, &model.VolumeAttachment{
					InstanceID:         instanceID,
					DeleteWithInstance: disk.DeleteWithInstance,
				})
			}
			volumes = append(volumes, model.Volume{
				VolumeID:    disk.DiskId,
				Name:        disk.DiskName,
				Attachments: attachments,
				Size:        &size,
				Type:        disk.DiskType,
				Status:      disk.DiskState,
				Zone:        disk.Placement.Zone,
				Profile:     profile,
				Region:      region,
				Tags:        &tags,
				CreatedTime: &createdTime,
			})
		}
		if len(volumes) >= int(*response.Response.TotalCount) {
			break
		}
		request.Offset = tea.Uint64(uint64(len(volumes)))
	}
	return volumes, nil
}

func (i *tencentClient) CreateVolume(profile, region string, input model.CreateVolumeInput) (string, error) {
	panic("not implemented")
}

func (i *tencentClient) ModifyVolume(profile, region string, input model.ModifyVolumeInput) error {
	panic("not implemented")
}

func (i *tencentClient) DeleteVolume(profile, region string, input model.DeleteVolumeInput) error {
	panic("not implemented")
}

func (i *tencentClient) AttachVolume(profile, region string, input model.AttachVolumeInput) error {
	panic("not implemented")
}

func (i *tencentClient) DetachVolume(profile, region string, input model.DetachVolumeInput) error {
	panic("not implemented")
}
