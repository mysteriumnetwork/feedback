package infra

import (
	"github.com/didip/tollbooth/v5"
	"github.com/didip/tollbooth/v5/limiter"
	"github.com/gin-gonic/gin"
	"github.com/mysteriumnetwork/feedback/infra/apierror"
)

// RateLimiter limits requests to certain endpoints
type RateLimiter struct {
	limiter *limiter.Limiter
}

// NewRateLimiter creates a new RateLimiter
func NewRateLimiter(requestsPerSecond float64) *RateLimiter {
	lmt := tollbooth.NewLimiter(requestsPerSecond, &limiter.ExpirableOptions{})
	return &RateLimiter{
		limiter: lmt,
	}
}

// Handler returns handler func for gin
func (r *RateLimiter) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		httpErr := tollbooth.LimitByRequest(r.limiter, c.Writer, c.Request)
		if httpErr != nil {
			c.JSON(httpErr.StatusCode, apierror.NewMsg(httpErr.Message).ToResponse())
			c.Abort()
		} else {
			c.Next()
		}
	}
}
