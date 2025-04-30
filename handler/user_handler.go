package handler

import (
	"log"
	"net/http"
	"net/url"
	"time"

	models "github.com/daivan18/user-management-service/model"
	"github.com/daivan18/user-management-service/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func ShowUserList(c *gin.Context) {
	var users []models.User

	if err := utils.GetDB().Select("id", "username").Find(&users).Error; err != nil {
		log.Println("查詢錯誤:", err)
		c.String(http.StatusInternalServerError, "DB error")
		return
	}

	var userList []map[string]interface{}
	for _, u := range users {
		userList = append(userList, map[string]interface{}{
			"id":       u.ID,
			"username": u.Username,
		})
	}

	// 成功訊息
	successMessage := ""
	switch c.Query("success") {
	case "created":
		successMessage = "User created successfully!"
	case "updated":
		successMessage = "User updated successfully!"
	case "deleted":
		successMessage = "User deleted successfully!"
	}

	errorMessage := c.Query("error")

	c.HTML(http.StatusOK, "user.html", gin.H{
		"users":   userList,
		"success": successMessage,
		"error":   errorMessage,
	})
}

func CreateUser(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	if isUsernameExists(username) {
		c.Redirect(http.StatusSeeOther, "/users?error="+url.QueryEscape("帳號已存在"))
		return
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		log.Println("Hash failed:", err)
		c.String(http.StatusInternalServerError, "Error hashing password")
		return
	}

	user := models.User{
		Username:      username,
		Password_Hash: hashedPassword,
		Create_Time:   time.Now(),
	}

	if err := utils.GetDB().Create(&user).Error; err != nil {
		log.Println("Insert failed:", err)
		c.String(http.StatusInternalServerError, "Error creating user")
		return
	}

	c.Redirect(http.StatusSeeOther, "/users?success=created")
}

func isUsernameExists(username string) bool {
	var count int64
	utils.GetDB().Model(&models.User{}).Where("username = ?", username).Count(&count)
	return count > 0
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

// ShowAdminLoginPage 顯示管理員登入頁面
func ShowAdminLoginPage(c *gin.Context) {
	c.HTML(200, "admin_login.html", nil)
}

// VerifyAdminLogin 驗證管理員帳號與密碼
func VerifyAdminLogin(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	// 獲取資料庫連線
	db := utils.GetDB()

	// 查詢資料庫中是否有此用戶
	var admin models.User
	if err := db.Where("username = ?", username).First(&admin).Error; err != nil {
		// 賬號不存在
		c.HTML(401, "admin_login.html", gin.H{
			"Error": "帳號或密碼錯誤",
		})
		return
	}

	// 驗證密碼（解密密碼）
	err := bcrypt.CompareHashAndPassword([]byte(admin.Password_Hash), []byte(password))
	if err != nil {
		// 密碼不匹配
		c.HTML(401, "admin_login.html", gin.H{
			"Error": "帳號或密碼錯誤",
		})
		return
	}

	// 密碼驗證通過，跳轉到管理員頁面
	c.Redirect(302, "/admin/users")
}

// ShowAllUsersWithPassword 顯示所有帳號與密碼（測試用，不建議正式使用）
func ShowAllUsersWithPassword(c *gin.Context) {
	// 定義一個 User 陣列來儲存所有用戶資料
	var users []models.User

	// 獲取資料庫連線
	db := utils.GetDB()

	// 使用 GORM 查詢所有用戶資料
	if err := db.Find(&users).Error; err != nil {
		c.JSON(500, gin.H{"error": "無法獲取用戶資料"})
		return
	}

	// 返回結果給前端，顯示所有用戶資料
	c.HTML(200, "admin.html", gin.H{
		"Users": users,
	})
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
