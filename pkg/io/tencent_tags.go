package io

import (
	"fmt"
	"strings"

	"github.com/alibabacloud-go/tea/tea"
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

func (c *tencentClient) AddTagsToResource(profile, region string, input model.AddTagsInput) error {
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

func (c *tencentClient) RemoveTagsFromResource(profile, region string, input model.RemoveTagsInput) error {
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

func (c *tencentClient) ModifyTagsForResource(profile, region string, input model.ModifyTagsInput) error {
	svc, err := c.io.GetTencentTagsClient(profile, region)
	if err != nil {
		return err
	}
	resource, err := getResourceById(region, *input.InstanceId)
	if err != nil {
		return err
	}
	request := tencentTag.NewUpdateResourceTagValueRequest()
	request.Resource = tea.String(resource)
	request.TagKey = input.Key
	request.TagValue = input.Value
	_, err = svc.UpdateResourceTagValue(request)
	return err
}

/*
兼容不同资源的 resource 六点写法 qcs:project_id:service_type:region:account:resource
根据资源ID前缀判断资源类型：
- ins-*: CVM 实例
- cdb-*: 云数据库
- clb-*: 负载均衡
*/
func getResourceById(region, id string) (string, error) {
	// 根据实例ID前缀判断资源类型
	var serviceType, resource string
	switch {
	case strings.HasPrefix(id, "ins-"):
		serviceType = "cvm"
		resource = "instance"
	case strings.HasPrefix(id, "cdb-"):
		serviceType = "cdb"
		resource = "instanceId"
	default:
		return "", fmt.Errorf("unsupported resource type for instance id: %s", id)
	}
	return fmt.Sprintf("qcs::%s:%s::%s/%s", serviceType, region, resource, id), nil
}
