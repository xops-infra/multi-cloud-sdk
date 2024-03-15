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
		},
		{
			Name:  "tencent",
			Cloud: model.TENCENT,
			AK:    os.Getenv("TENCENT_ACCESS_KEY"),
			SK:    os.Getenv("TENCENT_SECRET_KEY"),
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
	resp, err := dnsS.DescribeDomainList("aws",
		model.DescribeDomainListRequest{
			DomainKeyword: tea.String(os.Getenv("TEST_AWS_DOMAIN")),
		})
	if err != nil {
		t.Error(err)
	}
	fmt.Println(tea.Prettify(resp))
}

func TestTencentDomain(t *testing.T) {
	resp, err := dnsS.DescribeDomainList("tencent",
		model.DescribeDomainListRequest{
			// DomainKeyword: tea.String(""),
		})
	if err != nil {
		t.Error(err)
	}
	fmt.Println(tea.Prettify(resp))
}

func TestAwsRecord(t *testing.T) {
	resp, err := dnsS.DescribeRecordList("aws",
		model.DescribeRecordListRequest{
			Domain: tea.String(os.Getenv("TEST_AWS_DOMAIN")),
			Limit:  tea.Int64(2),
			// RecordType: tea.String("CNAME"),
			NextMarker: tea.String("cGF0c25hcC5pbmZvLixUWFR0aGlzaXNhc2NyZWF0a2V5"),
		})
	if err != nil {
		t.Error(err)
	}
	fmt.Println(tea.Prettify(resp.NextMarker))
}

func TestTencentRecord(t *testing.T) {
	resp, err := dnsS.DescribeRecordList("tencent",
		model.DescribeRecordListRequest{
			Domain:     tea.String(os.Getenv("TEST_TENCENT_DOMAIN")),
			Limit:      tea.Int64(2),
			NextMarker: tea.String("M3RoaXNpc2FzY3JlYXRrZXk="),
		})
	if err != nil {
		t.Error(err)
	}
	fmt.Println(tea.Prettify(resp))
}

func TestAwsRecordDetail(t *testing.T) {
	resp, err := dnsS.DescribeRecord("aws",
		model.DescribeRecordRequest{
			Domain:    tea.String(os.Getenv("TEST_AWS_DOMAIN")),
			SubDomain: tea.String("test"),
		})
	if err != nil {
		t.Error(err)
	}
	fmt.Println(tea.Prettify(resp))
}

func TestTencentRecordDetail(t *testing.T) {
	resp, err := dnsS.DescribeRecord("tencent",
		model.DescribeRecordRequest{
			Domain:    tea.String(os.Getenv("TEST_TENCENT_DOMAIN")),
			SubDomain: tea.String("tencent1"),
		})
	if err != nil {
		t.Error(err)
	}
	fmt.Println(tea.Prettify(resp))
}

func TestAwsCreateRecord(t *testing.T) {
	resp, err := dnsS.CreateRecord("aws",
		model.CreateRecordRequest{
			Domain:     tea.String(os.Getenv("TEST_AWS_DOMAIN")),
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
	resp, err := dnsS.CreateRecord("tencent",
		model.CreateRecordRequest{
			Domain:     tea.String(os.Getenv("TEST_TENCENT_DOMAIN")),
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
	err := dnsS.ModifyRecord("aws", true,
		model.ModifyRecordRequest{
			Domain:     tea.String(os.Getenv("TEST_AWS_DOMAIN")),
			SubDomain:  tea.String("test123"),
			RecordType: tea.String("A"),
			Value:      tea.String("192.168.3.121"),
		})
	if err != nil {
		t.Error(err)
	}
}

func TestTencentUpdateRecord(t *testing.T) {
	err := dnsS.ModifyRecord("tencent", true,
		model.ModifyRecordRequest{
			Domain:     tea.String(os.Getenv("TEST_TENCENT_DOMAIN")),
			SubDomain:  tea.String("test"),
			TTL:        tea.Uint64(600),
			RecordType: tea.String("A"),
			Value:      tea.String("192.168.3.235"),
			Status:     tea.Bool(false),
		})
	if err != nil {
		t.Error(err)
	}
}

func TestAwsDeleteRecord(t *testing.T) {
	resp, err := dnsS.DeleteRecord("aws",
		model.DeleteRecordRequest{
			Domain:    tea.String(os.Getenv("TEST_AWS_DOMAIN")),
			SubDomain: tea.String("test123"),
			// RecordType: tea.String("A"),
		})
	if err != nil {
		t.Error(err)
	}
	fmt.Println(tea.Prettify(resp))
}

func TestTencentDeleteRecord(t *testing.T) {
	resp, err := dnsS.DeleteRecord("tencent",
		model.DeleteRecordRequest{
			Domain:     tea.String(os.Getenv("TEST_TENCENT_DOMAIN")),
			SubDomain:  tea.String("tencent"),
			RecordType: tea.String("A"),
		})
	if err != nil {
		t.Error(err)
	}
	fmt.Println(tea.Prettify(resp))
}
