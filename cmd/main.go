// Package main Feedback API
//
// The purpose of this documentation is to provide developers an insight of how to
// interact with Mysterium Feedback API
//
//	schemes: https
//	host: localhost
//	basepath: /api/v1
//	license: GPLv3 https://www.gnu.org/licenses/gpl-3.0.html
//
//	consumes:
//	  - application/json
//
//	produces:
//	  - application/json
//
//	version: 0.0.1
//
// swagger:meta
package main

import (
	"errors"
	"flag"
	"os"

	"github.com/mysteriumnetwork/feedback/constants"
	"github.com/mysteriumnetwork/feedback/di"
	"github.com/mysteriumnetwork/feedback/params"
	mlog "github.com/mysteriumnetwork/logger"
	"github.com/rs/zerolog/log"
)

// @title Feedback
// @version 1.0
// @description This is a service dedicated to collecting feedback from Mysterium Network users
// @termsOfService https://docs.mysterium.network/en/latest/about/terms-and-conditions/

// @contact.name API Support
// @contact.url https://github.com/mysteriumnetwork/feedback/issues

// @BasePath /api
func main() {
	os.Exit(app())
}

func app() (retValue int) {
	logger := mlog.BootstrapDefaultLogger()

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
		log.Fatal().Err(err).Msg("Failed to initialize environment")
		return -1
	}

	eparams := params.Environment{}
	eparams.Init()

	flag.Parse()
	mlog.SetLevel(logger, *gparams.LogLevelFlag)

	log.Info().Msg("Starting feedback service")
	defer func() {
		log.Info().Msg("Stopping feedback service")
	}()

	container := &di.Container{}
	defer container.Cleanup()

	server, err := container.ConstructServer(gparams, eparams)
	if err != nil {
		log.Fatal().Err(err).Msg("Error constructing API server")
		return -1
	}

	err = server.Serve()
	if err != nil {
		log.Fatal().Err(err).Msg("Error running API server")
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
