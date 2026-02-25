package config

import (
	"Manager/domain/entity"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	db *gorm.DB
}

// 创建数据库连接实例
func NewDatabase() *Database {
	dsn := "host=127.0.0.1 port=5433 user=postgres password=Root123 dbname=manger_user sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		// 可选：配置日志级别，方便调试（生产环境可设为 logger.Error）
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}

	// 配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("获取 sql.DB 实例失败:", err)
	}
	sqlDB.SetMaxOpenConns(100)              // 最大打开连接数
	sqlDB.SetMaxIdleConns(10)               // 最大空闲连接数
	sqlDB.SetConnMaxLifetime(1 * time.Hour) // 新增：连接最大生命周期（避免无效连接）

	// 自动迁移
	db.AutoMigrate(&entity.Admins{})

	log.Println("Postgres 连接成功，数据库：manger_user")
	return &Database{db: db}
}

// 获取 GORM DB 实例的方法
func (d *Database) GetDB() *gorm.DB {
	return d.db
}

// Close 关闭数据库连接
func (d *Database) Close() {
	if d.db != nil {
		sqlDB, err := d.db.DB()
		if err != nil {
			log.Printf("获取 sql.DB 实例失败，关闭连接失败：%v\n", err)
			return
		}
		err = sqlDB.Close()
		if err != nil {
			log.Printf("关闭数据库连接失败：%v\n", err)
			return
		}
		log.Println("数据库连接已关闭")
	}
}
