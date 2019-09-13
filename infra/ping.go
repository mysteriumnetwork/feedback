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

// Ping responds to ping
// swagger:operation GET /ping ping
// ---
// summary: Responds to ping
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
