package service_test

import (
	"os"
	"testing"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/joho/godotenv"

	"github.com/xops-infra/multi-cloud-sdk/pkg/io"
	"github.com/xops-infra/multi-cloud-sdk/pkg/model"
	server "github.com/xops-infra/multi-cloud-sdk/pkg/service"
)

var s model.CommonContract

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
	s = server.NewCommonService(profiles, serverAws, serverTencent)
}

// teset tencent get monitor metric data
func TestTencentGetMonitorMetricData(t *testing.T) {
	profile := "tencent"
	region := "ap-shanghai"
	mt := model.MetricsTypeInfrequent
	metricData, err := s.GetMonitorMetricData(profile, region, model.GetMonitorMetricDataRequest{
		InstanceType: model.MetricInstanceTypeCOS,
		MetricsType:  &mt,
		Instances: []string{"xxx-xx",
			"xxx-xx"},
	})
	if err != nil {
		t.Fatalf("failed to get monitor metric data: %v", err)
	}
	t.Logf("monitor metric data: %s", tea.Prettify(metricData))
}
