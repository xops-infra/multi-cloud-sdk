package main

import (
	"fmt"
	"os"
	"time"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/xops-infra/multi-cloud-sdk/pkg/io"
	"github.com/xops-infra/multi-cloud-sdk/pkg/model"
	server "github.com/xops-infra/multi-cloud-sdk/pkg/service"
)

var serverS *server.ServerService

func init() {
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
	if profiles[0].AK == "" {
		panic("AWS_ACCESS_KEY_ID not found")
	}
	cloudIo := io.NewCloudClient(profiles)
	serverTencent := io.NewTencentClient(cloudIo)
	serverAws := io.NewAwsClient(cloudIo)
	serverS = server.NewServer(profiles, serverAws, serverTencent)
}

func main() {
	startTime := time.Now()
	instances := serverS.QueryInstances(model.InstanceQueryInput{
		// Status: model.InstanceStatusRunning,
		// Name: "as-tke-np-5qiueryt",
		// Ip: "10.",
		Ip: "10.40.40.",
		// Owner: "zhoushoujian",
		// Account: "aws",
	})
	for _, instance := range instances {
		fmt.Println(tea.Prettify(instance))
	}
	fmt.Printf("%s len: %d\n", time.Since(startTime), len(instances))
}
