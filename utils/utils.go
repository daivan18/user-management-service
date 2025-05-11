// 初始化函數，載入環境變數
package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv" // 載入 godotenv
)

var EncryptionKey string

func Init() {
	// 載入 .env 檔案
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	EncryptionKey = os.Getenv("ENCRYPTION_KEY")
	if EncryptionKey == "" {
		log.Fatal("ENCRYPTION_KEY 沒有設定")
	}
}
