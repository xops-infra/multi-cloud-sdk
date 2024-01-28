package model

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	tencentTag "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tag/v20180813"
	tencentVpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
)

type Tags []Tag

type CreateTagsInput struct {
	Tags Tags
}

// to string
func (t Tags) ToString() string {
	var tags string
	for _, tag := range t {
		tags += tag.Key + ":" + tag.Value + ","
	}
	return strings.TrimSuffix(tags, ",")
}

func (t Tags) GetName() *string {
	for _, tag := range t {
		if tag.Key == "Name" {
			return aws.String(tag.Value)
		}
	}
	return nil
}

// get Owner
func (t Tags) GetOwner() *string {
	for _, tag := range t {
		if tag.Key == "Owner" {
			return aws.String(tag.Value)
		}
	}
	return nil
}

// get EnvType
func (t Tags) GetEnvType() *string {
	for _, tag := range t {
		if tag.Key == "EnvType" {
			return aws.String(tag.Value)
		}
	}
	return nil
}

// get Team
func (t Tags) GetTeam() *string {
	for _, tag := range t {
		if tag.Key == "Team" {
			return aws.String(tag.Value)
		}
	}
	return nil
}

type Tag struct {
	Key   string
	Value string
}

// aws tags to model tags []*ec2.Tag
func AwsTagsToModelTags(tags []*ec2.Tag) *Tags {
	var modelTags Tags
	for _, tag := range tags {
		modelTags = append(modelTags, Tag{
			Key:   *tag.Key,
			Value: *tag.Value,
		})
	}
	return &modelTags
}

func (t *Tags) ToTencentTags() []*tencentTag.Tag {
	var tencentTags []*tencentTag.Tag
	for _, tag := range *t {
		tencentTags = append(tencentTags, &tencentTag.Tag{
			TagKey:   aws.String(tag.Key),
			TagValue: aws.String(tag.Value),
		})
	}
	return tencentTags
}

// tencent tags to model tags
func TencentTagsToModelTags(tags []*cvm.Tag) *Tags {
	var modelTags Tags
	for _, tag := range tags {
		modelTags = append(modelTags, Tag{
			Key:   *tag.Key,
			Value: *tag.Value,
		})
	}
	return &modelTags
}

func TencentVpcTagsFmt(tags []*tencentVpc.Tag) *Tags {
	var modelTags Tags
	for _, tag := range tags {
		modelTags = append(modelTags, Tag{
			Key:   *tag.Key,
			Value: *tag.Value,
		})
	}
	return &modelTags
}

type AddTagsInput struct {
	Tags         Tags      `json:"tags"`
	ResourceList []*string `json:"resource_list"` // 资源列表(qcs::cvm:ap-beijing:uin/1234567:instance/ins-123)
}

type RemoveTagsInput struct {
	ResourceList []*string `json:"resource_list"` // 资源列表(qcs::cvm:ap-beijing:uin/1234567:instance/ins-123)
	Keys         []*string `json:"keys"`
}

type ModifyTagsInput struct {
	Resource *string `json:"resource" binding:"required"`
	Key      *string `json:"key" binding:"required"`
	Value    *string `json:"value" binding:"required"`
}
