package model

// DescribeInstancePriceInput 定义查询实例价格的输入参数
type DescribeInstancePriceInput struct {
	Zone         *string `json:"zone" binding:"required"`         // 可用区
	ImageId      *string `json:"imageId" binding:"required"`      // 镜像ID
	InstanceType *string `json:"instanceType" binding:"required"` // 实例规格
	Period       *int64  `json:"period" binding:"required"`       // 购买时长，单位：月
	SystemDisk   *Disk   `json:"systemDisk" binding:"required"`   // 系统盘类型
	DataDisks    []Disk  `json:"dataDisks"`                       // 数据盘类型
}

// DescribeInstancePriceResponse 定义查询实例价格的返回结果
type DescribeInstancePriceResponse struct {
	OriginalPrice *float64 `json:"originalPrice"` // 原价，单位：元
	DiscountPrice *float64 `json:"discountPrice"` // 折扣价，单位：元
	Discount      *float64 `json:"discount"`      // 折扣，单位：%
	Currency      *string  `json:"currency"`      // 货币类型：CNY
}
