package main

import (
	"log"
	"my-blog-backend/config"
	"my-blog-backend/database"
	"my-blog-backend/routes" // 1. 引入 routes 包

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("配置加载失败: %v", err)
	}

	// 2. 初始化数据库
	database.InitDB(cfg)

	// 3. 设置 Gin 模式 (Debug/Release)
	gin.SetMode(cfg.Server.Mode)

	// 4. 初始化路由
	// 这里调用了 routes/routes.go 中的 SetupRouter 函数
	// 它会返回一个已经注册好 /api/v1/posts 等所有路由的 Gin 引擎
	r := routes.SetupRouter()

	// 5. 启动服务器
	log.Printf("服务器启动在端口 %s", cfg.Server.Port)
	if err := r.Run(cfg.Server.Port); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
