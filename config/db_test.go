package config

import (
	"Manager/domain/entity"
	"log"
	"testing"
)

func TestDatabaseConnection(t *testing.T) {
	dbInstance := NewDatabase()
	defer dbInstance.Close() // 测试后关闭连接

	// 2. 获取DB并查询用户表数据
	db := dbInstance.GetDB()
	//var users []entity.Users
	//if err := db.Find(&users).Error; err != nil {
	//	t.Fatal("查询用户表失败:", err)
	//}
	var admins []entity.Admins
	if err := db.Find(&admins).Error; err != nil {
		t.Fatal("查询用户表失败:", err)
	}

	// 3. 输出结果验证
	//log.Printf("成功查询到 %d 条用户数据", len(users))
	//t.Log("##############################")
	//if len(users) == 0 {
	//	t.Log("用户表暂无数据，但连接正常")
	//	t.Log("##############################")
	//}
	log.Printf("成功查询到 %d 条用户数据", len(admins))
	if len(admins) == 0 {
		t.Log("用户表暂无数据，但连接正常")
	}
}
