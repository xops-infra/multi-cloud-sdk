package model

import "fmt"

var (
	ErrCloudNotSupported = fmt.Errorf("profile cloud not supported") // 配置文件中的云不能匹配到
	ErrProfileNotFound   = fmt.Errorf("profile not found")
)
