package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// 设置路由
func setupRoutes() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// 简单的日志中间件
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// 健康检查
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	// API路由
	api := router.Group("/api")
	{

		// 钱包相关
		api.POST("/wallet/deposit", deposit)
		api.POST("/wallet/withdraw", withdraw)
		api.POST("/wallet/transfer", transfer)
		api.POST("/wallet/balance", getBalance)
		api.POST("/wallet/history", getHistory)

	}

	return router
}
