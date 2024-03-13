package io_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/xops-infra/multi-cloud-sdk/pkg/model"
)

// TEST DescribeRecordList
func TestDescribeAWSRecordList(t *testing.T) {
	req := model.DescribeRecordListRequest{
		// Limit:      tea.Int64(2),
		Domain: tea.String(os.Getenv("TEST_AWS_DOMAIN")),
		// NextMarker: tea.String("itbtZa/gGkn9H97wBqpq3fO8S4bgQitmCJirgIFR7BSR"),
	}
	fmt.Println(tea.Prettify(req))
	resp, err := AwsIo.DescribeRecordList(
		"aws",
		req,
	)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(tea.Prettify(resp))
}

// TEST DescribeRecordList
func TestDescribeTencentRecordList(t *testing.T) {
	req := model.DescribeRecordListRequest{
		Limit:   tea.Int64(2),
		Domain:  tea.String(os.Getenv("TEST_AWS_DOMAIN")),
		Keyword: tea.String("pop"),
		// NextMarker: tea.String("Ki5wYXRzbmFwLnByb2QuZXUuYXBpLXRlc3RpbmcucGF0c25hcC5pbmZvLixDTkFNRXRoaXNpc2FzY3JlYXRrZXk="),
	}
	fmt.Println(tea.Prettify(req))
	resp, err := AwsIo.DescribeRecordList(
		"aws",
		req,
	)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(tea.Prettify(resp))
}

// TEST CreateRecord
func TestCreateAWSRecord(t *testing.T) {
	req := model.CreateRecordRequest{
		Domain:     tea.String(os.Getenv("TEST_AWS_DOMAIN")),
		SubDomain:  tea.String("testA"),
		RecordType: tea.String("CNAME"),
		Value:      tea.String("test.com"),
	}
	resp, err := AwsIo.CreateRecord(
		"aws",
		req,
	)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(tea.Prettify(resp))
}

// TEST DescribeRecord
func TestDescribeAWSRecord(t *testing.T) {
	req := model.DescribeRecordRequest{
		Domain:     tea.String(os.Getenv("TEST_AWS_DOMAIN")),
		SubDomain:  tea.String("testA"),
		RecordType: tea.String("A"),
	}
	resp, err := AwsIo.DescribeRecord(
		"aws",
		req,
	)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(tea.Prettify(resp))
}

// TEST ModifyRecord
func TestModifyAWSRecord(t *testing.T) {
	req := model.ModifyRecordRequest{
		Domain:     tea.String(os.Getenv("TEST_AWS_DOMAIN")),
		SubDomain:  tea.String("test"),
		RecordType: tea.String("A"),
		Value:      tea.String("192.168.1.1"), // 修改记录值
	}
	err := AwsIo.ModifyRecord("aws", false, req)
	if err != nil {
		t.Error(err)
		return
	}
}

// TEST DeleteRecord
func TestDeleteAWSRecord(t *testing.T) {
	req := model.DeleteRecordRequest{
		Domain:     tea.String(os.Getenv("TEST_AWS_DOMAIN")),
		SubDomain:  tea.String("test"),
		RecordType: tea.String("A"),
	}
	_, err := AwsIo.DeleteRecord("aws", req)
	if err != nil {
		t.Error(err)
		return
	}
}
