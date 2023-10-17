package infra

import "github.com/gin-gonic/gin"

// Ping represents ping response
// swagger:model
type Ping struct {
	// example: pong
	Message string `json:"message"`
}

// PingEndpoint responds to ping
type PingEndpoint struct {
}

// NewPingEndpoint creates ping endpoint
func NewPingEndpoint() *PingEndpoint {
	return &PingEndpoint{}
}

// Ping check service health
//
// @Tags health
// @Summary Check service health
// @Description Endpoint is meant to check service health and returns pong response for every request
// @Accept  json
// @Produce  json
// @Success 200 {object} infra.Ping
// @Router /v1/ping [get]
func (p *PingEndpoint) Ping(c *gin.Context) {
	c.JSON(200, Ping{"pong"})
}

// RegisterRoutes registers ping route
func (p *PingEndpoint) RegisterRoutes(r gin.IRoutes) {
	r.GET("/ping", p.Ping)
}
