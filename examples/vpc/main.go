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

var vpcS model.VpcContract

func init() {
	profiles := []model.ProfileConfig{
		{
			Name:  "us1549",
			Cloud: model.AWS,
			AK:    os.Getenv("AWS_ACCESS_KEY_ID"),
			SK:    os.Getenv("AWS_SECRET_ACCESS_KEY"),
			Regions: []string{
				"us-east-2",
			},
		},
		// {
		// 	Name:  "tencent",
		// 	Cloud: model.TENCENT,
		// 	AK:    os.Getenv("TENCENT_ACCESS_KEY"),
		// 	SK:    os.Getenv("TENCENT_SECRET_KEY"),
		// 	Regions: []string{
		// 		"ap-shanghai",
		// 		// "na-ashburn",
		// 	},
		// },
	}
	if profiles[0].AK == "" {
		panic("AWS_ACCESS_KEY_ID not found")
	}
	cloudIo := io.NewCloudClient(profiles)
	serverTencent := io.NewTencentClient(cloudIo)
	serverAws := io.NewAwsClient(cloudIo)
	vpcS = server.NewVpcService(profiles, serverAws, serverTencent)
}

func main() {
	startTime := time.Now()
	fmt.Println("vpcs...")
	vpcs, _ := vpcS.QueryVPCs(model.CommonQueryInput{
		CloudProvider: model.TENCENT,
	})
	for _, vpc := range vpcs {
		fmt.Println(tea.Prettify(vpc))
	}
	fmt.Println("eips...")
	eips, _ := vpcS.QueryEIPs(model.CommonQueryInput{})
	for _, eip := range eips {
		fmt.Println(tea.Prettify(eip))
	}
	fmt.Println("nats...")
	nats, _ := vpcS.QueryNATs(model.CommonQueryInput{})
	for _, nat := range nats {
		fmt.Println(tea.Prettify(nat))
	}
	fmt.Println("subnets...")
	subnets, _ := vpcS.QuerySubnets(model.CommonQueryInput{})
	for _, subnet := range subnets {
		fmt.Println(tea.Prettify(subnet))
	}

	fmt.Printf("%s \n", time.Since(startTime))
}
