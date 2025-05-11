package main

import (
	"net/http"

	"github.com/daivan18/user-management-service/handler"
	"github.com/daivan18/user-management-service/middleware"
	"github.com/daivan18/user-management-service/utils"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {
	// 載入 .env 並初始化 utils
	utils.Init()

	// 初始化資料庫
	utils.InitDatabase()

	// 創建 Gin 引擎
	r := gin.Default()

	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	// 設定 static 資源路由
	r.Static("/static", "./static")

	// 設置 HTML 模板
	r.LoadHTMLGlob("templates/*")

	// 將首頁導向 login 頁面
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/login")
	})

	// 登入
	r.GET("/login", handler.ShowLoginPage)
	r.POST("/login", handler.VerifyLogin)

	// 註冊（新增的路由）
	r.GET("/register", handler.ShowRegisterPage)
	r.POST("/register", handler.RegisterHandler)

	// 登出
	r.GET("/logout", handler.Logout)

	// 登入後保護的管理路由
	auth := r.Group("/")
	auth.Use(middleware.RequireLogin())
	{
		auth.GET("/users", handler.ShowUserList)
		auth.POST("/users", handler.CreateUser)
		auth.GET("/users/:id/edit", handler.EditUserPage)
		auth.POST("/users/:id/update", handler.UpdateUser)
		auth.POST("/users/:id/delete", handler.DeleteUser)
	}

	// 啟動伺服器
	r.Run(":8080")
}
