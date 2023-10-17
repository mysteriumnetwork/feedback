package ci

import (
	"fmt"
	"os"
	"time"

	"github.com/magefile/mage/sh"
	"github.com/mysteriumnetwork/feedback/di"
	"github.com/mysteriumnetwork/feedback/params"
	"github.com/rs/zerolog"
)

// E2E runs the e2e tests
func E2E() error {
	_, err := startCompose()
	if err != nil {
		return fmt.Errorf("could not deploy containers")
	}
	defer teardownCompose()

	// wait for startup
	// TODO: ping s3-mock and wiremock and await startup, it is very slow in the pipeline
	time.Sleep(time.Second * 10)

	logLevel := zerolog.DebugLevel.String()
	requestsPerSecond := 9999999999.0
	logProxyBaseUrl := "http://someweb.com"
	githubBaseUrlOverride := "http://localhost:8090/github/"
	gparams := params.Generic{
		LogLevelFlag:          &logLevel,
		RequestsPerSecond:     &requestsPerSecond,
		LogProxyBaseUrl:       &logProxyBaseUrl,
		GithubBaseUrlOverride: &githubBaseUrlOverride,
	}

	envAWSEndpointURL := "http://localhost:9090"
	envAWSBucket := "node-user-reports"
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

	//set mock env vars
	err = os.Setenv("AWS_ACCESS_KEY_ID", "test_key_id")
	if err != nil {
		return err
	}
	err = os.Setenv("AWS_SECRET_ACCESS_KEY", "test_secret_key")
	if err != nil {
		return err
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
