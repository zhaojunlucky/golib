// Package httpmetrics exposes Prometheus HTTP server metrics for a plain
// net/http (or http.ServeMux) service, named and labeled the same way as
// pkg/ginmetrics so all Go services share one convention: metric names and
// labels matching Spring Boot / Micrometer's http_server_requests_seconds,
// with services distinguished by an "application" label (mirroring Spring
// Boot's management.metrics.tags.application) instead of a per-service
// metric-name prefix.
package httpmetrics

import (
	"net/http"
	"strconv"
	"strings"
	"time"

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

// New builds and registers the HTTP server metrics for a net/http service.
// It panics if Options.Application is empty, since every metric depends on
// that label to be distinguishable from other services.
func New(opts Options) *Metrics {
	if strings.TrimSpace(opts.Application) == "" {
		panic("httpmetrics: Options.Application is required")
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

// routePatterner is implemented by *http.ServeMux: it reports which
// registered pattern would handle a request without dispatching it, which
// Middleware uses as the low-cardinality "uri" label (analogous to gin's
// c.FullPath()).
type routePatterner interface {
	Handler(r *http.Request) (http.Handler, string)
}

// Middleware instruments every request served through mux with the
// registered metrics, skipping the configured metrics endpoints
// themselves. mux is consulted (via its Handler method) only to look up
// the matched route pattern for the "uri" label - it does not have to be
// the handler being wrapped.
func (m *Metrics) Middleware(mux routePatterner) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL != nil {
				if _, ok := m.skipPaths[r.URL.Path]; ok {
					next.ServeHTTP(w, r)
					return
				}
			}

			method := r.Method
			uri := uriLabel(mux, r)
			m.inFlight.WithLabelValues(method, uri).Inc()
			start := time.Now()
			rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}

			defer func() {
				m.inFlight.WithLabelValues(method, uri).Dec()

				status := rec.status
				finalURI := uri
				if finalURI == "UNKNOWN" && status == http.StatusNotFound {
					finalURI = "NOT_FOUND"
				}
				m.requestDuration.
					WithLabelValues(method, finalURI, strconv.Itoa(status), outcomeLabel(status)).
					Observe(time.Since(start).Seconds())
			}()

			next.ServeHTTP(rec, r)
		})
	}
}

type statusRecorder struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func (r *statusRecorder) WriteHeader(status int) {
	if !r.wroteHeader {
		r.status = status
		r.wroteHeader = true
	}
	r.ResponseWriter.WriteHeader(status)
}

func (r *statusRecorder) Write(b []byte) (int, error) {
	if !r.wroteHeader {
		r.wroteHeader = true
	}
	return r.ResponseWriter.Write(b)
}

// uriLabel resolves the low-cardinality route template for r via mux's
// pattern lookup. http.ServeMux patterns are "[METHOD ][HOST]/PATH"; the
// method (already captured in its own label) is stripped, keeping the
// host+path portion. Falls back to "UNKNOWN" when nothing matches.
func uriLabel(mux routePatterner, r *http.Request) string {
	_, pattern := mux.Handler(r)
	if pattern == "" {
		return "UNKNOWN"
	}
	if i := strings.IndexByte(pattern, ' '); i >= 0 {
		pattern = pattern[i+1:]
	}
	if pattern == "" {
		return "UNKNOWN"
	}
	return pattern
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
