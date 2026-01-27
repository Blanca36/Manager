package config

import (
	"Manager/domain/entity"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type Database struct {
	db *gorm.DB
}

// 创建数据库连接实例
func NewDatabase() *Database {
	dsn := "root:Root.123@(127.0.0.1:3306)/demo?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}

	// 配置连接池
	db.DB().SetMaxOpenConns(100)
	db.DB().SetMaxIdleConns(10)

	// 自动迁移
	db.AutoMigrate(&entity.User{})

	return &Database{db: db}
}

// 新增：获取 GORM DB 实例的方法
func (d *Database) GetDB() *gorm.DB {
	return d.db
}

// Close 关闭数据库连接
func (d *Database) Close() {
	if d.db != nil {
		d.db.Close()
		log.Println("数据库连接已关闭")
	}
}
