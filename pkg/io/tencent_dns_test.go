package io_test

import (
	"os"
	"testing"
	"time"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/stretchr/testify/assert"
	"github.com/xops-infra/multi-cloud-sdk/pkg/model"
)

var profile = "tencent"
var region = "" // 腾讯云不需要region

// TEST DescribeDomainList
func TestDescribeDomainList(t *testing.T) {
	timeStart := time.Now()
	resp, err := TencentIo.DescribeDomainList("tencent", "", model.DescribeDomainListRequest{
		DomainKeyword: tea.String(os.Getenv("TEST_TENCENT_DOMAIN")),
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("Success.", time.Since(timeStart), tea.Prettify(resp))
}

func TestDns(t *testing.T) {

	{
		// TEST DescribeDomainList
		_, err := TencentIo.DescribeDomainList(profile, region, model.DescribeDomainListRequest{})
		if err != nil {
			t.Error(err)
			return
		}
		assert.Nil(t, err)
	}

	{
		// TEST DescribeRecordList
		_, err := TencentIo.DescribeRecordList(profile, region, model.DescribeRecordListRequest{
			Domain:  tea.String(os.Getenv("TEST_TENCENT_DOMAIN")),
			Keyword: tea.String("test"),
		})
		assert.Nil(t, err)
	}

	{
		// TEST DescribeRecordList
		_, err := TencentIo.DescribeRecordListWithPages(profile, region, model.DescribeRecordListWithPageRequest{
			Domain: tea.String(os.Getenv("TEST_TENCENT_DOMAIN")),
			Limit:  tea.Int64(2),
			Page:   tea.Int64(2),
		})
		assert.Nil(t, err)
	}
}

func TestList(t *testing.T) {
	// TEST DescribeRecordList
	resp, err := TencentIo.DescribeRecordList(profile, region, model.DescribeRecordListRequest{
		Domain:  tea.String(os.Getenv("TEST_TENCENT_DOMAIN")),
		Keyword: tea.String("test"),
	})
	assert.Nil(t, err)
	t.Log(tea.Prettify(resp))
}

func TestCreate(t *testing.T) {
	// TEST CreateRecord
	resp, err := TencentIo.CreateRecord(profile, region, model.CreateRecordRequest{
		Domain:     tea.String(os.Getenv("TEST_TENCENT_DOMAIN")),
		SubDomain:  tea.String("zsj.test"),
		RecordType: tea.String("A"),
		Value:      tea.String("1.1.1.1"),
	})
	assert.Nil(t, err)
	t.Log(tea.Prettify(resp))
}

func TestDelete(t *testing.T) {
	// TEST DeleteRecord
	_, err := TencentIo.DeleteRecord(profile, region, model.DeleteRecordRequest{
		Domain:     tea.String(os.Getenv("TEST_TENCENT_DOMAIN")),
		SubDomain:  tea.String("zsj.test"),
		RecordType: tea.String("A"),
	})
	assert.Nil(t, err)
}
