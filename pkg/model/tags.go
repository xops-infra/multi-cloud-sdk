package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/emr"
	"github.com/aws/aws-sdk-go/service/s3"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	tencentEmr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/emr/v20190103"
	tencentTag "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tag/v20180813"
	tencentVpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	cos "github.com/tencentyun/cos-go-sdk-v5"
)

type Tags []Tag

// to aws s3 tags
func (t Tags) ToAWSS3Tags() []*s3.Tag {
	var tags []*s3.Tag
	for _, tag := range t {
		tags = append(tags, &s3.Tag{
			Key:   aws.String(tag.Key),
			Value: aws.String(tag.Value),
		})
	}
	return tags
}

// to aws emr tags
func (t Tags) ToAWSEmrTags() []*emr.Tag {
	var tags []*emr.Tag
	for _, tag := range t {
		tags = append(tags, &emr.Tag{
			Key:   aws.String(tag.Key),
			Value: aws.String(tag.Value),
		})
	}
	return tags
}

// to tencent cos tags
func (t Tags) ToTencentCosTags() []cos.BucketTaggingTag {
	var tags []cos.BucketTaggingTag
	for _, tag := range t {
		tags = append(tags, cos.BucketTaggingTag{
			Key:   tag.Key,
			Value: tag.Value,
		})
	}
	return tags
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

// 更通用的方法
func (t Tags) GetTagValueByKey(key string) *string {
	for _, tag := range t {
		if tag.Key == key {
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

// add
func (t *Tags) Add(key, value string) {
	*t = append(*t, Tag{
		Key:   key,
		Value: value,
	})
}

// remove
func (t *Tags) Remove(key string) {
	for i, tag := range *t {
		if tag.Key == key {
			*t = append((*t)[:i], (*t)[i+1:]...)
		}
	}
}

// for gorm
func (t *Tags) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, t)
}

// Value implements the Valuer interface for Tags
func (t Tags) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func NewTagsFromAWSS3Tags(tags []*s3.Tag) Tags {
	var modelTags Tags
	for _, tag := range tags {
		modelTags = append(modelTags, Tag{
			Key:   *tag.Key,
			Value: *tag.Value,
		})
	}
	return modelTags
}

func NewTagsFromTencentCosTags(tags []cos.BucketTaggingTag) Tags {
	var modelTags Tags
	for _, tag := range tags {
		modelTags = append(modelTags, Tag{
			Key:   tag.Key,
			Value: tag.Value,
		})
	}
	return modelTags
}

func NewTagsFromTencentEmrTags(tags []*tencentEmr.Tag) Tags {
	var modelTags Tags
	for _, tag := range tags {
		modelTags = append(modelTags, Tag{
			Key:   *tag.TagKey,
			Value: *tag.TagValue,
		})
	}
	return modelTags
}

func NewTagsFromAWSEmrTags(tags []*emr.Tag) Tags {
	var modelTags Tags
	for _, tag := range tags {
		modelTags = append(modelTags, Tag{
			Key:   *tag.Key,
			Value: *tag.Value,
		})
	}
	return modelTags
}

type CreateTagsInput struct {
	Tags Tags
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

// []*emr.Tag
func (t *Tags) ToTencentEmrTags() []*tencentEmr.Tag {
	var tencentTags []*tencentEmr.Tag
	for _, tag := range *t {
		tencentTags = append(tencentTags, &tencentEmr.Tag{
			TagKey:   aws.String(tag.Key),
			TagValue: aws.String(tag.Value),
		})
	}
	return tencentTags
}

func (t *Tags) ToRunInstanceTags() []*cvm.TagSpecification {
	var tencentTags []*cvm.TagSpecification
	for _, tag := range *t {
		tencentTags = append(tencentTags, &cvm.TagSpecification{
			ResourceType: aws.String("instance"),
			Tags: []*cvm.Tag{
				{
					Key:   aws.String(tag.Key),
					Value: aws.String(tag.Value),
				},
			},
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
	InstanceId *string `json:"instance_id" binding:"required"`
	Key        *string `json:"key" binding:"required"`
	Value      *string `json:"value" binding:"required"`
}
