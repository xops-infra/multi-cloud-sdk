package model

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
