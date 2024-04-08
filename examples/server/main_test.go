package main

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/joho/godotenv"

	"github.com/xops-infra/multi-cloud-sdk/pkg/io"
	"github.com/xops-infra/multi-cloud-sdk/pkg/model"
	server "github.com/xops-infra/multi-cloud-sdk/pkg/service"
)

var serverS model.InstanceContact

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
	serverS = server.NewServer(profiles, serverAws, serverTencent)
}

func TestDescribeServers(t *testing.T) {
	profile := "tencent"
	region := "ap-shanghai"
	timeStart := time.Now()
	filter := model.InstanceFilter{
		Size: tea.Int64(10),
		// Profile: tea.String("aws"),
		// Region:  tea.String("cn-northwest-1"),
	}
	// filter.Owner = tea.String("zhoushoujian")
	// filter.NextMarker = tea.String("xxx")
	var ids []model.Instance
	filter.Status = model.InstanceStatusRunning.TString()
	resp, err := serverS.DescribeInstances(profile, region, filter)
	if err != nil {
		t.Error(err)
		return
	}
	ids = append(ids, resp.Instances...)
	for {
		fmt.Println(tea.Prettify(filter))
		resp, err = serverS.DescribeInstances(profile, region, filter)
		if err != nil {
			t.Error(err)
			return
		}
		if resp.NextMarker == nil {
			break
		}
		filter.NextMarker = resp.NextMarker
		ids = append(ids, resp.Instances...)
	}

	t.Log(time.Since(timeStart), len(ids))
}
