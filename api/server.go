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

package api

import (
	"errors"
	"net/http"
	"time"

	log "github.com/cihub/seelog"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/mysteriumnetwork/feedback/apierr"
)

type routes interface {
	RegisterRoutes(r gin.IRoutes)
}

// Server represents API server
type Server struct {
	routes []routes
}

// NewServer creates a new API server
func NewServer(routes ...routes) *Server {
	return &Server{
		routes: routes,
	}
}

// Serve starts API server
func (s *Server) Serve() error {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(Logger())
	r.Use(cors.Default())
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, apierr.Single(errors.New("resource not found")))
	})

	v1 := r.Group("/api/v1")
	{
		for i := range s.routes {
			s.routes[i].RegisterRoutes(v1)
		}
	}

	return r.Run()
}

// Logger forces gin to use our logger
// Adapted from gin.Logger
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		end := time.Now()
		latency := end.Sub(start)

		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		comment := c.Errors.ByType(gin.ErrorTypePrivate).String()

		if raw != "" {
			path = path + "?" + raw
		}

		log.Infof("[GIN] %3d | %13v | %15s | %-7s | %v %v",
			statusCode,
			latency,
			clientIP,
			method,
			path,
			comment,
		)
	}
}
