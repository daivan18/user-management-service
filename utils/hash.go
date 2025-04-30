package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// bcrypt雜湊密碼
func HashPassword(password string) (string, error) {
	// bcrypt雜湊密碼較一般密碼Hash的方式安全性高：慢速且安全防止暴力破解、自動內建 salt、可用 cost factor調整難度成本
	// 使用bcrypt雜湊密碼，成本因子設為 14
	// bcrypt是一種密碼雜湊函數，使用成本因子來調整計算時間
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
