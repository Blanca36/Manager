package main

import (
	"Manager/config"
	"Manager/infrastructure"
)

func main() {

	database := config.NewDatabase() //请求一个数据库实例
	//依赖注入
	userRepo := infrastructure.NewPsgUserRepository(database.GetDB())
}
