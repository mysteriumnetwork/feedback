// Package main Feedback API
//
// The purpose of this documentation is to provide developers an insight of how to
// interact with Mysterium Feedback API
//
//     schemes: https
//     host: localhost
//     basepath: /api/v1
//     license: GPLv3 https://www.gnu.org/licenses/gpl-3.0.html
//
//     consumes:
//       - application/json
//
//     produces:
//       - application/json
//
//     version: 0.0.1
//
// swagger:meta
package main

import (
	"errors"
	"flag"
	"os"

	log "github.com/cihub/seelog"
	"github.com/mysteriumnetwork/feedback/constants"
	"github.com/mysteriumnetwork/feedback/di"
	"github.com/mysteriumnetwork/feedback/infra"
	"github.com/mysteriumnetwork/feedback/params"
)

func main() {
	os.Exit(app())
}

func app() (retValue int) {
	gparams := params.Generic{}
	gparams.Init()

	err := envPresent(
		constants.EnvAWSEndpointURL,
		constants.EnvAWSBucket,
		constants.EnvGithubAccessToken,
		constants.EnvGithubOwner,
		constants.EnvGithubRepository,
		constants.EnvIntercomAccessToken,
	)
	if err != nil {
		_ = log.Critical(err)
		return -1
	}

	eparams := params.Environment{}
	eparams.Init()

	flag.Parse()
	infra.ConfigureLogger(*gparams.LogLevelFlag)
	infra.BootstrapLogger(infra.CurrentLogOptions)

	log.Info("Starting feedback service")
	defer func() {
		log.Info("Stopping feedback service")
		log.Flush()
	}()

	container := &di.Container{}
	defer container.Cleanup()

	server, err := container.ConstructServer(gparams, eparams)
	if err != nil {
		_ = log.Critical("Error constructing API server: ", err)
		return -1
	}

	err = server.Serve()
	if err != nil {
		_ = log.Critical("Error running API server: ", err)
		return -1
	}
	return 0
}

func envPresent(vars ...string) error {
	for _, envkey := range vars {
		_, found := os.LookupEnv(envkey)
		if !found {
			return errors.New("required environment variable is not set: " + envkey)
		}
	}
	return nil
}
