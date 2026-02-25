package router

import (
	"Manager/handler"
	"Manager/interface/middleware"

	"github.com/gin-gonic/gin"
)

// 路由
func SetRouter(userHandler *handler.UserHandler, adminHandler *handler.AdminHandler) *gin.Engine {
	r := gin.Default()
	//设置全局中间件
	r.Use(middleware.Cors(), middleware.RequestID(), middleware.Counter())
	//路由组
	userMe := r.Group("user/me")
	{
		userMe.POST("login", userHandler.UserLogin) // 登录接口
		userMe.PUT("password/update", userHandler.UpdatePsd)
	}
	//管理员登录
	adminMe := r.Group("admin/me")
	{
		adminMe.POST("login", adminHandler.AdminLogin)
	}

	//管理员验证
	admins := r.Group("admins").Use(middleware.AdminAuth())
	{
		admins.PUT("/:id/password/update", adminHandler.AdminUpdateUPsd) ///password/update，名词在前动词在后
		admins.PUT("username/update", adminHandler.AdminUpdateUName)
		admins.PUT("user/delete", adminHandler.DeleteUser)
		admins.GET("user_list", adminHandler.GetUsersList)
	}

	return r
}
