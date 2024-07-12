package metrics

import "github.com/prometheus/client_golang/prometheus"

var HTTPResponseTime = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name: "feedback_response_time",
		Help: "Response time of HTTP endpoints",
	},
	[]string{"request_url", "method", "response_code"},
)

func InitMetrics() error {
	return prometheus.Register(HTTPResponseTime)
}
