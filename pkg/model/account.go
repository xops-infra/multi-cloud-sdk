package model

import (
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/emr"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/aws/aws-sdk-go/service/route53domains"
	"github.com/aws/aws-sdk-go/service/s3"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	dnspod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323"
	tencentEmr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/emr/v20190103"
	ocr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ocr/v20181119"
	privatedns "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/privatedns/v20201028"
	tag "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tag/v20180813"
	tiia "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tiia/v20190529"
	tencentVpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	cos "github.com/tencentyun/cos-go-sdk-v5"
)

type ClientIo interface {
	GetAwsEc2Client(profile, region string) (*ec2.EC2, error)
	GetAWSEmrClient(profile, region string) (*emr.EMR, error)
	GetAWSCredential(profile string) (*credentials.Credentials, error)
	GetAWSObjectStorageClient(profile, region string) (*s3.S3, error)
	GetAwsRoute53Client(profile string) (*route53.Route53, error)
	GetAwsRoute53DomainClient(profile string) (*route53domains.Route53Domains, error)

	GetTencentCvmClient(profile, region string) (*cvm.Client, error)
	GetTencentEmrClient(profile, region string) (*tencentEmr.Client, error)
	GetTencentVpcClient(profile, region string) (*tencentVpc.Client, error)
	GetTencentObjectStorageClient(profile, region string) (*cos.Client, error)
	GetTencentTagsClient(profile, region string) (*tag.Client, error)
	GetTencentOcrClient(profile, region string) (*ocr.Client, error)
	GetTencentOcrTiiaClient(profile, region string) (*tiia.Client, error)
	GetTencentDnsPodClient(profile string) (*dnspod.Client, error)
	GetTencentPrivateDNSClient(profile string) (*privatedns.Client, error)
}

type ProfileConfig struct {
	Name string `mapstructure:"name" binding:"required"`
	AK   string `mapstructure:"ak" binding:"required"`
	SK   string `mapstructure:"sk" binding:"required"`
	// Regions []string `mapstructure:"regions" binding:"required"` // init clinet must have region
	Cloud Cloud `mapstructure:"cloud" binding:"required"`
}
