package model

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
)

/*
aws sqs 创建接口，默认关闭加密
*/

var default_policy string = `{
  "Version": "2012-10-17",
  "Id": "__default_policy_ID",
  "Statement": [
    {
      "Sid": "__owner_statement",
      "Effect": "Allow",
      "Principal": {
        "AWS": "006694404643"
      },
      "Action": [
        "SQS:*"
      ],
      "Resource": "arn:aws:sqs:us-east-1:006694404643:"
    }
  ]
}`

type SqsConfig struct {
	// 可见性超时 0-12小时 from 0 to 43,200 (12 hours). Default: 30.
	VisibilityTimeout int `json:"visibility_timeout"`
	// 消息保留时间 应介于 60 秒至 1209600 秒之间。
	MessageRetention int `json:"message_retention"`
	// 最大消息大小 KB 应介于  1024-262144
	MaximumMessageSize int `json:"max_message_size"`
	// 接收消息等待时间 0-20秒
	ReceiveWaitTime int `json:"receive_wait_time"`
	// 交付延迟 0-15分钟
	DelaySeconds int `json:"delay_seconds"`
}

type CreateSqsRequest struct {
	QueueName  string            `json:"queue_name"`
	Type       string            `json:"type"`   // normal | fifo
	Policy     string            `json:"policy"` // 策略
	Config     SqsConfig         `json:"config"`
	Encryption bool              `json:"encryption"` // 是否开启加密
	Tags       map[string]string `json:"tags"`       // 标签
}

// to sqs.CreateQueueInput
func (c *CreateSqsRequest) ToCreateQueueInput() *sqs.CreateQueueInput {
	// 规范化队列名称（AWS要求只能包含字母数字、连字符和下划线）
	queueName := strings.ToLower(c.QueueName)
	queueName = strings.ReplaceAll(queueName, "_", "-")
	queueName = strings.ReplaceAll(queueName, " ", "-")

	r := &sqs.CreateQueueInput{
		QueueName: aws.String(queueName),
		Attributes: map[string]*string{
			"VisibilityTimeout":             aws.String(strconv.Itoa(c.Config.VisibilityTimeout)),
			"MessageRetentionPeriod":        aws.String(strconv.Itoa(c.Config.MessageRetention)),
			"MaximumMessageSize":            aws.String(strconv.Itoa(c.Config.MaximumMessageSize)),
			"ReceiveMessageWaitTimeSeconds": aws.String(strconv.Itoa(c.Config.ReceiveWaitTime)),
			"DelaySeconds":                  aws.String(strconv.Itoa(c.Config.DelaySeconds)),
			"SqsManagedSseEnabled":          aws.String(strconv.FormatBool(c.Encryption)),
		},
		Tags: map[string]*string{},
	}
	if c.Type == "fifo" {
		r.Attributes["FifoQueue"] = aws.String("true")
		r.QueueName = aws.String(strings.TrimSuffix(queueName, ".fifo") + ".fifo")
	}
	for k, v := range c.Tags {
		r.Tags[k] = aws.String(v)
	}
	r.Tags["Name"] = r.QueueName

	if c.Policy != "" {
		r.Attributes["Policy"] = aws.String(c.Policy)
	}

	fmt.Println(tea.Prettify(r))
	return r
}
