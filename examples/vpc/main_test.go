package main

import (
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"

	"github.com/xops-infra/multi-cloud-sdk/pkg/io"
	"github.com/xops-infra/multi-cloud-sdk/pkg/model"
	server "github.com/xops-infra/multi-cloud-sdk/pkg/service"
)

var vpcS model.CommonContract

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
	cloudIo := io.NewCloudClient(profiles)
	serverTencent := io.NewTencentClient(cloudIo)
	serverAws := io.NewAwsClient(cloudIo)
	vpcS = server.NewCommonService(profiles, serverAws, serverTencent)
}

func TestAws(t *testing.T) {
	{
		startTime := time.Now()
		vpcs, err := vpcS.QueryVPCs("aws", "cn-northwest-1", model.CommonFilter{})
		if err != nil {
			t.Error(err)
		}
		t.Log("aws vpc test ok", time.Since(startTime), len(vpcs))
	}
	{
		startTime := time.Now()
		eips, err := vpcS.QueryEIPs("aws", "cn-northwest-1", model.CommonFilter{})
		if err != nil {
			t.Error(err)
		}
		t.Log("aws eip test ok", time.Since(startTime), len(eips))
	}
	{
		startTime := time.Now()
		nats, err := vpcS.QueryNATs("aws", "cn-northwest-1", model.CommonFilter{})
		if err != nil {
			t.Error(err)
		}
		t.Log("aws nat test ok", time.Since(startTime), len(nats))
	}
	{
		startTime := time.Now()
		subnets, err := vpcS.QuerySubnets("aws", "cn-northwest-1", model.CommonFilter{})
		if err != nil {
			t.Error(err)
		}
		t.Log("aws subnets test ok", time.Since(startTime), len(subnets))
	}
}

func TestTencent(t *testing.T) {
	{
		startTime := time.Now()
		vpcs, err := vpcS.QueryVPCs("tencent", "ap-shanghai", model.CommonFilter{})
		if err != nil {
			t.Error(err)
		}
		t.Log("tencent vpc test ok", time.Since(startTime), len(vpcs))
	}
	{
		startTime := time.Now()
		eips, err := vpcS.QueryEIPs("tencent", "ap-shanghai", model.CommonFilter{})
		if err != nil {
			t.Error(err)
		}
		t.Log("tencent eip test ok", time.Since(startTime), len(eips))
	}
	{
		startTime := time.Now()
		nats, err := vpcS.QueryNATs("tencent", "ap-shanghai", model.CommonFilter{})
		if err != nil {
			t.Error(err)
		}
		t.Log("tencent nat test ok", time.Since(startTime), len(nats))
	}
	{
		startTime := time.Now()
		subnets, err := vpcS.QuerySubnets("tencent", "ap-shanghai", model.CommonFilter{})
		if err != nil {
			t.Error(err)
		}
		t.Log("tencent subnets test ok", time.Since(startTime), len(subnets))
	}
}
