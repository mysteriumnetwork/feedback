package params

import (
	"os"
	"strings"

	"github.com/mysteriumnetwork/feedback/constants"
)

// Environment represents all the environment variables
type Environment struct {
	EnvAWSEndpointURL *string
	// EnvAWSBucket AWS bucket for file upload
	EnvAWSBucket *string
	// EnvGithubAccessToken Github credentials for issue report
	EnvGithubAccessToken *string
	// EnvGithubOwner Github owner of the repository for issue report
	EnvGithubOwner *string
	// EnvGithubRepository Github repository for issue report
	EnvGithubRepository *string
	// EnvIntercomBaseURL Intercom base URL for API
	EnvIntercomBaseURL *string
	// EnvIntercomAccessToken Intercom token for API
	EnvIntercomAccessToken *string
}

// Init initializes all the environment variables
func (f *Environment) Init() {
	envAWSEndpointURL := os.Getenv(constants.EnvAWSEndpointURL)
	f.EnvAWSEndpointURL = &envAWSEndpointURL
	envAWSBucket := os.Getenv(constants.EnvAWSBucket)
	f.EnvAWSBucket = &envAWSBucket
	envGithubAccessToken := os.Getenv(constants.EnvGithubAccessToken)
	f.EnvGithubAccessToken = &envGithubAccessToken
	envGithubOwner := os.Getenv(constants.EnvGithubOwner)
	f.EnvGithubOwner = &envGithubOwner
	envGithubRepository := os.Getenv(constants.EnvGithubRepository)
	f.EnvGithubRepository = &envGithubRepository
	// remove trailing dash
	envIntercomBaseURL := strings.TrimSuffix(os.Getenv(constants.EnvIntercomBaseURL), "/")
	f.EnvIntercomBaseURL = &envIntercomBaseURL
	envIntercomAccessToken := os.Getenv(constants.EnvIntercomAccessToken)
	f.EnvIntercomAccessToken = &envIntercomAccessToken
}
