package utils

import (
	"log"

	models "github.com/daivan18/user-management-service/model"
	"gorm.io/gorm"
)

// MaskPhone 將手機號碼中間三碼遮蔽，例如 0912***131
func MaskPhone(phone string) string {
	if len(phone) < 7 {
		return phone // 資料太短無法遮蔽
	}
	return phone[:4] + "***" + phone[len(phone)-3:]
}

// 判斷手機號碼是否已存在（解密比對）
func IsCellPhoneExists(db *gorm.DB, cellPhone string) bool {
	rows, err := db.
		Model(&models.User{}).
		Select("cell_phone").
		Where("cell_phone IS NOT NULL AND cell_phone != ''").
		Rows()
	if err != nil {
		log.Println("DB query failed:", err)
		return false
	}
	defer rows.Close()

	for rows.Next() {
		var encryptedPhone string
		if err := rows.Scan(&encryptedPhone); err != nil {
			log.Println("Scan failed:", err)
			continue
		}

		decryptedPhone, err := Decrypt(encryptedPhone)
		if err != nil {
			log.Println("Decrypt failed:", err)
			continue
		}

		if decryptedPhone == cellPhone {
			return true
		}
	}

	return false
}
