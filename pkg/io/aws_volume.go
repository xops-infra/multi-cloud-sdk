package io

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/xops-infra/multi-cloud-sdk/pkg/model"
)

func (i *awsClient) DescribeVolumes(profile, region string, input model.DescribeVolumesInput) ([]model.Volume, error) {

	volumes := make([]model.Volume, 0)

	svc, err := i.io.GetAwsEc2Client(profile, region)
	if err != nil {
		return nil, err
	}
	request := &ec2.DescribeVolumesInput{}
	if input.VolumeIDs != nil {
		for _, volumeID := range input.VolumeIDs {
			request.VolumeIds = append(request.VolumeIds, aws.String(volumeID))
		}
	}
	for {
		out, err := svc.DescribeVolumes(request)
		if err != nil {
			return nil, err
		}
		for _, volume := range out.Volumes {
			tags := make(model.Tags, 0)
			for _, tag := range volume.Tags {
				tags = append(tags, model.Tag{
					Key:   *tag.Key,
					Value: *tag.Value,
				})
			}
			var attachments []*model.VolumeAttachment
			for _, attachment := range volume.Attachments {
				attachments = append(attachments, &model.VolumeAttachment{
					InstanceID:         attachment.InstanceId,
					DeleteWithInstance: attachment.DeleteOnTermination,
				})
			}
			volumes = append(volumes, model.Volume{
				VolumeID:    volume.VolumeId,
				Attachments: attachments,
				Name:        tags.GetTagValueByKey("Name"),
				Type:        volume.VolumeType,
				Size:        volume.Size,
				Zone:        volume.AvailabilityZone,
				Profile:     profile,
				Status:      volume.State,
				Tags:        &tags,
				CreatedTime: volume.CreateTime,
			})
		}
		if out.NextToken == nil {
			break
		}
		request.NextToken = out.NextToken
	}

	return volumes, nil
}

func (i *awsClient) CreateVolume(profile, region string, input model.CreateVolumeInput) (string, error) {
	panic("not implemented")
}

func (i *awsClient) ModifyVolume(profile, region string, input model.ModifyVolumeInput) error {
	panic("not implemented")
}

func (i *awsClient) DeleteVolume(profile, region string, input model.DeleteVolumeInput) error {
	panic("not implemented")
}

func (i *awsClient) AttachVolume(profile, region string, input model.AttachVolumeInput) error {
	panic("not implemented")
}

func (i *awsClient) DetachVolume(profile, region string, input model.DetachVolumeInput) error {
	panic("not implemented")
}
