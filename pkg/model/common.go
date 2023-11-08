package model

import "time"

type CommonQueryInput struct {
	ID            string        `json:"id"`
	CloudProvider CloudProvider `json:"cloud_provider"`
	Region        string        `json:"region"`
	Account       string        `json:"account"`
}

// 依据账号和区域过滤，完全符合条件返回true
func (i CommonQueryInput) Filter(profile, region string) bool {
	if i.Region != "" && i.Region != region {
		return false
	}
	if i.Account != "" && i.Account != profile {
		return false
	}
	return true
}

// 按照`ISO8601`标准表示，并且使用`UTC`时间。格式为：`YYYY-MM-DDThh:mm:ssZ` to time.Time
func TimeParse(t string) (time.Time, error) {
	return time.Parse(time.RFC3339, t)
}
