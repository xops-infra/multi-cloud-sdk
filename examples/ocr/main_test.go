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

var ocrS model.OcrContract

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	profiles := []model.ProfileConfig{
		{
			Name:  "tencent",
			Cloud: model.TENCENT,
			AK:    os.Getenv("TENCENT_ACCESS_KEY"),
			SK:    os.Getenv("TENCENT_SECRET_KEY"),
		},
	}
	if profiles[0].AK == "" {
		panic("Key not found")
	}
	cloudIo := io.NewCloudClient(profiles)
	serverTencent := io.NewTencentClient(cloudIo)
	serverAws := io.NewAwsClient(cloudIo)
	ocrS = server.NewOcrService(profiles, serverAws, serverTencent)
}

func TestOcr(t *testing.T) {
	resp, err := ocrS.QueryOcr("tencent", "ap-shanghai", model.OcrRequest{
		ImageUrl: tea.String("https://dayou.zhuangzhaimuye.com/wp-content/uploads/2021/03/%E7%99%BD%E6%A1%A6%E6%9C%A8-1%EF%BC%88DY199%EF%BC%89-1140x1747.jpg"),
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(tea.Prettify(resp))
}

func TestCreatePicture(t *testing.T) {
	resp, err := ocrS.CreatePicture("tencent", "ap-shanghai", model.CreatePictureRequest{
		GroupId:  tea.String("common"),
		EntityId: tea.String("test_entity"),
		PicName:  tea.String("test"),
		ImageUrl: tea.String("https://dayou.zhuangzhaimuye.com/wp-content/uploads/2021/03/%E7%99%BD%E6%A1%A6%E6%9C%A8-1%EF%BC%88DY199%EF%BC%89-1140x1747.jpg"),
		Tags:     tea.String(`{"time": "2021-03-25"}`),
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(tea.Prettify(resp))
}

func Test_getPictureByName(t *testing.T) {
	resp, err := ocrS.GetPictureByName("tencent", "ap-shanghai", model.CommonPictureRequest{
		GroupId:  tea.String("common"),
		EntityId: tea.String("test_entity"),
		PicName:  tea.String("test"),
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(tea.Prettify(resp))
}

func queryPicture() {
	resp, err := ocrS.QueryPicture("tencent", "ap-shanghai", model.QueryPictureRequest{
		GroupId: tea.String("common"),
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(tea.Prettify(resp))
}

func Test_deletePicture(t *testing.T) {
	resp, err := ocrS.DeletePicture("tencent", "ap-shanghai", model.CommonPictureRequest{
		GroupId:  tea.String("common"),
		EntityId: tea.String("test_entity"),
		PicName:  tea.String("ocr"),
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(tea.Prettify(resp))
}

func Test_updatePicture(t *testing.T) {
	tags := fmt.Sprintf(`{"%s": "%d"}`, "time", time.Now().Unix())
	resp, err := ocrS.UpdatePicture("tencent", "ap-shanghai", model.UpdatePictureRequest{
		GroupId:  tea.String("common"),
		EntityId: tea.String("test_entity"),
		PicName:  tea.String("test"),
		Tags:     tea.String(tags),
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(tea.Prettify(resp))
}

func Test_searchPicture(t *testing.T) {
	resp, err := ocrS.SearchPicture("tencent", "ap-shanghai", model.SearchPictureRequest{
		GroupId:  tea.String("common"),
		ImageUrl: tea.String("https://dayou.zhuangzhaimuye.com/wp-content/uploads/2021/03/%E7%99%BD%E6%A1%A6%E6%9C%A8-1%EF%BC%88DY199%EF%BC%89-1140x1747.jpg"),
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(tea.Prettify(resp))
}
