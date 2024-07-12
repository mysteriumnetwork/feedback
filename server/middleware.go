package server

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mysteriumnetwork/feedback/metrics"
)

func HandlerMetrics(c *gin.Context) {
	now := time.Now()
	c.Next()
	fullPath := c.FullPath()

	code := fmt.Sprint(c.Writer.Status())

	metrics.HTTPResponseTime.
		WithLabelValues(fullPath, c.Request.Method, code).
		Observe(time.Since(now).Seconds())
}
