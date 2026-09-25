package httpmetrics

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewRequiresApplication(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("expected New to panic when Application is empty")
		}
	}()
	New(Options{})
}

func newTestMux(t *testing.T, opts Options) (*http.ServeMux, *Metrics) {
	t.Helper()

	m := New(opts)
	mux := http.NewServeMux()
	mux.Handle("GET /metrics", m.Handler())
	mux.HandleFunc("GET /items/{id}", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	mux.HandleFunc("POST /boom", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "boom", http.StatusInternalServerError)
	})
	return mux, m
}

func TestMiddlewareRecordsRequestMetrics(t *testing.T) {
	mux, m := newTestMux(t, Options{Application: "exia-artifact-promoter"})
	handler := m.Middleware(mux)(mux)

	req := httptest.NewRequest(http.MethodGet, "/items/42", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusNoContent)
	}

	metricsReq := httptest.NewRequest(http.MethodGet, DefaultMetricsPath, nil)
	metricsW := httptest.NewRecorder()
	handler.ServeHTTP(metricsW, metricsReq)

	body := metricsW.Body.String()
	for _, want := range []string{
		`application="exia-artifact-promoter"`,
		`http_server_requests_seconds_count{`,
		`uri="/items/{id}"`,
		`status="204"`,
		`outcome="SUCCESS"`,
	} {
		if !strings.Contains(body, want) {
			t.Errorf("expected metrics output to contain %q, got:\n%s", want, body)
		}
	}
	if strings.Contains(body, `uri="/metrics"`) {
		t.Error("expected /metrics requests to be excluded from instrumentation")
	}
}

func TestMiddlewareServerErrorOutcome(t *testing.T) {
	mux, m := newTestMux(t, Options{Application: "exia-artifact-promoter"})
	handler := m.Middleware(mux)(mux)

	req := httptest.NewRequest(http.MethodPost, "/boom", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	metricsReq := httptest.NewRequest(http.MethodGet, DefaultMetricsPath, nil)
	metricsW := httptest.NewRecorder()
	handler.ServeHTTP(metricsW, metricsReq)

	body := metricsW.Body.String()
	if !strings.Contains(body, `outcome="SERVER_ERROR"`) {
		t.Errorf("expected outcome=SERVER_ERROR in metrics output, got:\n%s", body)
	}
}

func TestMiddlewareUnmatchedRouteIsUnknown(t *testing.T) {
	mux, m := newTestMux(t, Options{Application: "exia-artifact-promoter"})
	handler := m.Middleware(mux)(mux)

	req := httptest.NewRequest(http.MethodGet, "/nope", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusNotFound)
	}

	metricsReq := httptest.NewRequest(http.MethodGet, DefaultMetricsPath, nil)
	metricsW := httptest.NewRecorder()
	handler.ServeHTTP(metricsW, metricsReq)

	body := metricsW.Body.String()
	if !strings.Contains(body, `uri="NOT_FOUND"`) {
		t.Errorf("expected uri=NOT_FOUND for unmatched route, got:\n%s", body)
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
