package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func AddReqID(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		reqID := c.GetHeader("x-request-iD")
		if reqID == "" {
			reqID = uuid.New().String()
		}

		c.Set("x-request-id", reqID)
		logger = logger.With(zap.String("x-request-id", reqID))
		c.Set("logger", logger)

		c.Next()
	}
}
