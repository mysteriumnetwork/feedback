package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/mysteriumnetwork/feedback/docs"
	"github.com/mysteriumnetwork/feedback/infra/apierror"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
)

type routes interface {
	RegisterRoutes(r gin.IRoutes)
}

// Server represents API server
type Server struct {
	routes   []routes
	stopOnce sync.Once
	stop     chan struct{}
}

// New creates a new API server
func New(routes ...routes) *Server {
	return &Server{
		routes: routes,
		stop:   make(chan struct{}),
	}
}

// Serve starts API server
func (s *Server) Serve() error {
	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(Logger("/api/v1/ping"))
	r.Use(cors.Default())
	r.Use(HandlerMetrics)

	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, apierror.NewMsg("resource not found").ToResponse())
	})
	docs.NewSwagger().AttachHandlers(r)

	v1 := r.Group("/api/v1")
	{
		for i := range s.routes {
			s.routes[i].RegisterRoutes(v1)
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%v", port),
		Handler: r,
	}

	go func() {
		<-s.stop
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		srv.Shutdown(ctx)
	}()

	return srv.ListenAndServe()
}

// Stop stops the server
func (s *Server) Stop() {
	s.stopOnce.Do(func() {
		close(s.stop)
	})
}

// Logger forces gin to use our logger
// Adapted from gin.Logger
func Logger(ignorePaths ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		for _, igp := range ignorePaths {
			if path == igp {
				return
			}
		}

		end := time.Now()
		latency := end.Sub(start)

		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		comment := c.Errors.ByType(gin.ErrorTypePrivate).String()

		if raw != "" {
			path = path + "?" + raw
		}

		log.Debug().Int("status", statusCode).Str("method", method).
			Str("path", path).Str("client_ip", clientIP).
			Dur("latency", latency).Str("comment", comment).
			Msg("gin request logged")
	}
}
