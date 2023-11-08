package model

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	tencentVpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
)

type Instance struct {
	Name       *string   `json:"name"`
	InstanceID *string   `json:"instance_id" gorm:"primarykey"`
	Profile    string    `json:"profile"`
	KeyName    []*string `json:"key_name" gorm:"serializer:json"`
	Region     *string   `json:"region"`
	PrivateIP  []*string `json:"private_ip" gorm:"serializer:json"`
	Platform   *string   `json:"platform"`
	PublicIP   []*string `json:"public_ip" gorm:"serializer:json"`
	Status     *string   `json:"status"`
	Owner      *string   `json:"owner"`
	Tags       *Tags     `json:"tags" gorm:"serializer:json"`
}

type InstanceQueryInput struct {
	Name          string         `json:"name"`           // 机器名称，使用字符串包含匹配方式
	ID            string         `json:"id"`             // 机器ID，使用字符串包含匹配方式
	Region        string         `json:"region"`         // 区域
	Account       string         `json:"account"`        // 账号
	CloudProvider CloudProvider  `json:"cloud_provider"` // 云平台
	Ip            string         `json:"ip"`             // 机器IP，支持公网IP和内网IP，也是包含匹配方式
	Status        InstanceStatus `json:"status"`         // 机器状态
	Owner         string         `json:"owner"`          // 机器所有者，tags的Owner
}

// InstanceQueryInput filter
func (i InstanceQueryInput) Filter(instances []*Instance) []*Instance {
	var newInstances []*Instance
	instancesMap := make(map[string]*Instance)
	// 加速，如果没有任何条件，直接返回
	if i.Name == "" && i.Ip == "" && i.Status.ToString() == "" && i.Owner == "" && i.Region == "" && i.Account == "" && i.CloudProvider.ToString() == "" {
		return instances
	}
	for _, instance := range instances {
		if _, ok := instancesMap[*instance.InstanceID]; ok {
			continue
		}
		if i.Name != "" && instance.Name != nil && strings.Contains(*instance.Name, i.Name) {
			newInstances = append(newInstances, instance)
			instancesMap[*instance.InstanceID] = instance
			continue
		}
		if i.Ip != "" && instance.PrivateIP != nil {
			if checkIp(instance.PrivateIP, i) {
				newInstances = append(newInstances, instance)
				instancesMap[*instance.InstanceID] = instance
				continue
			}
		}
		if i.Ip != "" && instance.PublicIP != nil {
			if checkIp(instance.PublicIP, i) {
				newInstances = append(newInstances, instance)
				instancesMap[*instance.InstanceID] = instance
				continue
			}
		}
		if i.Status.ToString() != "" && instance.Status != nil && *instance.Status == i.Status.ToString() {
			newInstances = append(newInstances, instance)
			instancesMap[*instance.InstanceID] = instance
			continue
		}
		if i.Owner != "" && instance.Owner != nil && *instance.Owner == i.Owner {
			newInstances = append(newInstances, instance)
			instancesMap[*instance.InstanceID] = instance
			continue
		}
		if i.Region != "" && instance.Region != nil && *instance.Region == i.Region {
			newInstances = append(newInstances, instance)
			instancesMap[*instance.InstanceID] = instance
			continue
		}
		if i.Account != "" && instance.Profile == i.Account {
			newInstances = append(newInstances, instance)
			instancesMap[*instance.InstanceID] = instance
			continue
		}
		if i.CloudProvider != "" && instance.Profile == i.CloudProvider.ToString() {
			newInstances = append(newInstances, instance)
			instancesMap[*instance.InstanceID] = instance
			continue
		}
	}
	return newInstances
}

func checkIp(ips []*string, i InstanceQueryInput) bool {
	for _, ip := range ips {
		if ip != nil && strings.Contains(*ip, i.Ip) {
			return true
		}
	}
	return false
}

// InstanceStatus to string
func (i InstanceStatus) ToString() string {
	return string(i)
}

type InstanceResponse Instance

type Tags []Tag

// to string
func (t Tags) ToString() string {
	var tags string
	for _, tag := range t {
		tags += tag.Key + ":" + tag.Value + ","
	}
	strings.TrimSuffix(tags, ",")
	return tags
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
