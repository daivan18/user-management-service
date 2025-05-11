package handler

import (
	"log"
	"net/http"
	"net/url"
	"time"

	models "github.com/daivan18/user-management-service/model"
	"github.com/daivan18/user-management-service/utils"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func ShowUserList(c *gin.Context) {
	// 從 session 中提取 username
	session := sessions.Default(c)
	username := session.Get("username")
	if username == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	var users []models.User
	db := utils.GetDB()

	if username == "admin" {
		if err := db.Select("id", "username", "cell_phone").Find(&users).Error; err != nil {
			c.String(http.StatusInternalServerError, "DB error")
			return
		}
	} else {
		var user models.User
		if err := db.Select("id", "username", "cell_phone").Where("username = ?", username).First(&user).Error; err != nil {
			c.String(http.StatusInternalServerError, "DB error")
			return
		}
		users = append(users, user)
	}

	var userList []map[string]interface{}
	for _, u := range users {

		var err error

		// 解密手機號碼
		var decryptedPhone string

		if u.CellPhone != "" {
			// 如果手機號碼不為空，才進行解密
			decryptedPhone, err = utils.Decrypt(u.CellPhone)
			if err != nil {
				// 如果解密失敗，設為 "解密失敗"
				decryptedPhone = "解密失敗"
			}
		} else {
			// 如果手機號碼為空，直接設為空
			decryptedPhone = ""
		}

		// 對解密後的手機號碼進行遮蔽處理
		maskedPhone := decryptedPhone
		//maskedPhone := utils.MaskPhone(decryptedPhone)

		userList = append(userList, map[string]interface{}{
			"id":         u.ID,
			"username":   u.Username,
			"cell_phone": maskedPhone,
		})
	}

	c.HTML(http.StatusOK, "user.html", gin.H{
		"users":   userList,
		"success": c.Query("success"),
		"error":   c.Query("error"),
	})
}

func CreateUser(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	cellPhone := c.PostForm("cell_phone")

	// 確認手機號碼是否已存在
	if utils.IsCellPhoneExists(utils.GetDB(), cellPhone) {
		c.Redirect(http.StatusSeeOther, "/users?error="+url.QueryEscape("手機號碼已被註冊"))
		return
	}

	// 密碼加密
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		log.Println("Hash failed:", err)
		c.String(http.StatusInternalServerError, "Error hashing password")
		return
	}

	// 手機號碼加密
	encryptedPhone, err := utils.Encrypt(cellPhone)
	if err != nil {
		log.Println("Encrypt phone failed:", err)
		c.String(http.StatusInternalServerError, "Error encrypting phone")
		return
	}

	// 創建新用戶
	user := models.User{
		Username:      username,
		Password_Hash: hashedPassword,
		CellPhone:     encryptedPhone,
		Create_Time:   time.Now(),
	}

	if err := utils.GetDB().Create(&user).Error; err != nil {
		log.Println("Insert failed:", err)
		c.String(http.StatusInternalServerError, "Error creating user")
		return
	}

	c.Redirect(http.StatusSeeOther, "/users?success=created")
}

func EditUserPage(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	if err := utils.GetDB().First(&user, id).Error; err != nil {
		log.Println("Edit user not found:", err)
		c.String(http.StatusNotFound, "User not found")
		return
	}

	c.HTML(http.StatusOK, "edit.html", gin.H{
		"id":       user.ID,
		"username": user.Username,
	})
}

func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	username := c.PostForm("username")
	password := c.PostForm("password")
	//cellPhone := c.PostForm("cell_phone")

	// 確認手機號碼是否已存在
	/*
		if utils.IsCellPhoneExists(utils.GetDB(), cellPhone) {
			c.Redirect(http.StatusSeeOther, "/users?error="+url.QueryEscape("手機號碼已被註冊"))
			return
		}
	*/

	// 雜湊密碼
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		log.Println("Hash password failed:", err)
		c.String(http.StatusInternalServerError, "Hash failed")
		return
	}

	result := utils.GetDB().Exec(
		"UPDATE users SET username=$1, password_hash=$2, update_time=$3 WHERE id=$4",
		username, hashedPassword, time.Now(), id,
	)

	if result.Error != nil {
		log.Println("Update failed:", result.Error)
		c.String(http.StatusInternalServerError, "Update failed")
		return
	}

	c.Redirect(http.StatusSeeOther, "/users?success=updated")
}

func DeleteUser(c *gin.Context) {
	id := c.Param("id")

	db := utils.GetDB()
	result := db.Exec("DELETE FROM users WHERE id = ?", id)
	if result.Error != nil {
		log.Println("刪除用戶錯誤:", result.Error)
		c.JSON(500, gin.H{"error": "刪除用戶時發生錯誤"})
		return
	}

	c.Redirect(http.StatusSeeOther, "/users?success=deleted")
}

// 重設密碼
func ResetUserPassword(c *gin.Context) {
	id := c.PostForm("id")
	newPassword := c.PostForm("new_password")

	// 雜湊新密碼
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		log.Println("Hash password failed:", err)
		c.String(http.StatusInternalServerError, "密碼雜湊失敗")
		return
	}

	// 更新資料庫
	db := utils.GetDB()
	if err := db.Model(&models.User{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"password_hash": hashedPassword,
			"update_time":   time.Now(),
		}).Error; err != nil {
		log.Println("更新密碼失敗:", err)
		c.String(http.StatusInternalServerError, "更新密碼失敗")
		return
	}
	// 更新一個欄位(跟上面多個欄位的比較，比較完刪除)
	// db := utils.GetDB()
	// if err := db.Model(&models.User{}).
	// 	Where("id = ?", id).
	// 	Update("password_hash", hashedPassword).Error; err != nil {
	// 	log.Println("更新密碼失敗:", err)
	// 	c.String(http.StatusInternalServerError, "更新密碼失敗")
	// 	return
	// }

	c.Redirect(http.StatusSeeOther, "/admin")
}
