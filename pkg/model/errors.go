package model

import "fmt"

var (
	ErrCloudNotSupported = fmt.Errorf("cloud not supported")
	ErrProfileNotFound   = fmt.Errorf("profile not found")
)
