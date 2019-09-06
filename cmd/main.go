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

package main

import (
	"flag"
	"os"

	log "github.com/cihub/seelog"
	"github.com/mysteriumnetwork/feedback/api"
	"github.com/mysteriumnetwork/feedback/feedback"
	"github.com/mysteriumnetwork/feedback/logconfig"
)

func main() {
	os.Exit(app())
}

func app() (retValue int) {
	configureFromFlags()
	logconfig.BootstrapWith(logconfig.CurrentLogOptions)

	log.Info("Starting feedback service")
	defer func() {
		log.Info("Stopping feedback service")
		log.Flush()
	}()

	server := api.NewServer(
		feedback.NewEndpoint(),
	)

	err := server.Serve()
	if err != nil {
		_ = log.Critical("Critical error occurred: ", err)
		return -1
	}
	return 0
}

func configureFromFlags() {
	logconfig.RegisterFlags()
	flag.Parse()
	logconfig.Configure()
}
