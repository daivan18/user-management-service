package handler

import (
	"log"
	"net/http"
	"net/url"

	"github.com/daivan18/user-management-service/utils"
	"github.com/gin-gonic/gin"
)

// 顯示註冊頁面（使用 Gin 的 c.HTML）
func ShowRegisterPage(c *gin.Context) {
	errorMsg := c.Query("error")
	successMsg := c.Query("success")

	c.HTML(http.StatusOK, "register.html", gin.H{
		"Error":   errorMsg,
		"Success": successMsg,
	})
}

// 處理註冊請求（使用 Gin 的表單綁定和重定向）
func RegisterHandler(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	cellPhone := c.PostForm("cell_phone")

	if username == "" || password == "" || cellPhone == "" {
		c.Redirect(http.StatusSeeOther, "/register?error=欄位不能為空")
		return
	}

	db := utils.GetDB()

	// 確認手機號碼是否已存在
	if utils.IsCellPhoneExists(db, cellPhone) {
		c.Redirect(http.StatusSeeOther, "/register?error="+url.QueryEscape("手機號碼已被註冊"))
		return
	}

	// 密碼加密（使用 utils.HashPassword）
	hashedPw, err := utils.HashPassword(password)
	if err != nil {
		log.Println("Hash failed:", err)
		c.Redirect(http.StatusSeeOther, "/register?error="+url.QueryEscape("密碼加密失敗"))
		return
	}

	// 手機號碼加密（使用 utils.Encrypt）
	encryptedPhone, err := utils.Encrypt(cellPhone)
	if err != nil {
		log.Println("Encrypt phone failed:", err)
		c.Redirect(http.StatusSeeOther, "/register?error="+url.QueryEscape("手機號碼加密失敗"))
		return
	}

	// 建立新使用者
	err = db.Exec("INSERT INTO users (username, password_hash, cell_phone) VALUES (?, ?, ?)",
		username, hashedPw, encryptedPhone).Error
	if err != nil {
		log.Println("新增使用者錯誤:", err)
		c.Redirect(http.StatusSeeOther, "/register?error="+url.QueryEscape("使用者建立失敗"))
		return
	}

	c.Redirect(http.StatusSeeOther, "/login?success="+url.QueryEscape("註冊成功，請登入"))
}
