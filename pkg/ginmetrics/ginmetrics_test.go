package ginmetrics

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestNewRequiresApplication(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("expected New to panic when Application is empty")
		}
	}()
	New(Options{})
}

func newTestRouter(t *testing.T, opts Options) (*gin.Engine, *Metrics) {
	t.Helper()
	gin.SetMode(gin.TestMode)

	m := New(opts)
	r := gin.New()
	r.Use(m.Middleware())
	r.GET(DefaultMetricsPath, gin.WrapH(m.Handler()))
	r.GET("/ping/:id", func(c *gin.Context) {
		c.String(200, "pong")
	})
	r.GET("/boom", func(c *gin.Context) {
		c.String(500, "boom")
	})
	return r, m
}

func TestMiddlewareRecordsRequestMetrics(t *testing.T) {
	r, _ := newTestRouter(t, Options{Application: "markdown-writer"})

	req := httptest.NewRequest("GET", "/ping/42", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	metricsReq := httptest.NewRequest("GET", DefaultMetricsPath, nil)
	metricsW := httptest.NewRecorder()
	r.ServeHTTP(metricsW, metricsReq)

	body := metricsW.Body.String()
	for _, want := range []string{
		`application="markdown-writer"`,
		`http_server_requests_seconds_count{`,
		`uri="/ping/:id"`,
		`status="200"`,
		`outcome="SUCCESS"`,
	} {
		if !strings.Contains(body, want) {
			t.Errorf("expected metrics output to contain %q, got:\n%s", want, body)
		}
	}
}

func TestMiddlewareServerErrorOutcome(t *testing.T) {
	r, _ := newTestRouter(t, Options{Application: "markdown-writer"})

	req := httptest.NewRequest("GET", "/boom", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	metricsReq := httptest.NewRequest("GET", DefaultMetricsPath, nil)
	metricsW := httptest.NewRecorder()
	r.ServeHTTP(metricsW, metricsReq)

	body := metricsW.Body.String()
	if !strings.Contains(body, `outcome="SERVER_ERROR"`) {
		t.Errorf("expected outcome=SERVER_ERROR in metrics output, got:\n%s", body)
	}
}

func TestMiddlewareSkipsMetricsPaths(t *testing.T) {
	r, _ := newTestRouter(t, Options{Application: "markdown-writer"})

	req := httptest.NewRequest("GET", DefaultMetricsPath, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	metricsReq := httptest.NewRequest("GET", DefaultMetricsPath, nil)
	metricsW := httptest.NewRecorder()
	r.ServeHTTP(metricsW, metricsReq)

	if strings.Contains(metricsW.Body.String(), `uri="/metrics"`) {
		t.Error("expected /metrics requests to be excluded from instrumentation")
	}
}

func TestOutcomeLabel(t *testing.T) {
	cases := map[int]string{
		100: "INFORMATIONAL",
		200: "SUCCESS",
		301: "REDIRECTION",
		404: "CLIENT_ERROR",
		500: "SERVER_ERROR",
	}
	for status, want := range cases {
		if got := outcomeLabel(status); got != want {
			t.Errorf("outcomeLabel(%d) = %q, want %q", status, got, want)
		}
	}
}

func TestEnvLabel(t *testing.T) {
	if got := envLabel(""); got != "unknown" {
		t.Errorf("envLabel(\"\") = %q, want \"unknown\"", got)
	}
	if got := envLabel(" production "); got != "production" {
		t.Errorf("envLabel(\" production \") = %q, want \"production\"", got)
	}
}
