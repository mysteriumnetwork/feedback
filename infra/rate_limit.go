/*
 * Copyright (C) 2019 The "MysteriumNetwork/feedback" Authors.
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

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
