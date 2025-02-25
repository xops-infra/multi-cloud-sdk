package io

import (
	"fmt"

	"github.com/xops-infra/multi-cloud-sdk/pkg/model"
)

type tencentClient struct {
	io model.ClientIo
}

func NewTencentClient(io model.ClientIo) model.CloudIO {
	return &tencentClient{
		io: io,
	}
}

func (c *tencentClient) CreateSqs(profile, region string, input model.CreateSqsRequest) error {
	return fmt.Errorf("not implemented")
}
