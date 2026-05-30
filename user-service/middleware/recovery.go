package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/example/user-service/logger"
	"github.com/example/user-service/util"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Recovery panic 恢复中间件
// 对比 Java: @ControllerAdvice + @ExceptionHandler(Throwable.class)
//
// 比喻：像电路的保险丝，万一短路（panic）了，自动跳闸保护整个系统
// 不会因为一个请求的 panic 导致整个服务崩溃
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 记录 panic 详情 + 堆栈
				// 对比 Java: log.error("Unhandled exception", throwable);
				logger.Log.Error("panic recovered",
					zap.Any("error", err),
					zap.String("stack", string(debug.Stack())),
					zap.String("path", c.Request.URL.Path),
					zap.String("method", c.Request.Method),
				)

				util.Error(c, http.StatusInternalServerError, 500, "服务器内部错误")
				c.Abort()
			}
		}()

		c.Next()
	}
}
