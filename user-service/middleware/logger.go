package middleware

import (
	"time"

	"github.com/example/user-service/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Logger 请求日志中间件（zap 结构化版本）
// 对比 Java: @Component public class LoggingFilter implements Filter { ... }
//
// 比喻：像商场的智能监控系统，不仅记录"谁来了"，还记录"待了多久"、"结果如何"
// 并且所有日志都是结构化的 JSON，方便 ELK/Grafana 采集分析
//
// 输出示例（json 格式）：
//   {"level":"info","method":"GET","path":"/api/v1/users","status":200,"latency":"0.5ms"}
//
// 对比之前的标准库 log.Printf：
//   [GET] /api/v1/users 200 500µs  ← 机器很难解析
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		// zap 的结构化日志：每个字段都是独立的 key-value
		// 对比 Java SLF4J: log.info("method={}, path={}, status={}, latency={}", ...);
		// zap 的方式更适合机器解析，不需要字符串拼接
		logger.Log.Info("request",
			zap.String("method", method),
			zap.String("path", path),
			zap.Int("status", status),
			zap.Duration("latency", latency),
			zap.String("client_ip", c.ClientIP()),
		)
	}
}
