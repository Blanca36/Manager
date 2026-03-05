package main

import (
	"Manager/config"
	"Manager/domain/service"
	"Manager/handler"
	"Manager/infrastructure"
	"Manager/interface/router"
	"Manager/tool"
	"log"

	"go.uber.org/zap"
)

func main() {
	tool.InitLogger()
	tool.Info("程序启动，日志初始化完成")

	// 初始化数据库连接
	//（*sql.DB/gorm.DB） → PsgUserRepo（仓储层，直接操作数据库） → UserDomainImpl（业务层，调用仓储层处理业务逻辑）
	database := config.NewDatabase() //请求一个数据库实例

	// ========== 普通用户==========
	userRepo := infrastructure.NewPsgUserRepo(database.GetDB())
	userService := service.NewUserService(userRepo)
	userHandle := handler.NewUserHandler(userService)

	// ========== 管理员==========
	adminRepo := infrastructure.NewPsgAdminRepo(database.GetDB(), userRepo)
	adminService := service.NewAdminServiceImpl(adminRepo) //注入用户相关依赖
	adminHandler := handler.NewAdminHandler(adminService)

	userRouter := router.SetUserRouter(userHandle)

	go func() {
		tool.Info("用户接口服务初始化完成，启动HTTP服务，端口：8886")
		if err := userRouter.Run(":8886"); err != nil {
			tool.Error("用户服务启动失败", zap.Error(err))
			log.Panicf("用户服务启动失败：%s", err.Error())
		}
	}()

	adminRouter := router.SetAdminRouter(adminHandler)
	tool.Info("管理员接口服务初始化完成，启动HTTP服务，端口：8888")
	if err := adminRouter.Run(":8888"); err != nil {
		tool.Error("管理员服务启动失败", zap.Error(err))
		panic("管理员服务启动失败：" + err.Error())
	}
}
