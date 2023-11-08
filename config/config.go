package config

import (
	"github.com/xops-infra/multi-cloud-sdk/pkg/model"
)

type ServerConfig struct {
	Profiles []model.ProfileConfig `mapstructure:"profiles"`
}
