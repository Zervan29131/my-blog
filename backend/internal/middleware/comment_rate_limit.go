package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type rateLimitWindow struct {
	startedAt time.Time
	count     int
}

type CommentRateLimiter struct {
	mutex       sync.Mutex
	clients     map[string]rateLimitWindow
	limit       int
	window      time.Duration
	lastCleanup time.Time
	now         func() time.Time
}

func NewCommentRateLimiter(limit int, window time.Duration) *CommentRateLimiter {
	return &CommentRateLimiter{
		clients:     make(map[string]rateLimitWindow),
		limit:       limit,
		window:      window,
		lastCleanup: time.Now(),
		now:         time.Now,
	}
}

func (limiter *CommentRateLimiter) Handler() gin.HandlerFunc {
	return func(context *gin.Context) {
		if limiter.allow(context.ClientIP()) {
			context.Next()
			return
		}

		context.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
			"error": gin.H{
				"code":    "RATE_LIMITED",
				"message": "评论提交过于频繁，请稍后再试",
				"details": gin.H{},
			},
		})
	}
}

func (limiter *CommentRateLimiter) allow(clientIP string) bool {
	limiter.mutex.Lock()
	defer limiter.mutex.Unlock()

	now := limiter.now()
	if now.Sub(limiter.lastCleanup) >= limiter.window {
		for ip, client := range limiter.clients {
			if now.Sub(client.startedAt) >= limiter.window {
				delete(limiter.clients, ip)
			}
		}
		limiter.lastCleanup = now
	}

	client, exists := limiter.clients[clientIP]
	if !exists || now.Sub(client.startedAt) >= limiter.window {
		limiter.clients[clientIP] = rateLimitWindow{startedAt: now, count: 1}
		return true
	}
	if client.count >= limiter.limit {
		return false
	}
	client.count++
	limiter.clients[clientIP] = client
	return true
}
