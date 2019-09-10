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
