package params

import (
	"flag"
	"fmt"
	"strings"

	"github.com/cihub/seelog"
	"github.com/mysteriumnetwork/feedback/infra"
)

// Generic represents all the generic parameters
type Generic struct {
	LogLevelFlag      *string
	RequestsPerSecond *float64
}

// Init initialized the generic parameters with flags
func (f *Generic) Init() {
	f.LogLevelFlag = flag.String("log-level", seelog.DebugStr, fmt.Sprintf("Service logging level (%s)", strings.Join(infra.AllLevels, "|")))
	f.RequestsPerSecond = flag.Float64("requestsPerSecond", 0.0166, "Maximum number of requests per second (default: 1 per minute)")
}
