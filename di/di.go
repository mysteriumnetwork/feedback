package di

import (
	"sync"

	log "github.com/cihub/seelog"
	"github.com/mysteriumnetwork/feedback/docs"
	"github.com/mysteriumnetwork/feedback/feedback"
	"github.com/mysteriumnetwork/feedback/infra"
	"github.com/mysteriumnetwork/feedback/params"
	"github.com/mysteriumnetwork/feedback/server"
)

// Container represents our dependency container
type Container struct {
	cleanup []func()
	lock    sync.Mutex
}

// ConstructServer creates a server for us
func (c *Container) ConstructServer(gparams params.Generic, eparams params.Environment) (*server.Server, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	storage, err := feedback.New(&feedback.NewStorageOpts{
		EndpointURL: *eparams.EnvAWSEndpointURL,
		Bucket:      *eparams.EnvAWSBucket,
	})
	if err != nil {
		_ = log.Critical("Failed to initialize storage: ", err)
		return nil, err
	}

	githubReporter := feedback.NewGithubReporter(&feedback.NewGithubReporterOpts{
		Token:      *eparams.EnvGithubAccessToken,
		Owner:      *eparams.EnvGithubOwner,
		Repository: *eparams.EnvGithubRepository,
	})
	rateLimiter := infra.NewRateLimiter(*gparams.RequestsPerSecond)

	intercomReporter := feedback.NewIntercomReporter(&feedback.NewIntercomReporterOpts{
		Token:           *eparams.EnvIntercomAccessToken,
		IntercomBaseURL: *eparams.EnvIntercomBaseURL,
	})

	srvr := server.New(
		feedback.NewEndpoint(githubReporter, intercomReporter, storage, rateLimiter),
		infra.NewPingEndpoint(),
		docs.NewEndpoint(),
	)

	return srvr, nil
}

// Cleanup performs the cleanup required
func (c *Container) Cleanup() {
	c.lock.Lock()
	defer c.lock.Unlock()
	for i := len(c.cleanup) - 1; i >= 0; i-- {
		c.cleanup[i]()
	}
}
