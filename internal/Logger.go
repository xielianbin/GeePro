package internal

import (
	"log"
	"time"
)

// 中间件
func Logger() HandlerFunc {
	return func(c IContext) {
		// Start timer
		t := time.Now()
		// Process request
		c.Next()
		// Calculate resolution time
		log.Printf("[%d] %s in %v", c.GetStatus(), c.GetRequest().RequestURI, time.Since(t))
	}
}
