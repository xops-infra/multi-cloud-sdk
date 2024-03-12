package io_test

import (
	"os"
	"testing"
	"time"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/xops-infra/multi-cloud-sdk/pkg/model"
)

// TEST DescribeDomainList
func TestDescribeDomainList(t *testing.T) {
	timeStart := time.Now()
	resp, err := TencentIo.DescribeDomainList("tencent", model.DescribeDomainListRequest{
		DomainKeyword: tea.String(os.Getenv("TEST_TENCENT_DOMAIN")),
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("Success.", time.Since(timeStart), tea.Prettify(resp))
}

// TEST DescribeRecordList
func TestDescribeRecordList(t *testing.T) {
	timeStart := time.Now()
	resp, err := TencentIo.DescribeRecordList("tencent", model.DescribeRecordListRequest{
		Limit:      tea.Int64(2),
		Domain:     tea.String(os.Getenv("TEST_TENCENT_DOMAIN")),
		NextMarker: tea.String("MnRoaXNpc2FzY3JlYXRrZXk="),
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("Success.", time.Since(timeStart), tea.Prettify(resp))
}

// TEST DescribeRecord
func TestDescribeRecord(t *testing.T) {
	timeStart := time.Now()
	record, err := TencentIo.DescribeRecord("tencent", model.DescribeRecordRequest{
		Domain: tea.String(os.Getenv("TEST_TENCENT_DOMAIN")),
		// SubDomain: tea.String("test"),
		// RecordType: tea.String("CNAME"),
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("Success.", time.Since(timeStart), tea.Prettify(record))
}

// TEST CreateRecord
func TestCreateRecord(t *testing.T) {
	timeStart := time.Now()
	subDomain := "abcde"
	{
		resp, err := TencentIo.CreateRecord("tencent", model.CreateRecordRequest{
			Domain:     tea.String(os.Getenv("TEST_TENCENT_DOMAIN")),
			SubDomain:  tea.String(subDomain),
			RecordType: tea.String("CNAME"),
			Value:      tea.String("test.com"),
		})
		if err != nil {
			t.Error(err)
			return
		}
		t.Log("Success.", time.Since(timeStart), tea.Prettify(resp))
	}
	{
		// ModifyRecord
		err := TencentIo.ModifyRecord("tencent", false, model.ModifyRecordRequest{
			Domain:     tea.String(os.Getenv("TEST_TENCENT_DOMAIN")),
			SubDomain:  tea.String(subDomain),
			RecordType: tea.String("CNAME"),
			Value:      tea.String("test.com"),
		})
		if err != nil {
			t.Error(err)
			return
		}
	}
	{
		// DeleteRecord
		resp, err := TencentIo.DeleteRecord("tencent", model.DeleteRecordRequest{
			Domain:     tea.String(os.Getenv("TEST_TENCENT_DOMAIN")),
			SubDomain:  tea.String(subDomain),
			RecordType: tea.String("CNAME"),
		})
		if err != nil {
			t.Error(err)
			return
		}
		t.Log("Success.", time.Since(timeStart), tea.Prettify(resp))
	}

}
