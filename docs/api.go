package docs

import (
	"github.com/gin-gonic/gin"
)

// Endpoint API documentation endpoint
type Endpoint struct {
}

// NewEndpoint creates new Endpoint
func NewEndpoint() *Endpoint {
	return &Endpoint{}
}

// SwaggerJSON return API schema
func (e *Endpoint) SwaggerJSON(c *gin.Context) {
	c.Writer.Write(MustAsset(SwaggerJSONFilepath))
	c.Status(200)
}

// RegisterRoutes registers API documentation routes
func (e *Endpoint) RegisterRoutes(r gin.IRoutes) {
	r.GET("/swagger.json", e.SwaggerJSON)
}
