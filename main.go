package main

import (
	"github.com/daivan18/user-management-service/handler" // 引用 handler 包
	"github.com/daivan18/user-management-service/utils"
	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化資料庫
	utils.InitDatabase()

	// 創建 Gin 引擎
	r := gin.Default()

	// 設定 static 資源路由 告訴 Gin 去哪裡找 static 資源路由
	r.Static("/static", "./static")

	// 設置 HTML 模板 (告訴 Gin 去哪裡找 HTML 模板)
	r.LoadHTMLGlob("templates/*")

	// 設置路由，引用 handler 包中的函式
	r.GET("/users", handler.ShowUserList)           // 顯示用戶清單
	r.POST("/users", handler.CreateUser)            // 創建新用戶
	r.GET("/users/:id/edit", handler.EditUserPage)  // 顯示編輯畫面
	r.POST("/users/:id/update", handler.UpdateUser) // 處理更新資料
	r.POST("/users/:id/delete", handler.DeleteUser) // 刪除用戶

	// 管理者功能
	// 新增顯示帳密頁面 (受密碼保護)
	r.GET("/admin", handler.ShowAdminLoginPage)
	r.POST("/admin", handler.VerifyAdminLogin)
	r.GET("/admin/users", handler.ShowAllUsersWithPassword)
	// 修改密碼
	r.POST("/admin/reset-password", handler.ResetUserPassword)

	// 啟動伺服器
	r.Run(":8080")
}
