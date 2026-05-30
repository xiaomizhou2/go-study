package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/example/user-service/util"
)

// Auth 认证中间件
// 对比 Java: @Component public class AuthFilter implements Filter { ... }
//
// 比喻：像小区门禁卡，只有刷卡才能进，没卡的一律拦在门外
//
// 核心机制：
//   - c.Next() = 放行，继续处理
//   - c.Abort() = 拦截，停止后续所有处理
//   - c.Set()   = 往上下文里存数据，后续 handler 可以用 c.Get() 取出
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")

		if token == "" || !strings.HasPrefix(token, "Bearer ") {
			util.Error(c, http.StatusUnauthorized, 401, "未授权：缺少有效的 Token")
			c.Abort() // 拦截！不继续往下走了
			return
		}

		// 简化版：提取 token 作为用户 ID
		// 实际项目这里会做 JWT 解析、校验签名等
		userID := strings.TrimPrefix(token, "Bearer ")
		c.Set("userID", userID) // 存到上下文，后续 handler 可用

		c.Next() // 放行！继续处理下一个中间件或 handler
	}
}
