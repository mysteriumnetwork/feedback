package params

import (
	"flag"

	"github.com/rs/zerolog"
)

// Generic represents all the generic parameters
type Generic struct {
	LogLevelFlag          *string
	RequestsPerSecond     *float64
	LogProxyBaseUrl       *string
	GithubBaseUrlOverride *string
}

// Init initialized the generic parameters with flags
func (f *Generic) Init() {
	f.LogLevelFlag = flag.String("log-level", zerolog.DebugLevel.String(), "Service logging level")
	f.RequestsPerSecond = flag.Float64("requestsPerSecond", 0.0166, "Maximum number of requests per second (default: 1 per minute)")
	f.LogProxyBaseUrl = flag.String("logProxyBaseUrl", "https://mystnodes.com/support/logs", "A url which serves as proxy for accessing logs easily by support.")
}
