package io

import "github.com/xops-infra/multi-cloud-sdk/pkg/model"

func (i *awsClient) DescribeVolumes(profile, region string, input model.DescribeVolumesInput) ([]model.Volume, error) {
	panic("not implemented")
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
