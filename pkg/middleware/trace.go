package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// TraceIDMiddleware 生成并注入 Trace ID
func TraceIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := uuid.New().String()
		c.Set("TraceID", traceID)
		// 将 TraceID 添加到响应头
		c.Writer.Header().Set("X-Trace-Id", traceID)
		c.Next()
	}
}
