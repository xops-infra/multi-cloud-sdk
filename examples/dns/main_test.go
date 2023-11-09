package main_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/joho/godotenv"
	"github.com/xops-infra/multi-cloud-sdk/pkg/io"
	"github.com/xops-infra/multi-cloud-sdk/pkg/model"
	server "github.com/xops-infra/multi-cloud-sdk/pkg/service"
)

var dnsS model.DnsContract

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	profiles := []model.ProfileConfig{
		{
			Name:  "aws",
			Cloud: model.AWS,
			AK:    os.Getenv("AWS_ACCESS_KEY_ID"),
			SK:    os.Getenv("AWS_SECRET_ACCESS_KEY"),
			Regions: []string{
				"cn-northwest-1",
			},
		},
		{
			Name:  "tencent",
			Cloud: model.TENCENT,
			AK:    os.Getenv("TENCENT_ACCESS_KEY"),
			SK:    os.Getenv("TENCENT_SECRET_KEY"),
			Regions: []string{
				"ap-shanghai",
				// "na-ashburn",
			},
		},
	}
	if profiles[0].AK == "" {
		panic("AWS_ACCESS_KEY_ID is empty")
	}
	cloudIo := io.NewCloudClient(profiles)
	serverTencent := io.NewTencentClient(cloudIo)
	serverAws := io.NewAwsClient(cloudIo)
	dnsS = server.NewDnsService(profiles, serverAws, serverTencent)
}

func TestAwsDomain(t *testing.T) {
	resp, err := dnsS.DescribeDomainList("aws", "cn-northwest-1",
		model.DescribeDomainListRequest{
			DomainKeyword: tea.String("test.com"),
		})
	if err != nil {
		t.Error(err)
	}
	fmt.Println(tea.Prettify(resp))
}

func TestTencentDomain(t *testing.T) {
	resp, err := dnsS.DescribeDomainList("tencent", "ap-shanghai",
		model.DescribeDomainListRequest{
			// DomainKeyword: tea.String(""),
		})
	if err != nil {
		t.Error(err)
	}
	fmt.Println(tea.Prettify(resp))
}

func TestAwsRecord(t *testing.T) {
	resp, err := dnsS.DescribeRecordList("aws", "cn-northwest-1",
		model.DescribeRecordListRequest{
			Domain:     tea.String("test.com"),
			RecordType: tea.String("CNAME"),
		})
	if err != nil {
		t.Error(err)
	}
	fmt.Println(tea.Prettify(resp))
}

func TestTencentRecord(t *testing.T) {
	resp, err := dnsS.DescribeRecordList("tencent", "ap-shanghai",
		model.DescribeRecordListRequest{
			Domain: tea.String("test.com"),
			// RecordType: tea.String("CNAME"),
		})
	if err != nil {
		t.Error(err)
	}
	fmt.Println(tea.Prettify(resp))
}

func TestAwsRecordDetail(t *testing.T) {
	resp, err := dnsS.DescribeRecord("aws", "cn-northwest-1",
		model.DescribeRecordRequest{
			Domain:    tea.String("test.com"),
			SubDomain: tea.String("test"),
		})
	if err != nil {
		t.Error(err)
	}
	fmt.Println(tea.Prettify(resp))
}

func TestTencentRecordDetail(t *testing.T) {
	resp, err := dnsS.DescribeRecord("tencent", "ap-shanghai",
		model.DescribeRecordRequest{
			Domain:    tea.String("test.com"),
			SubDomain: tea.String("tencent1"),
		})
	if err != nil {
		t.Error(err)
	}
	fmt.Println(tea.Prettify(resp))
}

func TestAwsCreateRecord(t *testing.T) {
	resp, err := dnsS.CreateRecord("aws", "cn-northwest-1",
		model.CreateRecordRequest{
			Domain:     tea.String("test.com"),
			SubDomain:  tea.String("test"),
			RecordType: tea.String("A"),
			Value:      tea.String("192.168.3.222"),
			Info:       tea.String("test,中文测试"),
		})
	if err != nil {
		t.Error(err)
	}
	fmt.Println(tea.Prettify(resp))
}

func TestTencentCreateRecord(t *testing.T) {
	resp, err := dnsS.CreateRecord("tencent", "ap-shanghai",
		model.CreateRecordRequest{
			Domain:     tea.String("test.co"),
			SubDomain:  tea.String("tencent"),
			RecordType: tea.String("A"),
			Value:      tea.String("192.168.3.233"),
			Info:       tea.String("test,中文测试"),
		})
	if err != nil {
		t.Error(err)
	}
	fmt.Println(tea.Prettify(resp))
}

func TestAwsUpdateRecord(t *testing.T) {
	resp, err := dnsS.ModifyRecord("aws", "cn-northwest-1",
		model.ModifyRecordRequest{
			Domain:     tea.String("test.com"),
			SubDomain:  tea.String("test123"),
			RecordType: tea.String("A"),
			Value:      tea.String("192.168.3.121"),
		})
	if err != nil {
		t.Error(err)
	}
	fmt.Println(tea.Prettify(resp))
}

func TestTencentUpdateRecord(t *testing.T) {
	resp, err := dnsS.ModifyRecord("tencent", "ap-shanghai",
		model.ModifyRecordRequest{
			Domain:     tea.String("test.com"),
			SubDomain:  tea.String("test"),
			TTL:        tea.Uint64(600),
			RecordType: tea.String("A"),
			Value:      tea.String("192.168.3.235"),
			Status:     tea.Bool(false),
		})
	if err != nil {
		t.Error(err)
	}
	fmt.Println(tea.Prettify(resp))
}

func TestAwsDeleteRecord(t *testing.T) {
	resp, err := dnsS.DeleteRecord("aws", "cn-northwest-1",
		model.DeleteRecordRequest{
			Domain:    tea.String("test.com"),
			SubDomain: tea.String("test123"),
			// RecordType: tea.String("A"),
		})
	if err != nil {
		t.Error(err)
	}
	fmt.Println(tea.Prettify(resp))
}

func TestTencentDeleteRecord(t *testing.T) {
	resp, err := dnsS.DeleteRecord("tencent", "ap-shanghai",
		model.DeleteRecordRequest{
			Domain:     tea.String("test.co"),
			SubDomain:  tea.String("tencent"),
			RecordType: tea.String("A"),
		})
	if err != nil {
		t.Error(err)
	}
	fmt.Println(tea.Prettify(resp))
}
