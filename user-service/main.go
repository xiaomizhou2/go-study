package main

import (
	"fmt"
	"log"

	"github.com/example/user-service/config"
	"github.com/example/user-service/handler"
	"github.com/example/user-service/middleware"
	"github.com/example/user-service/repository"
	"github.com/example/user-service/service"
	"github.com/gin-gonic/gin"
)

func main() {
	// ============ 1. 加载配置 ============
	// 对比 Java: Spring 自动读取 application.yml
	// Go: 手动调用，但只有一行
	cfg, err := config.Load("./config.yaml")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}
	log.Printf("配置加载成功: 端口=%s, 数据库=%s", cfg.Server.Port, cfg.Database.DSN)

	// ============ 2. 初始化数据库 ============
	// 对比 Java: @Bean DataSource 根据 application.yml 自动配置
	// Go: 把 Config 传进去，让 database.go 自己读取参数
	db, err := config.InitDB(cfg.Database)
	if err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}
	defer db.Close()

	// ============ 3. 组装依赖链 ============
	repo := repository.NewUserRepository(db)
	svc := service.NewUserService(repo)
	userHandler := handler.NewUserHandler(svc)

	// ============ 4. 创建路由引擎 ============
	r := gin.New()

	// ============ 5. 全局中间件 ============
	r.Use(middleware.CORS())
	r.Use(middleware.Logger())
	r.Use(gin.Recovery())

	// ============ 6. 注册路由 ============
	r.GET("/api/v1/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "user-service",
		})
	})

	users := r.Group("/api/v1/users")
	{
		users.GET("", userHandler.GetAll)
		users.GET("/:id", userHandler.GetByID)

		auth := users.Group("")
		auth.Use(middleware.Auth())
		{
			auth.POST("", userHandler.Create)
			auth.PUT("/:id", userHandler.Update)
			auth.DELETE("/:id", userHandler.Delete)
		}
	}

	// ============ 7. 启动服务 ============
	printRoutes(cfg)
	log.Printf("服务启动在 http://localhost%s", cfg.Server.Port)
	if err := r.Run(cfg.Server.Port); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}

func printRoutes(cfg *config.Config) {
	fmt.Println("🚀 user-service 已启动")
	fmt.Println("📍 可用的端点:")
	fmt.Println("   GET    /api/v1/users       - 获取所有用户")
	fmt.Println("   GET    /api/v1/users/:id   - 按 ID 获取用户")
	fmt.Println("   POST   /api/v1/users       - 创建新用户 [需认证]")
	fmt.Println("   PUT    /api/v1/users/:id   - 更新用户 [需认证]")
	fmt.Println("   DELETE /api/v1/users/:id   - 删除用户 [需认证]")
	fmt.Println("   GET    /api/v1/health      - 健康检查")
	fmt.Println()
	fmt.Println("🔑 认证方式: 请求头加 Authorization: Bearer <任意字符串>")
	fmt.Printf("💾 数据库: %s (%s)\n", cfg.Database.Driver, cfg.Database.DSN)
}
