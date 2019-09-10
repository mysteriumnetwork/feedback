package server

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

// swagger:operation GET /ping ping
// ---
// summary: Ping responds to ping
// responses:
//   '200':
//     description: Ping successful
//     schema:
//       "$ref": "#/definitions/Ping"
//
func (p *PingEndpoint) Ping(c *gin.Context) {
	c.JSON(200, Ping{"pong"})
}

// RegisterRoutes registers ping route
func (p *PingEndpoint) RegisterRoutes(r gin.IRoutes) {
	r.GET("/ping", p.Ping)
}
