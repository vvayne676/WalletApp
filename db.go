package main

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// 初始化数据库
func initDB() {
	config := loadConfig()

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=Asia/Shanghai",
		config.DBHost, config.DBPort, config.DBUser, config.DBPassword, config.DBName, config.DBSSLMode,
	)

	// 重试连接数据库
	maxRetries := 30
	retryInterval := 2 * time.Second

	log.Println("Connecting to database...")

	var err error
	for i := 0; i < maxRetries; i++ {
		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})

		if err == nil {
			// 测试连接
			sqlDB, err := DB.DB()
			if err == nil {
				if err := sqlDB.Ping(); err == nil {
					// 连接成功，设置连接池
					sqlDB.SetMaxIdleConns(10)
					sqlDB.SetMaxOpenConns(100)
					sqlDB.SetConnMaxLifetime(time.Hour)

					log.Println("Database connected successfully")
					return
				}
			}
		}

		log.Printf("Database connection attempt %d/%d failed: %v", i+1, maxRetries, err)
		log.Printf("Retrying in %v...", retryInterval)
		time.Sleep(retryInterval)
	}

	log.Fatal("Failed to connect to database after", maxRetries, "attempts:", err)
}
