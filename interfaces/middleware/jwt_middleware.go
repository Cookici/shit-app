package middleware

import (
	"net/http"
	"record-project/infrastructure/auth"
	"strings"

	"github.com/gin-gonic/gin"
)

// JWTAuthMiddleware JWT认证中间件
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未提供认证令牌"})
			c.Abort()
			return
		}

		// 检查Bearer前缀
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "认证格式错误"})
			c.Abort()
			return
		}

		// 解析token
		claims, err := auth.ParseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的认证令牌: " + err.Error()})
			c.Abort()
			return
		}

		// 将用户信息保存到上下文
		c.Set("userID", claims.UserID)
		c.Set("openID", claims.OpenID)
		c.Set("nickname", claims.Nickname)

		c.Next()
	}
}