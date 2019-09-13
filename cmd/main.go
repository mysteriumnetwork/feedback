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
	"github.com/mysteriumnetwork/feedback/docs"
	"github.com/mysteriumnetwork/feedback/feedback"
	"github.com/mysteriumnetwork/feedback/infra"
	"github.com/mysteriumnetwork/feedback/server"
)

const (
	// EnvAWSEndpointURL AWS URL for file upload
	EnvAWSEndpointURL = "AWS_ENDPOINT_URL"
	// EnvAWSBucket AWS bucket for file upload
	EnvAWSBucket = "AWS_BUCKET"
	// EnvGithubAccessToken Github credentials for issue report
	EnvGithubAccessToken = "GITHUB_ACCESS_TOKEN"
	// EnvGithubOwner Github owner of the repository for issue report
	EnvGithubOwner = "GITHUB_OWNER"
	// EnvGithubRepository Github repository for issue report
	EnvGithubRepository = "GITHUB_REPOSITORY"
)

func main() {
	os.Exit(app())
}

func app() (retValue int) {
	configureFromFlags()
	infra.BootstrapLogger(infra.CurrentLogOptions)

	log.Info("Starting feedback service")
	defer func() {
		log.Info("Stopping feedback service")
		log.Flush()
	}()

	err := envPresent(
		EnvAWSEndpointURL,
		EnvAWSBucket,
		EnvGithubAccessToken,
		EnvGithubOwner,
		EnvGithubRepository,
	)
	if err != nil {
		_ = log.Critical(err)
		return -1
	}

	storage, err := feedback.New(&feedback.NewStorageOpts{
		EndpointURL: os.Getenv(EnvAWSEndpointURL),
		Bucket:      os.Getenv(EnvAWSBucket),
	})
	if err != nil {
		_ = log.Critical("Failed to initialize storage: ", err)
		return -1
	}

	githubReporter := feedback.NewReporter(&feedback.NewReporterOpts{
		Token:      os.Getenv(EnvGithubAccessToken),
		Owner:      os.Getenv(EnvGithubOwner),
		Repository: os.Getenv(EnvGithubRepository),
	})
	rateLimiter := infra.NewRateLimiter(0.0166) // 1/minute

	srvr := server.New(
		feedback.NewEndpoint(githubReporter, storage, rateLimiter),
		infra.NewPingEndpoint(),
		docs.NewEndpoint(),
	)

	err = srvr.Serve()
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

func configureFromFlags() {
	infra.RegisterLoggerFlags()
	flag.Parse()
	infra.ConfigureLogger()
}
