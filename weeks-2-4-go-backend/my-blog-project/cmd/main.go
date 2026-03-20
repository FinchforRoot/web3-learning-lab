package main

import (
	"fmt"
	"my-blog-project/config"
	"my-blog-project/database"
	"my-blog-project/router"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	// 1. 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		logrus.Errorf("加载配置失败: %v", err)
	}

	// 2. 初始化日志
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.Info("Starting blog application...")

	// 3. 初始化数据库
	database.InitDatabase()

	// 4. 设置 Gin 运行模式
	gin.SetMode(cfg.Server.Mode)

	// 5. 设置路由
	r := router.SetupRoutes()

	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	if err := r.Run(addr); err != nil {
		logrus.Error("Failed to start server:", err)
	}
}
