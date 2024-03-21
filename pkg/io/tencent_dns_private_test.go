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
	resp, err := TencentIo.DescribePrivateDomainList("tencent", model.DescribeDomainListRequest{})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("Success.", time.Since(timeStart), tea.Prettify(resp))
}

// TEST DescribePrivateRecordList
func TestDescribePrivateRecordList(t *testing.T) {
	timeStart := time.Now()
	resp, err := TencentIo.DescribePrivateRecordList("tencent", model.DescribePrivateRecordListRequest{
		Domain:  tea.String("zone-3m8hlc6o"),
		Keyword: tea.String("zsj"),
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("Success.", time.Since(timeStart), tea.Prettify(resp))
}

// TEST DescribePrivateRecordListWithPages
func TestDescribePrivateRecordListWithPages(t *testing.T) {
	timeStart := time.Now()
	resp, err := TencentIo.DescribePrivateRecordListWithPages("tencent", model.DescribePrivateDnsRecordListWithPageRequest{
		Domain: tea.String("zone-3m8hlc6o"),
		// Limit:  tea.Int64(2),
		// Page:   tea.Int64(2),
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
		Domain:     tea.String("zone-3m8hlc6o"),
		SubDomain:  tea.String("zsj11"),
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
