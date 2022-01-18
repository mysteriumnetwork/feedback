package ci

import (
	"fmt"
	"time"

	"github.com/cihub/seelog"
	"github.com/magefile/mage/sh"
	"github.com/mysteriumnetwork/feedback/di"
	"github.com/mysteriumnetwork/feedback/params"
)

// E2E runs the e2e tests
func E2E() error {
	_, err := startCompose()
	if err != nil {
		return fmt.Errorf("could not deploy containers")
	}
	defer teardownCompose()

	//wait for startup
	time.Sleep(time.Second * 2)

	logLevel := seelog.DebugStr
	skipFileUpload := true
	requestsPerSecond := 9999999999.0
	gparams := params.Generic{
		LogLevelFlag:      &logLevel,
		SkipFileUpload:    &skipFileUpload,
		RequestsPerSecond: &requestsPerSecond,
	}

	envAWSEndpointURL := "http://localhost:8090"
	envAWSBucket := "random-bucket"
	envGithubAccessToken := "github-token"
	envGithubOwner := "github-owner"
	envGithubRepository := "github-repo"
	envIntercomAccessToken := "hi-i-am-a-token!"
	envIntercomBaseURL := "http://localhost:8090"
	eparams := params.Environment{
		EnvAWSEndpointURL:      &envAWSEndpointURL,
		EnvAWSBucket:           &envAWSBucket,
		EnvIntercomAccessToken: &envIntercomAccessToken,
		EnvIntercomBaseURL:     &envIntercomBaseURL,
		EnvGithubAccessToken:   &envGithubAccessToken,
		EnvGithubOwner:         &envGithubOwner,
		EnvGithubRepository:    &envGithubRepository,
	}

	fmt.Println("starting server...")

	container := &di.Container{}
	defer container.Cleanup()

	server, err := container.ConstructServer(gparams, eparams)
	if err != nil {
		return err
	}

	go server.Serve()
	defer server.Stop()

	err = sh.RunV("go",
		"test",
		"-v",
		"-tags=e2e",
		"./e2e",
	)
	return err
}

func startCompose() (string, error) {
	return sh.Output("docker-compose", "-f", "docker-compose.e2e.yml", "up", "-d")
}

func teardownCompose() (string, error) {
	return sh.Output("docker-compose", "-f", "docker-compose.e2e.yml", "down")
}
