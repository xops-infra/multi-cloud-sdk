package io

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/emr"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/aws/aws-sdk-go/service/route53domains"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/xops-infra/multi-cloud-sdk/pkg/model"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	dnspod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323"
	tencentEmr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/emr/v20190103"
	ocr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ocr/v20181119"
	privatedns "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/privatedns/v20201028"
	tag "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tag/v20180813"
	tiia "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tiia/v20190529"
	tencentVpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
)

type cloudClient struct {
	profiles          map[string]model.ProfileConfig
	awsCredential     map[string]*credentials.Credentials
	tencentCredential map[string]*common.Credential
}

func NewCloudClient(profiles []model.ProfileConfig) model.ClientIo {
	awsCredential := make(map[string]*credentials.Credentials)
	tencentCredential := make(map[string]*common.Credential)
	_profiles := make(map[string]model.ProfileConfig)
	for _, profile := range profiles {
		_profiles[profile.Name] = profile
		switch profile.Cloud {
		case model.AWS:
			awsCredential[profile.Name] = credentials.NewStaticCredentials(profile.AK, profile.SK, "")
		case model.TENCENT:
			tencentCredential[profile.Name] = common.NewTokenCredential(profile.AK, profile.SK, "")
		default:
		}
	}
	return &cloudClient{
		profiles:          _profiles,
		awsCredential:     awsCredential,
		tencentCredential: tencentCredential,
	}
}

func (c *cloudClient) GetTencentCvmClient(accountId, region string) (*cvm.Client, error) {
	if _, ok := c.tencentCredential[accountId]; !ok {
		return nil, fmt.Errorf("tencent credential %s not found", accountId)
	}
	credential := c.tencentCredential[accountId]
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
func (c *cloudClient) getAWSSession(accountId string) (*session.Session, error) {
	if c.awsCredential[accountId] == nil {
		return nil, fmt.Errorf("aws credential %s not found", accountId)
	}
	sess, err := session.NewSession(aws.NewConfig().WithCredentials(c.awsCredential[accountId]))
	if err != nil {
		return nil, err
	}
	return sess, nil
}

func (c *cloudClient) GetAWSCredential(accountId string) (*credentials.Credentials, error) {
	if _, ok := c.awsCredential[accountId]; !ok {
		return nil, fmt.Errorf("aws credential %s not found", accountId)
	}
	return c.awsCredential[accountId], nil
}

// get awsEmrClient
func (c *cloudClient) GetAWSEmrClient(accountId, region string) (*emr.EMR, error) {
	sess, err := c.getAWSSession(accountId)
	if err != nil {
		return nil, err
	}
	sess.Config.Region = aws.String(region) // set default region
	return emr.New(sess), nil
}

// getAwsVpcClient
func (c *cloudClient) GetAwsEc2Client(accountId, region string) (*ec2.EC2, error) {
	sess, err := c.getAWSSession(accountId)
	if err != nil {
		return nil, err
	}
	sess.Config.Region = aws.String(region)
	return ec2.New(sess), nil
}

func (c *cloudClient) GetAwsSqsClient(accountId, region string) (*sqs.SQS, error) {
	sess, err := c.getAWSSession(accountId)
	if err != nil {
		return nil, err
	}
	sess.Config.Region = aws.String(region)
	return sqs.New(sess), nil
}

// getAwsObjectStorageClient
func (c *cloudClient) GetAWSS3Client(accountId, region string) (*s3.S3, error) {
	// s3 不需要指定 region，但是需要指定 endpoint
	// endpoint 不能是 eu-central-1 否则无响应，直到超时
	sess, err := c.getAWSSession(accountId)
	if err != nil {
		return nil, err
	}
	sess.Config.Region = tea.String(region)
	return s3.New(sess), nil
}

func (c *cloudClient) GetAwsRoute53Client(accountId, region string) (*route53.Route53, error) {
	if region == "" {
		return nil, fmt.Errorf("region is empty")
	}
	sess, err := c.getAWSSession(accountId)
	if err != nil {
		return nil, err
	}
	sess.Config.Region = tea.String(region) // must has region
	return route53.New(sess), nil
}

// getAwsRoute53DomainClient
func (c *cloudClient) GetAwsRoute53DomainClient(accountId string) (*route53domains.Route53Domains, error) {
	sess, err := c.getAWSSession(accountId)
	if err != nil {
		return nil, err
	}
	client := route53domains.New(sess)

	return client, nil
}

func (c *cloudClient) getTencentCredential(accountId string) (*common.Credential, error) {
	credential, ok := c.tencentCredential[accountId]
	if !ok {
		var keys string
		for k, _ := range c.tencentCredential {
			keys += k + ","
		}
		return nil, fmt.Errorf("tencent credential %s not found, Keys: %s", accountId, strings.TrimSuffix(keys, ","))
	}
	return credential, nil
}

// get TencnetEmrClient
func (c *cloudClient) GetTencentEmrClient(accountId, region string) (*tencentEmr.Client, error) {
	credential, err := c.getTencentCredential(accountId)
	if err != nil {
		return nil, err
	}
	clientProfile := profile.NewClientProfile()
	return tencentEmr.NewClient(credential, region, clientProfile)
}

// getTencentVpcClient
func (c *cloudClient) GetTencentVpcClient(accountId, region string) (*tencentVpc.Client, error) {
	credential, err := c.getTencentCredential(accountId)
	if err != nil {
		return nil, err
	}
	clientProfile := profile.NewClientProfile()
	return tencentVpc.NewClient(credential, region, clientProfile)
}

func (c *cloudClient) GetTencentCosClient(accountId, region string) (*cos.Client, error) {
	host := "https://service.cos.myqcloud.com"
	if region != "" {
		host = fmt.Sprintf("https://cos.%s.myqcloud.com", region)
	}

	credential, err := c.getTencentCredential(accountId)
	if err != nil {
		return nil, err
	}
	u, _ := url.Parse(host)
	b := &cos.BaseURL{
		ServiceURL: u,
	}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  credential.SecretId,
			SecretKey: credential.SecretKey,
		},
	})
	return client, nil
}

