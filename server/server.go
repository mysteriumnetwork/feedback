package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	log "github.com/cihub/seelog"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/mysteriumnetwork/feedback/infra/apierror"
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
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, apierror.NewMsg("resource not found").ToResponse())
	})

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
