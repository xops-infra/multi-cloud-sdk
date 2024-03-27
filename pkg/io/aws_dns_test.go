package io_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/stretchr/testify/assert"
	"github.com/xops-infra/multi-cloud-sdk/pkg/model"
)

// TEST ALL
func TestMain(t *testing.T) {
	profile := "aws"
	region := "cn-northwest-1"
	// TEST DescribePrivateDomainList
	{
		req := model.DescribeRecordListWithPageRequest{
			Limit:  tea.Int64(10),
			Page:   tea.Int64(1),
			Domain: tea.String(os.Getenv("TEST_AWS_DOMAIN")),
		}
		_, err := AwsIo.DescribeRecordListWithPages(
			profile,
			region,
			req,
		)
		assert.Nil(t, err)
	}
	// TEST DescribeRecordList
	{
		req := model.DescribeRecordListRequest{
			Domain:  tea.String(os.Getenv("TEST_AWS_DOMAIN")),
			Keyword: tea.String("pop"),
		}
		resp, err := AwsIo.DescribeRecordList(profile, region, req)
		if err != nil {
			t.Error(err)
			return
		}
		assert.Nil(t, err)
		fmt.Println(tea.Prettify(resp))
	}
	// TEST DescribeRecord
	{
		_, err := AwsIo.DescribeRecord(profile, region, model.DescribeRecordRequest{
			Domain:     tea.String(os.Getenv("TEST_AWS_DOMAIN")),
			SubDomain:  tea.String("test"),
			RecordType: tea.String("CNAME"),
		})
		assert.Nil(t, err)
	}
}

// TEST DescribePrivateDomainList
func TestDescribeList(t *testing.T) {
	// TEST DescribeRecordList
	{
		req := model.DescribeRecordListRequest{
			Domain: tea.String("patsnap.info"),
		}
		resp, err := AwsIo.DescribeRecordList("aws", "cn-notrhwest-1", req)
		if err != nil {
			t.Error(err)
			return
		}
		assert.Nil(t, err)
		fmt.Println(tea.Prettify(resp))
	}
}

// TEST DescribePrivateDomainList
func TestDescribeListWithPages(t *testing.T) {
	// TEST DescribeRecordList
	req := model.DescribeDomainListRequest{}
	resp, err := AwsIo.DescribeDomainList("aws-cas", "us-east-2", req)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(tea.Prettify(resp))
}