func (c *cloudClient) GetTencentCosLifecycleClient(accountId, region, bucket string) (*cos.Client, error) {
	u, _ := url.Parse(fmt.Sprintf("https://%s.cos.%s.myqcloud.com", bucket, region))
	b := &cos.BaseURL{BucketURL: u}

	credential, err := c.getTencentCredential(accountId)
	if err != nil {
		return nil, err
	}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  credential.SecretId,
			SecretKey: credential.SecretKey,
		},
	})
	return client, nil
}

// GetTencentTagsClient
func (c *cloudClient) GetTencentTagsClient(accountId, region string) (*tag.Client, error) {
	credential, err := c.getTencentCredential(accountId)
	if err != nil {
		return nil, err
	}
	clientProfile := profile.NewClientProfile()
	return tag.NewClient(credential, region, clientProfile)
}

// GetTencentOcrClient
func (c *cloudClient) GetTencentOcrClient(accountId, region string) (*ocr.Client, error) {
	credential, err := c.getTencentCredential(accountId)
	if err != nil {
		return nil, err
	}
	clientProfile := profile.NewClientProfile()
	return ocr.NewClient(credential, region, clientProfile)
}

// GetTencentOcrTiiaClient
func (c *cloudClient) GetTencentOcrTiiaClient(accountId, region string) (*tiia.Client, error) {
	credential, err := c.getTencentCredential(accountId)
	if err != nil {
		return nil, err
	}
	clientProfile := profile.NewClientProfile()
	return tiia.NewClient(credential, region, clientProfile)
}

func (c *cloudClient) GetTencentDnsPodClient(accountId string) (*dnspod.Client, error) {
	credential, err := c.getTencentCredential(accountId)
	if err != nil {
		return nil, err
	}
	clientProfile := profile.NewClientProfile()
	return dnspod.NewClient(credential, "", clientProfile)
}

func (c *cloudClient) GetTencentPrivateDNSClient(accountId string) (*privatedns.Client, error) {
	credential, ok := c.tencentCredential[accountId]
	if !ok {
		return nil, fmt.Errorf("tencent credential %s not found", accountId)
	}
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "privatedns.tencentcloudapi.com"
	client, err := privatedns.NewClient(credential, "", cpf)
	if err != nil {
		return nil, err
	}
	return client, nil
}
