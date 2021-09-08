package models

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

// DB 数据库链接单例
var DB *gorm.DB

// Init 初始化SQLLite数据库链接
func Init() {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second,  // 慢 SQL 阈值
			LogLevel:      logger.Error, // 日志级别
		},
	)

	db, err := gorm.Open(sqlite.Open("data.db"), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Panicln("无法连接到SQLite数据库.")
	}

	// Score表
	err = db.AutoMigrate(&Score{})
	if err != nil {
		return
	}

	// Usage表
	err = db.AutoMigrate(&Usage{})
	if err != nil {
		return
	}

	DB = db
}
