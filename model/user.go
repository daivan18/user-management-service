package models

import "time"

type User struct {
	ID            int       `json:"id"`
	Username      string    `json:"username"`
	Password_Hash string    `json:"password_hash"` // 不回傳密碼 hash 給前端
	Create_Time   time.Time `json:"create_time"`   // 時間格式
	Update_Time   time.Time `json:"update_time"`   // 時間格式
}
