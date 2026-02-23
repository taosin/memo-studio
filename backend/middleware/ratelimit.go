package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiter 简单速率限制器
type RateLimiter struct {
	requests map[string][]time.Time
	mu       sync.RWMutex
	limit    int           // 每时间窗口最大请求数
	window   time.Duration // 时间窗口
}

// NewRateLimiter 创建速率限制器
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}
}

// isAllowed 检查请求是否允许
func (rl *RateLimiter) isAllowed(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-rl.window)

	// 清理过期记录
	if times, exists := rl.requests[key]; exists {
		var valid []time.Time
		for _, t := range times {
			if t.After(windowStart) {
				valid = append(valid, t)
			}
		}
		if len(valid) < len(times) {
			rl.requests[key] = valid
		}
	}

	// 检查是否超过限制
	if len(rl.requests[key]) >= rl.limit {
		return false
	}

	// 记录请求
	rl.requests[key] = append(rl.requests[key], now)
	return true
}

// GetRemaining 获取剩余请求数
func (rl *RateLimiter) getRemaining(key string) int {
	rl.mu.RLock()
	defer rl.mu.RUnlock()

	now := time.Now()
	windowStart := now.Add(-rl.window)

	count := 0
	if times, exists := rl.requests[key]; exists {
		for _, t := range times {
			if t.After(windowStart) {
				count++
			}
		}
	}

	remaining := rl.limit - count
	if remaining < 0 {
		remaining = 0
	}
	return remaining
}

var (
	globalLimiter     *RateLimiter
	globalLimiterOnce sync.Once
)

// GetGlobalLimiter 获取全局速率限制器
func GetGlobalLimiter() *RateLimiter {
	globalLimiterOnce.Do(func() {
		globalLimiter = NewRateLimiter(50, time.Minute) // 每分钟50次
	})
	return globalLimiter
}

// RateLimitMiddleware 返回速率限制中间件
func RateLimitMiddleware() gin.HandlerFunc {
	limiter := GetGlobalLimiter()
	
	return func(c *gin.Context) {
		// 获取客户端 IP 作为 key
		clientIP := c.ClientIP()
		
		if !limiter.isAllowed(clientIP) {
			c.Header("Retry-After", "60")
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "请求过于频繁，请稍后再试",
				"code":  "RATE_LIMIT_EXCEEDED",
			})
			c.Abort()
			return
		}

		// 添加速率限制头信息
		remaining := limiter.getRemaining(clientIP)
		c.Header("X-RateLimit-Limit", "50")
		c.Header("X-RateLimit-Remaining", string(rune('0'+remaining)))
		
		c.Next()
	}
}

// StrictRateLimitMiddleware 严格速率限制（每分钟30次）
func StrictRateLimitMiddleware() gin.HandlerFunc {
	limiter := NewRateLimiter(30, time.Minute)
	
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		
		if !limiter.isAllowed(clientIP) {
			c.Header("Retry-After", "60")
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "请求过于频繁，请稍后再试",
				"code":  "RATE_LIMIT_EXCEEDED",
			})
			c.Abort()
			return
		}
		
		c.Next()
	}
}
