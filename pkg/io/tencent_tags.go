package io

import (
	"strings"

	tencentTag "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tag/v20180813"
	"github.com/xops-infra/multi-cloud-sdk/pkg/model"
)

func (c *tencentClient) CreateTags(profile, region string, input model.CreateTagsInput) error {
	svc, err := c.io.GetTencentTagsClient(profile, region)
	if err != nil {
		return err
	}
	request := tencentTag.NewCreateTagRequest()
	for _, tag := range input.Tags {
		request.TagKey = &tag.Key
		request.TagValue = &tag.Value
		_, err = svc.CreateTag(request)
		if err != nil {
			if strings.Contains(err.Error(), "Message=tagKey-tagValue have exists.") {
				continue
			}
			return err
		}
	}
	return nil
}

func (c *tencentClient) AddTagsFromResource(profile, region string, input model.AddTagsInput) error {
	svc, err := c.io.GetTencentTagsClient(profile, region)
	if err != nil {
		return err
	}
	request := tencentTag.NewTagResourcesRequest()
	request.ResourceList = input.ResourceList
	request.Tags = input.Tags.ToTencentTags()
	_, err = svc.TagResources(request)
	return err
}

func (c *tencentClient) DeleteTagsFromResource(profile, region string, input model.RemoveTagsInput) error {
	svc, err := c.io.GetTencentTagsClient(profile, region)
	if err != nil {
		return err
	}
	request := tencentTag.NewUnTagResourcesRequest()
	request.ResourceList = input.ResourceList
	request.TagKeys = input.Keys
	_, err = svc.UnTagResources(request)
	return err
}

func (c *tencentClient) ModifyTagsFromResource(profile, region string, input model.ModifyTagsInput) error {
	svc, err := c.io.GetTencentTagsClient(profile, region)
	if err != nil {
		return err
	}
	request := tencentTag.NewUpdateResourceTagValueRequest()
	request.Resource = input.Resource
	request.TagKey = input.Key
	request.TagValue = input.Value
	_, err = svc.UpdateResourceTagValue(request)
	return err
}
