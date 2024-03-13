package io_test

import (
	"testing"
	"time"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/xops-infra/multi-cloud-sdk/pkg/model"
)

// TEST DescribePrivateDomainList
func TestDescribePrivateDomainList(t *testing.T) {
	timeStart := time.Now()
	resp, err := TencentIo.DescribePrivateDomainList("tencent", model.DescribeDomainListRequest{
		DomainKeyword: tea.String("com"),
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("Success.", time.Since(timeStart), tea.Prettify(resp))
}

// TEST DescribePrivateRecordList
func TestDescribePrivateRecordList(t *testing.T) {
	timeStart := time.Now()
	resp, err := TencentIo.DescribePrivateRecordList("tx-dev", model.DescribeRecordListRequest{
		Domain:     tea.String("domain.com"),
		Limit:      tea.Int64(6),
		NextMarker: tea.String("M3RoaXNpc2FzY3JlYXRrZXk="),
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("Success.", time.Since(timeStart), tea.Prettify(resp))
}

// TEST CreatePrivateRecord
func TestCreatePrivateRecord(t *testing.T) {
	timeStart := time.Now()
	resp, err := TencentIo.CreatePrivateRecord("tencent", model.CreateRecordRequest{
		Domain:     tea.String("domain.com"),
		SubDomain:  tea.String("zsj1"),
		RecordType: tea.String("A"),
		Value:      tea.String("192.168.1.1"),
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("Success.", time.Since(timeStart), tea.Prettify(resp))
}

// TEST ModifyPrivateRecord
func TestModifyPrivateRecord(t *testing.T) {
	timeStart := time.Now()
	err := TencentIo.ModifyPrivateRecord("tencent", model.ModifyRecordRequest{
		Domain:     tea.String("domain.com"),
		RecordId:   tea.Uint64(1965530),
		SubDomain:  tea.String("zsj"),
		RecordType: tea.String("A"),
		Value:      tea.String("1.1.1.1"), // new value
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("Success.", time.Since(timeStart))
}

// TEST DeletePrivateRecord
func TestDeletePrivateRecord(t *testing.T) {
	timeStart := time.Now()
	err := TencentIo.DeletePrivateRecord("tencent", model.DeletePrivateRecordRequest{
		Domain:   tea.String("domain.com"),
		RecordId: tea.String("1965587"),
		// RecordIds: []*string{tea.String("1965530")},
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("Success.", time.Since(timeStart))
}
