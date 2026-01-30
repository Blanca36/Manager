package main

import (
	"Manager/application"
	"Manager/config"
	"Manager/domain/service"
	"Manager/handle"
	"Manager/infrastructure"
	"Manager/router"
)

func main() {
	// 初始化数据库连接
	//（*sql.DB/gorm.DB） → PsgUserRepository（仓储层，直接操作数据库） → UserDomainImpl（业务层，调用仓储层处理业务逻辑）
	database := config.NewDatabase() //请求一个数据库实例

	// ========== 普通用户==========
	ur := infrastructure.NewPsgUserRepository(database.GetDB())
	ud := service.NewUserDomainImpl(ur)
	up := application.NewServiceImpl(ud)
	uh := handle.NewUserHandler(up)

	// ========== 管理员==========
	ar := infrastructure.NewPsgUserRepository(database.GetDB())
	ad := service.NewUserDomainImpl(ar)
	ap := application.NewServiceImpl(ad)
	ah := handle.NewUserHandler(ap)

	router := router.SetRouter(uh, ah)

	if err := router.Run(":8888"); err != nil {
		panic("服务启动失败：" + err.Error())
	}
}
