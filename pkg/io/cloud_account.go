package io

import (
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/emr"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/aws/aws-sdk-go/service/route53domains"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/xops-infra/multi-cloud-sdk/pkg/model"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	dnspod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323"
	tencentEmr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/emr/v20190103"
	ocr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ocr/v20181119"
	tag "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tag/v20180813"
	tiia "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tiia/v20190529"
	tencentVpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
)

type cloudClient struct {
	profiles          []model.ProfileConfig
	awsCredential     map[string]*credentials.Credentials
	tencentCredential map[string]*common.Credential
}

func NewCloudClient(profiles []model.ProfileConfig) model.ClientIo {
	awsCredential := make(map[string]*credentials.Credentials)
	tencentCredential := make(map[string]*common.Credential)
	for _, profile := range profiles {
		switch profile.Cloud {
		case model.AWS:
			for _, region := range profile.Regions {
				credentials := credentials.NewStaticCredentials(profile.AK, profile.SK, "")
				awsCredential[profile.Name+region] = credentials
			}
		case model.TENCENT:
			for _, region := range profile.Regions {
				tencentCredential[profile.Name+region] = common.NewTokenCredential(profile.AK, profile.SK, "")
			}
		default:
		}
	}
	return &cloudClient{
		profiles:          profiles,
		awsCredential:     awsCredential,
		tencentCredential: tencentCredential,
	}
}

func (c *cloudClient) GetTencentCvmClient(accountId, region string) (*cvm.Client, error) {
	if _, ok := c.tencentCredential[accountId+region]; !ok {
		return nil, fmt.Errorf("tencent credential %s-%s not found", accountId, region)
	}
	credential := c.tencentCredential[accountId+region]
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.ReqMethod = "GET"
	cpf.SignMethod = "HmacSHA1"
	client, err := cvm.NewClient(
		credential,
		region,
		cpf)
	if err != nil {
		println("fail to init client: %v", err)
		// fmt.Sprintfs("fail to init client: %v", err)
	}
	return client, nil
}

// get awsSession
func (c *cloudClient) getAWSSession(accountId, region string) (*session.Session, error) {
	if c.awsCredential[accountId+region] == nil {
		return nil, fmt.Errorf("aws credential %s-%s not found", accountId, region)
	}
	sess, err := session.NewSession(aws.NewConfig().WithCredentials(c.awsCredential[accountId+region]))
	if err != nil {
		return nil, err
	}
	// set default region form profile
	sess.Config.Region = aws.String(region)
	return sess, nil
}

func (c *cloudClient) GetAWSCredential(accountId, region string) (*credentials.Credentials, error) {
	if _, ok := c.awsCredential[accountId+region]; !ok {
		return nil, fmt.Errorf("aws credential %s-%s not found", accountId, region)
	}
	return c.awsCredential[accountId+region], nil
}

// get awsEmrClient
func (c *cloudClient) GetAWSEmrClient(accountId, region string) (*emr.EMR, error) {
	sess, err := c.getAWSSession(accountId, region)
	if err != nil {
		return nil, err
	}
	return emr.New(sess), nil
}

// getAwsVpcClient
func (c *cloudClient) GetAwsEc2Client(accountId, region string) (*ec2.EC2, error) {
	sess, err := c.getAWSSession(accountId, region)
	if err != nil {
		return nil, err
	}
	return ec2.New(sess), nil
}

// getAwsObjectStorageClient
func (c *cloudClient) GetAWSObjectStorageClient(accountId, region string) (*s3.S3, error) {
	// s3 不需要指定 region，但是需要指定 endpoint
	// endpoint 不能是 eu-central-1 否则无响应，直到超时
	sess, err := c.getAWSSession(accountId, region)
	if err != nil {
		return nil, err
	}
	if *sess.Config.Region == "eu-central-1" {
		// The authorization header is malformed; the region 'eu-central-1' is wrong; expecting 'us-east-1'
		sess.Config.Region = aws.String("us-east-1")
	}
	return s3.New(sess), nil
}

// getAwsRoute53Client
func (c *cloudClient) GetAwsRoute53Client(accountId, region string) (*route53.Route53, error) {
	sess, err := c.getAWSSession(accountId, region)
	if err != nil {
		return nil, err
	}
	return route53.New(sess), nil
}

// getAwsRoute53DomainClient
func (c *cloudClient) GetAwsRoute53DomainClient(accountId, region string) (*route53domains.Route53Domains, error) {
	sess, err := c.getAWSSession(accountId, region)
	if err != nil {
		return nil, err
	}
	client := route53domains.New(sess)

	return client, nil
}

func (c *cloudClient) GetTencentCredential(accountId, region string) (*common.Credential, error) {
	if _, ok := c.tencentCredential[accountId+region]; !ok {
		return nil, fmt.Errorf("tencent credential %s-%s not found", accountId, region)
	}
	return c.tencentCredential[accountId+region], nil
}

// get TencnetEmrClient
func (c *cloudClient) GetTencentEmrClient(accountId, region string) (*tencentEmr.Client, error) {
	credential := c.tencentCredential[accountId+region]
	clientProfile := profile.NewClientProfile()
	return tencentEmr.NewClient(credential, region, clientProfile)
}

// getTencentVpcClient
func (c *cloudClient) GetTencentVpcClient(accountId, region string) (*tencentVpc.Client, error) {
	credential := c.tencentCredential[accountId+region]
	clientProfile := profile.NewClientProfile()
	return tencentVpc.NewClient(credential, region, clientProfile)
}

func (c *cloudClient) GetTencentObjectStorageClient(accountId, region string) (*cos.Client, error) {
	credential := c.tencentCredential[accountId+region]
	client := cos.NewClient(nil, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  credential.SecretId,
			SecretKey: credential.SecretKey,
		},
	})
	return client, nil
}

// GetTencentTagsClient
func (c *cloudClient) GetTencentTagsClient(accountId, region string) (*tag.Client, error) {
	credential := c.tencentCredential[accountId+region]
	clientProfile := profile.NewClientProfile()
	return tag.NewClient(credential, region, clientProfile)
}

// GetTencentOcrClient
func (c *cloudClient) GetTencentOcrClient(accountId, region string) (*ocr.Client, error) {
	credential := c.tencentCredential[accountId+region]
	clientProfile := profile.NewClientProfile()
	return ocr.NewClient(credential, region, clientProfile)
}

// GetTencentOcrTiiaClient
func (c *cloudClient) GetTencentOcrTiiaClient(accountId, region string) (*tiia.Client, error) {
	credential := c.tencentCredential[accountId+region]
	clientProfile := profile.NewClientProfile()
	return tiia.NewClient(credential, region, clientProfile)
}

func (c *cloudClient) GetTencentDnsPodClient(accountId, region string) (*dnspod.Client, error) {
	credential := c.tencentCredential[accountId+region]
	clientProfile := profile.NewClientProfile()
	return dnspod.NewClient(credential, region, clientProfile)
}
