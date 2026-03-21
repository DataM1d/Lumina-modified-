package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

func RateLimit(requestsPerMinute int) gin.HandlerFunc {
	var mu sync.Mutex
	clients := make(map[string][]time.Time)

	return func(c *gin.Context) {
		ip := c.ClientIP()
		mu.Lock()
		defer mu.Unlock()

		now := time.Now()
		var valid []time.Time
		for _, t := range clients[ip] {
			if now.Sub(t) < time.Minute {
				valid = append(valid, t)
			}
		}

		if len(valid) >= requestsPerMinute {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded",
			})
			return
		}

		clients[ip] = append(valid, now)
		c.Next()
	}
}
