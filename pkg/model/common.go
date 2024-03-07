package model

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/alibabacloud-go/tea/tea"
)

const (
	slot = "12345678901234567890123456789012"
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
	_data, _ := encrypt(data, slot)
	return _data
}

func DecodeAwsNextMaker(data string) (*string, *string) {
	decodeData, _ := decrypt(data, slot)
	dataSet := strings.SplitN(string(decodeData), ",", 2)
	return tea.String(dataSet[0]), tea.String(dataSet[1])
}

func ToTencentNextMaker(data string) *string {
	_data, _ := encrypt(data, slot)
	return tea.String(_data)
}

func DecodeTencentNextMaker(data string) string {
	_data, _ := decrypt(data, slot)
	return _data
}

func encrypt(text, key string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(text))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(text))
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func decrypt(encryptedText, key string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	ciphertext, err := base64.StdEncoding.DecodeString(encryptedText)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < aes.BlockSize {
		return "", fmt.Errorf("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(ciphertext, ciphertext)
	return string(ciphertext), nil
}
