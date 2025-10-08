package core

import (
	"fmt"
	"go-ai-assistant/config"
	"go-ai-assistant/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() {
	var err error
	cfg := config.AppConfig.Database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{ // 3. gorm.Open 的参数也换成 mysql.Open
		Logger: logger.Default.LogMode(logger.Info), // 打印所有 SQL 语句
	})
	if err != nil {
		panic(fmt.Errorf("failed to connect database: %w", err))
	}
	fmt.Println("MySQL Database connection successful.")

	// 执行自动迁移
	autoMigrate()
}

// autoMigrate 自动迁移数据库表
func autoMigrate() {
	err := DB.AutoMigrate(
		&models.User{},
		&models.KnowledgeBase{},
		&models.Document{},
		&models.ChatSession{},
		&models.ChatMessage{},
	)
	if err != nil {
		panic(fmt.Errorf("failed to auto migrate database: %w", err))
	}
	fmt.Println("Database migration successful.")
}
