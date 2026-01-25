package middleware

import (
	"net/http"
	"memo-studio/backend/models"
	"memo-studio/backend/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware 认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未提供认证令牌"})
			c.Abort()
			return
		}

		// 提取 token (Bearer <token>)
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "认证令牌格式错误"})
			c.Abort()
			return
		}

		token := parts[1]
		claims, err := utils.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的认证令牌"})
			c.Abort()
			return
		}

		// 将用户信息存储到上下文
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		// 兼容旧 token：如果 claim 未带 is_admin，则从 DB 兜底一次
		isAdmin := claims.IsAdmin
		if !isAdmin {
			if u, err := models.GetUserByID(claims.UserID); err == nil {
				isAdmin = u.IsAdmin
			}
		}
		c.Set("isAdmin", isAdmin)

		c.Next()
	}
}

// AdminOnly 需要管理员权限
func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		v, ok := c.Get("isAdmin")
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "无权限"})
			c.Abort()
			return
		}
		if b, ok := v.(bool); !ok || !b {
			c.JSON(http.StatusForbidden, gin.H{"error": "无权限"})
			c.Abort()
			return
		}
		c.Next()
	}
}
