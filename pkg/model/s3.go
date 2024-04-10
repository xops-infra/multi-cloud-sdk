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
