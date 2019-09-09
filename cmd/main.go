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
	"errors"
	"flag"
	"os"

	log "github.com/cihub/seelog"
	"github.com/mysteriumnetwork/feedback/feedback"
	"github.com/mysteriumnetwork/feedback/logconfig"
	"github.com/mysteriumnetwork/feedback/server"
	"github.com/mysteriumnetwork/feedback/storage"
	"github.com/mysteriumnetwork/feedback/target/github"
)

const (
	// ENV_AWS_ENDPOINT_URL AWS URL for file upload
	ENV_AWS_ENDPOINT_URL = "AWS_ENDPOINT_URL"
	// ENV_AWS_BUCKET AWS bucket for file upload
	ENV_AWS_BUCKET = "AWS_BUCKET"
	// ENV_GITHUB_ACCESS_TOKEN Github credentials for issue report
	ENV_GITHUB_ACCESS_TOKEN = "GITHUB_ACCESS_TOKEN"
	// ENV_GITHUB_OWNER Github owner of the repository for issue report
	ENV_GITHUB_OWNER = "GITHUB_OWNER"
	// ENV_GITHUB_REPOSITORY Github repository for issue report
	ENV_GITHUB_REPOSITORY = "GITHUB_REPOSITORY"
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

	err := envPresent(
		ENV_AWS_ENDPOINT_URL,
		ENV_AWS_BUCKET,
		ENV_GITHUB_ACCESS_TOKEN,
		ENV_GITHUB_OWNER,
		ENV_GITHUB_REPOSITORY,
	)
	if err != nil {
		_ = log.Critical(err)
		return -1
	}

	storage, err := storage.New(&storage.NewStorageOpts{
		EndpointURL: os.Getenv(ENV_AWS_ENDPOINT_URL),
		Bucket:      os.Getenv(ENV_AWS_BUCKET),
	})
	if err != nil {
		_ = log.Critical("Failed to initialize storage: ", err)
		return -1
	}

	githubReporter := github.NewReporter(&github.NewReporterOpts{
		Token:      os.Getenv(ENV_GITHUB_ACCESS_TOKEN),
		Owner:      os.Getenv(ENV_GITHUB_OWNER),
		Repository: os.Getenv(ENV_GITHUB_REPOSITORY),
	})

	srvr := server.New(
		feedback.NewEndpoint(githubReporter, storage),
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
	logconfig.RegisterFlags()
	flag.Parse()
	logconfig.Configure()
}
