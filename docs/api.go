package docs

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Swagger object
type Swagger struct{}

var docsHandler = ginSwagger.WrapHandler(swaggerfiles.Handler)

// NewSwaggerHandler creates a new handler for swagger docs
func NewSwagger() *Swagger {
	return &Swagger{}
}

// Index redirects root route to swagger docs
func (s *Swagger) Index(context *gin.Context) {
	context.Redirect(301, "/swagger/index.html")
}

// Docs use ginSwagger middleware to serve the API docs
func (s *Swagger) Docs(context *gin.Context) {
	docsHandler(context)
}

func (s *Swagger) AttachHandlers(g *gin.Engine) {
	g.GET("/", s.Index)
	g.GET("/swagger/*any", s.Docs)
}
