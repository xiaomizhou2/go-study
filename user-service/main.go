package main

import (
	"fmt"
	"log"

	"github.com/example/user-service/config"
	"github.com/example/user-service/handler"
	"github.com/example/user-service/logger"
	"github.com/example/user-service/middleware"
	"github.com/example/user-service/repository"
	"github.com/example/user-service/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	cfg, err := config.Load("./config.yaml")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	if err := logger.Init(cfg.Log); err != nil {
		log.Fatalf("日志初始化失败: %v", err)
	}
	defer logger.Sync()

	logger.Log.Info("配置加载成功",
		zap.String("port", cfg.Server.Port),
		zap.String("db", cfg.Database.DSN),
		zap.String("log_level", cfg.Log.Level),
	)

	db, err := config.InitDB(cfg.Database)
	if err != nil {
		logger.Log.Fatal("数据库初始化失败", zap.Error(err))
	}
	defer db.Close()

	repo := repository.NewUserRepository(db)
	svc := service.NewUserService(repo)
	userHandler := handler.NewUserHandler(svc)

	r := setupRouter(userHandler)

	printRoutes(cfg)
	logger.Log.Info("服务启动", zap.String("addr", fmt.Sprintf("http://localhost%s", cfg.Server.Port)))
	if err := r.Run(cfg.Server.Port); err != nil {
		logger.Log.Fatal("服务启动失败", zap.Error(err))
	}
}

// setupRouter 注册所有路由和中间件
// 对比 Java: WebMvcConfigurer 或 SecurityFilterChain 配置类
// 提取为独立函数，测试时可以直接调用拿到 *gin.Engine
func setupRouter(h *handler.UserHandler) *gin.Engine {
	r := gin.New()

	r.Use(middleware.CORS())
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())

	r.GET("/api/v1/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "user-service",
		})
	})

	users := r.Group("/api/v1/users")
	{
		users.GET("", h.GetAll)
		users.GET("/:id", h.GetByID)

		auth := users.Group("")
		auth.Use(middleware.Auth())
		{
			auth.POST("", h.Create)
			auth.PUT("/:id", h.Update)
			auth.DELETE("/:id", h.Delete)
		}
	}

	return r
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
