package model

import "time"

type VpcContract interface {
	QueryVPCs(input CommonQueryInput) ([]*VPC, error)
	GetVPC(vpc_id string) (*VPC, error)
	QuerySubnets(input CommonQueryInput) ([]*Subnet, error)
	QueryEIPs(input CommonQueryInput) ([]*EIP, error)
	QueryNATs(input CommonQueryInput) ([]*NAT, error)
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
