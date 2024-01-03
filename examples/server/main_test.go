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
				"na-ashburn",
			},
		},
	}
	cloudIo := io.NewCloudClient(profiles)
	serverTencent := io.NewTencentClient(cloudIo)
	serverAws := io.NewAwsClient(cloudIo)
	serverS = server.NewServer(profiles, serverAws, serverTencent)
}

func TestDescribeServers(t *testing.T) {
	timeStart := time.Now()
	filter := model.InstanceFilter{
		// Profile: tea.String("aws"),
		// Region:  tea.String("cn-northwest-1"),
	}
	// filter.Owner = tea.String("zhoushoujian")
	filter.Status = model.InstanceStatusRunning.TString()
	instances, err := serverS.DescribeInstances("aws", "cn-northwest-1", filter)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(time.Since(timeStart), len(instances.Instances))
}
