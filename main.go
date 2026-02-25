package main

import (
	"Manager/config"
	"Manager/domain/service"
	"Manager/handler"
	"Manager/infrastructure"
	"Manager/interface/router"
	"Manager/tool"

	"go.uber.org/zap"
)

func main() {
	tool.InitLogger()
	tool.Info("程序启动，日志初始化完成")

	// 初始化数据库连接
	//（*sql.DB/gorm.DB） → PsgUserRepo（仓储层，直接操作数据库） → UserDomainImpl（业务层，调用仓储层处理业务逻辑）
	database := config.NewDatabase() //请求一个数据库实例

	// ========== 普通用户==========
	ur := infrastructure.NewPsgUserRepo(database.GetDB())
	ud := service.NewUserSrv(ur)
	up := service.NewUserAppSvcImpl(ud)
	uh := handler.NewUserHandler(up)

	// ========== 管理员==========
	ar := infrastructure.NewPsgAdminRepo(database.GetDB())
	ad := service.NewAdminSrv(ar)
	ap := service.NewAdminAppSrvImpl(ad, ud) //注入用户相关依赖
	ah := handler.NewAdminHandler(ap)

	router := router.SetRouter(uh, ah)
	tool.Info("路由初始化完成，启动HTTP服务，端口：8889")
	if err := router.Run(":8889"); err != nil {
		tool.Error("服务启动失败", zap.Error(err))
		panic("服务启动失败：" + err.Error())
	}
}
