package io

import (
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
