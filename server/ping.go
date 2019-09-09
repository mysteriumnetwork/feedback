package server

import "github.com/gin-gonic/gin"

// Ping represents ping response
type Ping struct {
	Message string `json:"message"`
}

// PingEndpoint responds to ping
type PingEndpoint struct {
}

// NewPingEndpoint creates ping endpoint
func NewPingEndpoint() *PingEndpoint {
	return &PingEndpoint{}
}

// Ping responds to ping
func (p *PingEndpoint) Ping(c *gin.Context) {
	c.JSON(200, Ping{"pong"})
}

// RegisterRoutes registers ping route
func (p *PingEndpoint) RegisterRoutes(r gin.IRoutes) {
	r.GET("/ping", p.Ping)
}
