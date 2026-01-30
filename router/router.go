package router

import (
	"Manager/handle"

	"github.com/gin-gonic/gin"
)

// 路由
func SetRouter(userHandler *handle.UserHandler, adminHandler *handle.UserHandler) *gin.Engine {
	r := gin.Default()
	// 初始化数据库连接

	//路由组
	rGroup := r.Group("/user")
	{
		rGroup.POST("/login", userHandler.LocalLoginUserHandler) // 登录接口
	}

	adminGroup := r.Group("/admin")
	{
		adminGroup.POST("/login", adminHandler.LocalLoginUserHandler) // 登录接口
	}
	return r
}
