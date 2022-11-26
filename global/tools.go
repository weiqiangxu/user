package global

import (
	gonanoid "github.com/matoous/go-nanoid"
)

// GenerateUniqueId 生成唯一码
func GenerateUniqueId(size int) string {
	uuid, err := gonanoid.Generate("123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWSYZ", size)
	if err != nil {
		return ""
	}
	return uuid
}
