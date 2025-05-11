package handler

import (
	"net/http"

	models "github.com/daivan18/user-management-service/model"
	"github.com/daivan18/user-management-service/utils"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// 顯示登入頁面
func ShowLoginPage(c *gin.Context) {
	c.HTML(200, "login.html", nil)
}

// 驗證登入
func VerifyLogin(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	var user models.User
	if err := utils.GetDB().Where("username = ?", username).First(&user).Error; err != nil {
		c.HTML(401, "login.html", gin.H{"Error": "帳號或密碼錯誤"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password_Hash), []byte(password)); err != nil {
		c.HTML(401, "login.html", gin.H{"Error": "帳號或密碼錯誤"})
		return
	}

	// 儲存登入狀態到 session
	session := sessions.Default(c)
	session.Set("username", username)
	session.Save()

	// 重定向到 /users 頁面
	c.Redirect(302, "/users")
}

// 登出
func Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete("username")
	session.Save()
	c.Redirect(http.StatusFound, "/login")
}
