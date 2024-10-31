package model

import (
	"encoding/xml"
	"fmt"

	"github.com/alibabacloud-go/tea/tea"
	cos "github.com/tencentyun/cos-go-sdk-v5"
)

type Bucket struct {
	Name       string `json:"name"`
	CreateTime string `json:"create_time"` // "2006-01-02 15:04:05" localTime
	Location   string `json:"location"`
	Tags       Tags   `json:"tags"`
	// Meta any    `json:"meta"`
}

type CreateBucketRequest struct {
	BucketName *string `json:"name" binding:"required"` // 腾讯云必须带上 appid， examplebucket-1250000000
	Tags       Tags    `json:"tags"`
}

type DeleteBucketRequest struct {
	BucketName *string `json:"name" binding:"required"`
}

type DeleteBucketResponse struct {
	Meta any `json:"meta"`
}

type ListBucketRequest struct {
	KeyWord *string `json:"keyword"`
}

type ListBucketResponse struct {
	Buckets []*Bucket `json:"buckets"`
	Total   int64     `json:"total"`
}

type ObjectPregisnRequest struct {
	Bucket *string `json:"bucket" binding:"required"`
	Key    *string `json:"key" binding:"required"`
	Expire *int64  `json:"expire"` // 默认 1 小时。 签名最多支持7天(604800秒)，控制台上最多 12小时(43200秒)
}

type ObjectPregisnResponse struct {
	Url string `json:"url"`
}

type CreateBucketLifecycleRequest struct {
	Bucket     *string
	Lifecycles []Lifecycle
}

type Lifecycle struct {
	ID                             *string                                  `xml:"ID" json:"id"`
	Filter                         *LifecycleFilter                         `xml:"Filter" json:"filter"`
	Status                         *bool                                    `xml:"Status" json:"status"`
	Expiration                     *LifecycleExpiration                     `xml:"Expiration" json:"expiration"`
	Transition                     *LifecycleTransition                     `xml:"Transition" json:"transition"`
	AbortIncompleteMultipartUpload *LifecycleAbortIncompleteMultipartUpload `xml:"AbortIncompleteMultipartUpload" json:"abort_incomplete_multipart_upload"`
}

type LifecycleFilter struct {
	Prefix *string
}

type LifecycleExpiration struct {
	Days *int
}

type LifecycleTransition struct {
	Days *int
}

type LifecycleAbortIncompleteMultipartUpload struct {
	DaysAfterInitiation *int
}

/*
<LifecycleConfiguration>
	<Rule>
		<ID>huggingface模型定期删除</ID>
		<Filter>
			<Prefix>hg/</Prefix>
		</Filter>
		<Status>Enabled</Status>
		<Expiration>
			<Days>5</Days>
		</Expiration>
	</Rule>
	<Rule>
		<ID>OPS_BASE</ID>
		<Filter/>
		<Status>Enabled</Status>
		<AbortIncompleteMultipartUpload>
			<DaysAfterInitiation>30</DaysAfterInitiation>
		</AbortIncompleteMultipartUpload>
	</Rule>
</LifecycleConfiguration>
*/
// to COSLifecycle
func (c *CreateBucketLifecycleRequest) ToCOSLifecycle() *cos.BucketPutLifecycleOptions {
	cosRules := make([]cos.BucketLifecycleRule, len(c.Lifecycles))

	for i, lifecycle := range c.Lifecycles {
		rule := cos.BucketLifecycleRule{}

		if lifecycle.ID != nil {
			rule.ID = *lifecycle.ID
		}
		if lifecycle.Filter != nil && lifecycle.Filter.Prefix != nil {
			rule.Filter = &cos.BucketLifecycleFilter{
				Prefix: *lifecycle.Filter.Prefix,
			}
		}
		if lifecycle.Status != nil || *lifecycle.Status {
			rule.Status = "Enabled" // 默认开启
		} else {
			rule.Status = "Disabled"
		}
		if lifecycle.AbortIncompleteMultipartUpload != nil {
			rule.AbortIncompleteMultipartUpload = &cos.BucketLifecycleAbortIncompleteMultipartUpload{
				DaysAfterInitiation: *lifecycle.AbortIncompleteMultipartUpload.DaysAfterInitiation,
			}
		}

		if lifecycle.Expiration != nil {
			rule.Expiration = &cos.BucketLifecycleExpiration{
				Days: *lifecycle.Expiration.Days,
			}
		}
		if lifecycle.Transition != nil {
			rule.Transition = []cos.BucketLifecycleTransition{
				{
					Days: *lifecycle.Transition.Days,
				},
			}
		}

		cosRules[i] = rule
	}
	fmt.Println(tea.Prettify(cosRules))
	return &cos.BucketPutLifecycleOptions{
		XMLName: xml.Name{Local: "LifecycleConfiguration"},
		Rules:   cosRules,
	}
}
