package observability

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestManagementHandlerExposesMetricsAndReadiness(t *testing.T) {
	metrics, registry := NewMetrics()
	metrics.Observe("deny", "openfga_denied", 10*time.Millisecond)
	handler := Handler(registry, func(context.Context) error { return nil })
	for _, path := range []string{"/livez", "/readyz", "/metrics"} {
		request := httptest.NewRequest(http.MethodGet, path, nil)
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, request)
		if response.Code != http.StatusOK {
			t.Fatalf("%s returned %d", path, response.Code)
		}
		if path == "/metrics" && !strings.Contains(response.Body.String(), "reefops_authorizer_checks_total") {
			t.Fatal("authorizer metric missing")
		}
	}
}

func TestReadinessFailsClosed(t *testing.T) {
	handler := Handler(mustRegistry(), func(context.Context) error { return context.DeadlineExceeded })
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/readyz", nil))
	if response.Code != http.StatusServiceUnavailable {
		t.Fatalf("readiness returned %d", response.Code)
	}
}
func mustRegistry() *prometheus.Registry { _, registry := NewMetrics(); return registry }
