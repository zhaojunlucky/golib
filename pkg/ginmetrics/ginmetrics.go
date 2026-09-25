// Package ginmetrics exposes Prometheus HTTP server metrics for a gin
// application, named and labeled to match Spring Boot / Micrometer's
// http_server_requests_seconds so the same dashboards and alert rules
// can be reused across Go and Spring Boot services. Instead of a
// per-service metric-name prefix, services are distinguished by an
// "application" label applied to every metric, mirroring Spring Boot's
// management.metrics.tags.application.
package ginmetrics

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// DefaultMetricsPath and DefaultAPIMetricsPath are excluded from
// instrumentation by default since they are the metrics endpoint itself.
const (
	DefaultMetricsPath    = "/metrics"
	DefaultAPIMetricsPath = "/api/metrics"
)

var defaultBuckets = []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10}

// Options configures a Metrics instance.
type Options struct {
	// Application is the value of the "application" label applied to every
	// metric in the registry, including the Go/process/build-info
	// collectors. Required.
	Application string

	// Env is reported as the "env" label on app_info (e.g. "production").
	// Optional; defaults to "unknown".
	Env string

	// Buckets overrides the request-duration histogram buckets. Optional;
	// defaults to a set spanning 5ms-10s.
	Buckets []float64

	// SkipPaths lists additional request paths excluded from
	// instrumentation, besides DefaultMetricsPath and DefaultAPIMetricsPath.
	SkipPaths []string
}

// Metrics holds a dedicated Prometheus registry and the HTTP server
// metrics registered on it.
type Metrics struct {
	registry        *prometheus.Registry
	requestDuration *prometheus.HistogramVec
	inFlight        *prometheus.GaugeVec
	skipPaths       map[string]struct{}
}

// New builds and registers the HTTP server metrics for a gin application.
// It panics if Options.Application is empty, since every metric depends on
// that label to be distinguishable from other services.
func New(opts Options) *Metrics {
	if strings.TrimSpace(opts.Application) == "" {
		panic("ginmetrics: Options.Application is required")
	}

	buckets := opts.Buckets
	if len(buckets) == 0 {
		buckets = defaultBuckets
	}

	base := prometheus.NewRegistry()
	registerer := prometheus.WrapRegistererWith(prometheus.Labels{"application": opts.Application}, base)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_server_requests_seconds",
			Help:    "Duration of HTTP server requests in seconds.",
			Buckets: buckets,
		},
		[]string{"method", "uri", "status", "outcome"},
	)

	inFlight := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "http_server_active_requests",
			Help: "Current number of in-flight HTTP server requests.",
		},
		[]string{"method", "uri"},
	)

	appInfo := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "app_info",
			Help: "Static information about the application instance.",
		},
		[]string{"env"},
	)

	registerer.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		collectors.NewBuildInfoCollector(),
		requestDuration,
		inFlight,
		appInfo,
	)
	appInfo.WithLabelValues(envLabel(opts.Env)).Set(1)

	skipPaths := map[string]struct{}{
		DefaultMetricsPath:    {},
		DefaultAPIMetricsPath: {},
	}
	for _, p := range opts.SkipPaths {
		skipPaths[p] = struct{}{}
	}

	return &Metrics{
		registry:        base,
		requestDuration: requestDuration,
		inFlight:        inFlight,
		skipPaths:       skipPaths,
	}
}

// Handler serves the registry in the Prometheus exposition format.
func (m *Metrics) Handler() http.Handler {
	return promhttp.HandlerFor(m.registry, promhttp.HandlerOpts{})
}

// Middleware instruments every gin request with the registered metrics,
// skipping the configured metrics endpoints themselves.
func (m *Metrics) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request != nil && c.Request.URL != nil {
			if _, ok := m.skipPaths[c.Request.URL.Path]; ok {
				c.Next()
				return
			}
		}

		method := c.Request.Method
		uri := initialURILabel(c)
		m.inFlight.WithLabelValues(method, uri).Inc()
		start := time.Now()

		defer func() {
			m.inFlight.WithLabelValues(method, uri).Dec()

			status := c.Writer.Status()
			finalURI := finalURILabel(c, uri)
			m.requestDuration.
				WithLabelValues(method, finalURI, strconv.Itoa(status), outcomeLabel(status)).
				Observe(time.Since(start).Seconds())
		}()

		c.Next()
	}
}

func initialURILabel(c *gin.Context) string {
	if route := c.FullPath(); route != "" {
		return route
	}
	return "UNKNOWN"
}

func finalURILabel(c *gin.Context, fallback string) string {
	if route := c.FullPath(); route != "" {
		return route
	}
	if c.Writer.Status() == http.StatusNotFound {
		return "NOT_FOUND"
	}
	return fallback
}

// outcomeLabel mirrors Spring Boot's http_server_requests_seconds "outcome"
// tag (SUCCESS, CLIENT_ERROR, ...) derived from the HTTP status class.
func outcomeLabel(status int) string {
	switch {
	case status >= 100 && status < 200:
		return "INFORMATIONAL"
	case status >= 200 && status < 300:
		return "SUCCESS"
	case status >= 300 && status < 400:
		return "REDIRECTION"
	case status >= 400 && status < 500:
		return "CLIENT_ERROR"
	case status >= 500 && status < 600:
		return "SERVER_ERROR"
	default:
		return "UNKNOWN"
	}
}

func envLabel(env string) string {
	trimmed := strings.TrimSpace(env)
	if trimmed == "" {
		return "unknown"
	}
	return trimmed
}
