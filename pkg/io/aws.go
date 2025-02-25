package io

import (
	"fmt"

	"github.com/xops-infra/multi-cloud-sdk/pkg/model"
)

type awsClient struct {
	io model.ClientIo
}

func NewAwsClient(io model.ClientIo) model.CloudIO {
	return &awsClient{
		io: io,
	}
}

func (c *awsClient) CreateTags(profile, region string, input model.CreateTagsInput) error {
	return fmt.Errorf("not support for aws")
}

func (c *awsClient) AddTagsToResource(profile, region string, input model.AddTagsInput) error {
	return fmt.Errorf("not support for aws")
}

func (c *awsClient) RemoveTagsFromResource(profile, region string, input model.RemoveTagsInput) error {
	return fmt.Errorf("not support for aws")
}

func (c *awsClient) ModifyTagsForResource(profile, region string, input model.ModifyTagsInput) error {
	return fmt.Errorf("not support for aws")
}

// CommonOCR
func (c *awsClient) CommonOCR(profile, region string, input model.OcrRequest) (model.OcrResponse, error) {
	return model.OcrResponse{}, nil
}

// CreatePicture
func (c *awsClient) CreatePicture(profile, region string, input model.CreatePictureRequest) (model.CreatePictureResponse, error) {
	return model.CreatePictureResponse{}, nil
}

// GetPictureByName
func (c *awsClient) GetPictureByName(profile, region string, input model.CommonPictureRequest) (model.GetPictureByNameResponse, error) {
	return model.GetPictureByNameResponse{}, nil
}

// QueryPicture
func (c *awsClient) QueryPicture(profile, region string, input model.QueryPictureRequest) (model.QueryPictureResponse, error) {
	return model.QueryPictureResponse{}, nil
}

// DeletePicture
func (c *awsClient) DeletePicture(profile, region string, input model.CommonPictureRequest) (model.CommonPictureResponse, error) {
	return model.CommonPictureResponse{}, nil
}

// UpdatePicture
func (c *awsClient) UpdatePicture(profile, region string, input model.UpdatePictureRequest) (model.CommonPictureResponse, error) {
	return model.CommonPictureResponse{}, nil
}

// SearchPicture
func (c *awsClient) SearchPicture(profile, region string, input model.SearchPictureRequest) (model.SearchPictureResponse, error) {
	return model.SearchPictureResponse{}, nil
}

func (c *awsClient) CreateSqs(profile, region string, input model.CreateSqsRequest) error {
	client, err := c.io.GetAwsSqsClient(profile, region)
	if err != nil {
		return fmt.Errorf("get aws sqs client failed, %v", err)
	}

	_, err = client.CreateQueue(input.ToCreateQueueInput())
	if err != nil {
		return fmt.Errorf("create aws sqs failed, %v", err)
	}

	return nil
}
