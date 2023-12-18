package main

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/joho/godotenv"

	"github.com/xops-infra/multi-cloud-sdk/pkg/io"
	"github.com/xops-infra/multi-cloud-sdk/pkg/model"
	server "github.com/xops-infra/multi-cloud-sdk/pkg/service"
)

var serverS *server.ServerService

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
			Regions: strings.Split(os.Getenv("AWS_REGIONS"),
				","),
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

func TestQueryInstances(t *testing.T) {
	startTime := time.Now()
	instances := serverS.QueryInstances(model.InstanceQueryInput{
		Ip: tea.String("10.1.1.1"),
	})
	for _, instance := range instances {
		fmt.Printf("%+v", tea.Prettify(instance))
	}

	fmt.Printf("%s len: %d\n", time.Since(startTime), len(instances))
}

func TestDescribeServers(t *testing.T) {
	instances, err := serverS.DescribeInstances(model.DescribeInstancesInput{
		Profile: "tencent",
		Region:  "ap-shanghai",
		InstanceIds: []*string{
			tea.String("ins-xxx"),
		},
	})
	if err != nil {
		t.Error(err)
		return
	}
	for _, instance := range instances {
		fmt.Printf("%+v", tea.Prettify(instance))
	}
	t.Log("success")
}
