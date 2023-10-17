package e2e

import (
	"os"
	"testing"

	"github.com/mysteriumnetwork/feedback/client"
	"github.com/rs/zerolog/log"
)

var (
	apiClient *client.FeedbackAPI
)

func TestMain(m *testing.M) {
	var err error
	apiClient, err = client.NewFeedbackAPI("http://localhost:8080")
	if err != nil {
		log.Error().Err(err).Msg("could not create feedback api client")
		os.Exit(-1)
	}

	os.Exit(m.Run())
}
