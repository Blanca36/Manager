package router

import (
	"Manager/handler"
	"Manager/interface/middleware"

	"github.com/gin-gonic/gin"
)

// 路由
func SetUserRouter(userHandler *handler.UserHandler) *gin.Engine {
	r := gin.Default()
	//设置全局中间件
	r.Use(middleware.Cors(), middleware.RequestID(), middleware.Counter())
	//路由组,PS:登录是「认证动作」
	authGroup := r.Group("auth")
	{
		authGroup.POST("users/login", userHandler.UserLogin) // 登录接口
	}
	userGroup := r.Group("users/me")
	{
		userGroup.PUT("password", userHandler.UpdatePsd)
	}
	return r
}

func SetAdminRouter(adminHandler *handler.AdminHandler) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.Cors(), middleware.RequestID(), middleware.Counter())
	//管理员登录
	authGroup := r.Group("auth")
	{
		authGroup.POST("admins/login", adminHandler.AdminLogin)
	}

	//管理员验证
	//adminsGroup := r.Group("admins").Use(middleware.AdminAuth())
	//{
	//
	//}
	//管理员管理普通用户

	adminUserGroup := r.Group("admin/users").Use(middleware.AdminAuth())
	{
		adminUserGroup.DELETE("/:id", adminHandler.DeleteUser)
		adminUserGroup.GET("", adminHandler.GetUsersList)
		//users.POST("", xxxx)
		adminUserGroup.PUT("/:id/password", adminHandler.UpdateUserPassword) ///password/update，名词在前动词在后
		//admins.PUT("/password/update", adminHandler.UpdateUserPassword)
		adminUserGroup.PUT("/:id/username", adminHandler.UpdateUsername)

	}
	return r
}
