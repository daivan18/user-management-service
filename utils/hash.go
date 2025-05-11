package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword 將明碼轉換為 bcrypt 雜湊
func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}
