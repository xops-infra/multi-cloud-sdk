package model

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/alibabacloud-go/tea/tea"
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

func ToAwsNextMaker(name, nextType *string) string {
	data := fmt.Sprintf("%s,%s", *name, *nextType)
	return encrypt(data, slot)
}

func DecodeAwsNextMaker(data string) (*string, *string) {
	decodeData := decrypt(data, slot)
	dataSet := strings.SplitN(string(decodeData), ",", 2)
	return tea.String(dataSet[0]), tea.String(dataSet[1])
}

func ToTencentNextMaker(data string) *string {
	return tea.String(encrypt(data, slot))
}

func DecodeTencentNextMaker(data string) string {
	return decrypt(data, slot)
}

func encrypt(text, key string) string {
	return base64.StdEncoding.EncodeToString([]byte(text + key))
}

func decrypt(encryptedText, key string) string {
	base64Data, _ := base64.StdEncoding.DecodeString(encryptedText)
	return strings.TrimRight(string(base64Data), key)
}
