package model

import (
	"time"
)

const (
	slot = "thisisascreatkey"
)

type CommonFilter struct {
	ID string `json:"id"`
}

// 按照`ISO8601`标准表示，并且使用`UTC`时间。格式为：`YYYY-MM-DDThh:mm:ssZ` to time.Time
func TimeParse(t string) (time.Time, error) {
	return time.Parse(time.RFC3339, t)
}
