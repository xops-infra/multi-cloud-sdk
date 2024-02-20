package model

import (
	"time"

	tencentVpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
)

type VpcContract interface {
	QueryVPCs(profile, region string, input CommonFilter) ([]VPC, error)
	QuerySubnets(profile, region string, input CommonFilter) ([]Subnet, error)
	QueryEIPs(profile, region string, input CommonFilter) ([]EIP, error)
	QueryNATs(profile, region string, input CommonFilter) ([]NAT, error)
}

type VPC struct {
	ID            string `json:"id"`
	Region        string `json:"region"`
	CloudProvider Cloud  `json:"cloud_provider"`
	Account       string `json:"account"`
	Tags          *Tags  `json:"tags"`
	IsDefault     bool   `json:"is_default"`
	CidrBlock     string `json:"cidr_block"`
}

type Subnet struct {
	ID                      *string    `json:"id"`
	Region                  string     `json:"region"`
	Account                 string     `json:"account"`
	CloudProvider           Cloud      `json:"cloud_provider"`
	Tags                    *Tags      `json:"tags"`
	VpcID                   *string    `json:"vpc_id"`
	Name                    *string    `json:"name"`
	CidrBlock               *string    `json:"cidr_block"`
	AvailableIpAddressCount int64      `json:"available_ip_address_count"`
	IsDefault               *bool      `json:"is_default"`
	Zone                    *string    `json:"zone"`
	RouteTableId            *string    `json:"route_table_id"`
	CreatedTime             *time.Time `json:"created_time"`
	NetworkAclId            *string    `json:"network_acl_id"`
}

type EIP struct {
	ID            *string `json:"id"`
	Region        string  `json:"region"`
	CloudProvider Cloud   `json:"cloud_provider"`
	Account       string  `json:"account"`
	Tags          *Tags   `json:"tags"`
	Name          *string `json:"name"`
	// `EIP`状态，包含'CREATING'(创建中),'BINDING'(绑定中),'BIND'(已绑定),'UNBINDING'(解绑中),'UNBIND'(已解绑),'OFFLINING'(释放中),'BIND_ENI'(绑定悬空弹性网卡)
	Status             *string    `json:"status"`
	AddressIp          *string    `json:"address_ip"`
	InstanceId         *string    `json:"instance_id"`
	CreatedTime        *time.Time `json:"created_time"`
	NetworkInterfaceId string     `json:"network_interface_id"`
	PrivateAddressIp   string     `json:"private_address_ip"`
	Bandwidth          *int64     `json:"bandwidth"`
	InternetChargeType *string    `json:"internet_charge_type"`
}

type NAT struct {
	ID            string    `json:"id"`
	Region        string    `json:"region"`
	Account       string    `json:"account"`
	CloudProvider Cloud     `json:"cloud_provider"`
	Tags          *Tags     `json:"tags"`
	Name          string    `json:"name"`
	CreatedTime   time.Time `json:"created_time"`
	Status        string    `json:"status"`
	AddressIps    []string  `json:"address_ips"`
	VpcID         string    `json:"vpc_id"`
	Zone          *string   `json:"zone"`
	SubnetID      string    `json:"subnet_id"`
}

type CreateSecurityGroupWithPoliciesInput struct {
	GroupName        *string   `json:"group_name" binding:"required"`
	GroupDescription *string   `json:"group_description"`
	PolicySet        PolicySet `json:"policy_set"`
}

type PolicySet struct {
	Egress  []SecurityGroupPolicy `json:"egress"`  // 出站规则
	Ingress []SecurityGroupPolicy `json:"ingress"` // 入站规则
}

// to *tencentVpc.SecurityGroupPolicySet
func (p *PolicySet) ToTencentPolicySet() *tencentVpc.SecurityGroupPolicySet {
	var egress []*tencentVpc.SecurityGroupPolicy
	for _, policy := range p.Egress {
		egress = append(egress, policy.ToTencentPolicy())
	}
	var ingress []*tencentVpc.SecurityGroupPolicy
	for _, policy := range p.Ingress {
		ingress = append(ingress, policy.ToTencentPolicy())
	}
	return &tencentVpc.SecurityGroupPolicySet{
		Egress:  egress,
		Ingress: ingress,
	}
}

type SecurityGroupPolicy struct {
	SecurityGroupId   *string `json:"security_group_id"`
	Protocol          *string `json:"protocol"`           // 协议,取值: TCP,UDP,ICMP,ICMPv6,ALL。
	Port              *string `json:"port"`               // 端口范围，取值:1~65535。示例值：22
	CidrBlock         *string `json:"cidr_block"`         // 来源IP或CIDR 示例值：0.0.0.0/16
	Action            *string `json:"action"`             // ACCEPT 或者 DROP
	PolicyDescription *string `json:"policy_description"` // 描述
	ModifyTime        *string `json:"modify_time"`        // 修改时间
}

// to *tencentVpc.SecurityGroupPolicy
func (policy *SecurityGroupPolicy) ToTencentPolicy() *tencentVpc.SecurityGroupPolicy {
	return &tencentVpc.SecurityGroupPolicy{
		SecurityGroupId:   policy.SecurityGroupId,
		Protocol:          policy.Protocol,
		Port:              policy.Port,
		CidrBlock:         policy.CidrBlock,
		Action:            policy.Action,
		PolicyDescription: policy.PolicyDescription,
		ModifyTime:        policy.ModifyTime,
	}
}

type CreateSecurityGroupWithPoliciesResponse struct {
	Data any
}
