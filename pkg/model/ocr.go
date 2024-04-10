package model

/*
https://console.cloud.tencent.com/api/explorer?Product=tiia&Version=2019-05-29&Action=SearchImage
*/

type OcrResponse struct {
	TextDetections []TextDetection `json:"text_detections"`
}

type TextDetection struct {
	DetectedText *string  `json:"detected_text"` // 识别出的文本行内容
	Confidence   *int64   `json:"confidence"`    // 置信度 0 ~100
	Polygon      []*Coord `json:"polygon"`       // 文本行在图像中的四点坐标
}

type Coord struct {
	X *int64 `json:"x"`
	Y *int64 `json:"y"`
}

type OcrRequest struct {
	ImageBase64  *string `json:"image_base64"`
	ImageUrl     *string `json:"image_url"`
	LanguageType *string `json:"language_type"`
}

type CreatePictureRequest struct {
	GroupId     *string `json:"group_id" binding:"required"`
	EntityId    *string `json:"entity_id" binding:"required"`
	PicName     *string `json:"pic_name" binding:"required"`
	ImageUrl    *string `json:"image_url"`
	ImageBase64 *string `json:"image_base64"`
	Tags        *string `json:"tags"`
}

type SearchPictureRequest struct {
	GroupId        *string    `json:"group_id" binding:"required"`
	ImageUrl       *string    `json:"image_url"`
	ImageBase64    *string    `json:"image_base64"`
	Limit          *int64     `json:"limit"`
	Offset         *int64     `json:"offset"`
	MatchThreshold *int64     `json:"match_threshold"`
	Filter         *string    `json:"filter"`
	ImageRect      *ImageRect `json:"image_rect"`
	EnableDetect   *bool      `json:"enable_detect"` // 是否需要启用主体识别，默认为TRUE 。
	CategoryId     *int64     `json:"category_id"`   // 识别的商品类别ID，若不填则为全部类别。
}

type ImageRect struct {
	X      *int64 `json:"x"`
	Y      *int64 `json:"y"`
	Width  *int64 `json:"width"`
	Height *int64 `json:"height"`
}

type SearchPictureResponse struct {
	RequestId  *string     `json:"request_id"`
	Count      *int64      `json:"count"`
	ImageInfos []ImageInfo `json:"image_infos"`
	Object     Object      `json:"object"`
}

type Object struct {
	Box        *Box     `json:"box"`
	Colors     []*Color `json:"colors"`
	CategoryId *int64   `json:"category_id"`
	Attributes []*Attr  `json:"attributes"`
	AllBox     []*Box   `json:"all_box"`
}

type Box struct {
	Rect  *ImageRect `json:"rect"`
	Score *int64     `json:"score"`
}

type Color struct {
	Color      *string  `json:"color"`
	Percentage *float64 `json:"percentage"`
	Label      *string  `json:"label"`
}

type Attr struct {
	Type    *string `json:"type"`
	Details *string `json:"details"`
}

type UpdatePictureRequest struct {
	GroupId  *string `json:"group_id" binding:"required"`
	EntityId *string `json:"entity_id" binding:"required"`
	PicName  *string `json:"pic_name"`
	Tags     *string `json:"tags"` // 新的自定义标签，最多不超过10个，格式为JSON
}

type CommonPictureResponse struct {
	RequestId *string `json:"request_id"`
}

type QueryPictureRequest struct {
	Offset  *int64  `json:"offset"`
	Limit   *int64  `json:"limit"`
	GroupId *string `json:"group_id"`
}

type QueryPictureResponse struct {
	Groups    []Group `json:"groups"`
	RequestId *string `json:"request_id"`
}

type Group struct {
	CreateTime  *string `json:"create_time"` // "2023-02-02 18:04:08"
	UpdateTime  *string `json:"update_time"` // "2023-02-02 18:04:08"
	GroupId     *string `json:"group_id"`
	GroupName   *string `json:"group_name"`
	GroupType   *int64  `json:"group_type"`
	MaxCapacity *int64  `json:"max_capacity"`
	MaxQps      *int64  `json:"max_qps"`
	PicCount    *int64  `json:"pic_count"`
	Brief       *string `json:"brief"`
}

type CreatePictureResponse struct {
	Object    Object  `json:"object"`
	RequestId *string `json:"request_id"`
}

type CommonPictureRequest struct {
	GroupId  *string `json:"group_id" binding:"required"`
	EntityId *string `json:"entity_id" binding:"required"`
	PicName  *string `json:"pic_name"`
}

type GetPictureByNameResponse struct {
	GroupId    *string     `json:"group_id"`
	EntityId   *string     `json:"entity_id"`
	ImageInfos []ImageInfo `json:"image_infos"`
	RequestId  *string     `json:"request_id"`
}

type ImageInfo struct {
	CustomContent *string `json:"custom_content"`
	EntityId      *string `json:"entity_id"`
	PicName       *string `json:"pic_name"`
	Score         *int64  `json:"score"`
	Tags          *string `json:"tags"`
}
