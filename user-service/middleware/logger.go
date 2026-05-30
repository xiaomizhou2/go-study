package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger 请求日志中间件
// 对比 Java: @Component public class LoggingFilter implements Filter { ... }
//
// 比喻：像商场的监控摄像头，记录每个访客什么时候来、待了多久、做了什么
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		// c.Next() 是关键！它把控制权交给下一个中间件或最终的处理函数
		// 对比 Java: chain.doFilter(request, response) 继续传递请求
		c.Next()

		// c.Next() 返回后，请求已经被处理完了，可以拿到响应状态
		latency := time.Since(start)
		status := c.Writer.Status()

		log.Printf("[%s] %s %d %v",
			method,
			path,
			status,
			latency,
		)
	}
}
