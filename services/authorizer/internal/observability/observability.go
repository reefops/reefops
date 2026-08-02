package observability

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	"google.golang.org/grpc/credentials/insecure"
)

type Metrics struct {
	checks   *prometheus.CounterVec
	duration *prometheus.HistogramVec
}

func NewMetrics() (*Metrics, *prometheus.Registry) {
	registry := prometheus.NewRegistry()
	registry.MustRegister(prometheus.NewGoCollector(), prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}))
	m := &Metrics{
		checks:   prometheus.NewCounterVec(prometheus.CounterOpts{Name: "reefops_authorizer_checks_total", Help: "Authorization checks by closed result and reason code."}, []string{"result", "reason"}),
		duration: prometheus.NewHistogramVec(prometheus.HistogramOpts{Name: "reefops_authorizer_check_duration_seconds", Help: "End-to-end ext_authz check latency.", Buckets: prometheus.DefBuckets}, []string{"result"}),
	}
	registry.MustRegister(m.checks, m.duration)
	return m, registry
}

func (m *Metrics) Observe(result, reason string, elapsed time.Duration) {
	m.checks.WithLabelValues(result, reason).Inc()
	m.duration.WithLabelValues(result).Observe(elapsed.Seconds())
}

type Readiness func(context.Context) error

func Handler(registry *prometheus.Registry, ready Readiness) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	mux.HandleFunc("/livez", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok\n"))
	})
	mux.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), time.Second)
		defer cancel()
		if err := ready(ctx); err != nil {
			http.Error(w, "not ready", http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok\n"))
	})
	return mux
}

func ConfigureLogger() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})))
}

func SetupTracing(ctx context.Context, endpoint, environment, version string) (func(context.Context) error, error) {
	if endpoint == "" {
		return func(context.Context) error { return nil }, nil
	}
	exporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithEndpoint(endpoint), otlptracegrpc.WithTLSCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("create OTLP trace exporter: %w", err)
	}
	res, err := resource.New(ctx, resource.WithAttributes(semconv.ServiceName("reefops-authorizer"), semconv.ServiceVersion(version), semconv.DeploymentEnvironmentName(environment), attribute.String("reefops.environment", environment)))
	if err != nil {
		return nil, fmt.Errorf("create trace resource: %w", err)
	}
	provider := tracesdk.NewTracerProvider(tracesdk.WithBatcher(exporter), tracesdk.WithResource(res))
	otel.SetTracerProvider(provider)
	return provider.Shutdown, nil
}

func ServeHTTP(server *http.Server, errorsChannel chan<- error) {
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		errorsChannel <- fmt.Errorf("serve management HTTP: %w", err)
	}
}
