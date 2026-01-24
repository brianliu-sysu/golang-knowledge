package observability

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	WSConnections = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ws_connections",
		Help: "Number of active WebSocket connections",
	})

	WSReadFrames = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "ws_read_frames",
		Help: "Number of frames read from WebSocket connections",
	})

	WSSentFrames = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "ws_sent_frames",
		Help: "Number of frames sent to WebSocket connections",
	})

	WSBadProto = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "ws_bad_proto",
		Help: "Number of WebSocket connections with bad protocol",
	})
)

func init() {
	prometheus.MustRegister(WSConnections)
	prometheus.MustRegister(WSReadFrames)
	prometheus.MustRegister(WSSentFrames)
	prometheus.MustRegister(WSBadProto)
}

func MetricsHandler() http.Handler {
	return promhttp.Handler()
}
