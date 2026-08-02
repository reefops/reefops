package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/reefops/reefops/services/authorizer/internal/actorcontext"
	auditdb "github.com/reefops/reefops/services/authorizer/internal/audit/postgresql"
	"github.com/reefops/reefops/services/authorizer/internal/authorization"
	"github.com/reefops/reefops/services/authorizer/internal/observability"
	"github.com/reefops/reefops/services/authorizer/internal/openfga"
	"github.com/reefops/reefops/services/authorizer/internal/transport/grpcauthz"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	grpc_health_v1 "google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	observability.ConfigureLogger()
	if err := run(); err != nil {
		slog.Error("authorizer stopped", "error", err)
		os.Exit(1)
	}
}
func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	listenAddress := envDefault("AUTHORIZER_GRPC_LISTEN_ADDRESS", ":9002")
	managementAddress := envDefault("AUTHORIZER_MANAGEMENT_LISTEN_ADDRESS", ":9003")
	env, err := requiredEnvironment("OPENFGA_API_URL", "OPENFGA_STORE_ID", "OPENFGA_AUTHORIZATION_MODEL_ID", "OPENFGA_API_TOKEN", "OPENFGA_ENVIRONMENT_ID", "OPENFGA_MODEL_SHA256", "AUDIT_DATABASE_URL", "ACTOR_CONTEXT_PRIVATE_KEY_PKCS8_PEM_B64", "ACTOR_CONTEXT_PUBLIC_KEY_PEM_B64", "ACTOR_CONTEXT_ACTIVE_KID", "ACTOR_CONTEXT_ALGORITHM")
	if err != nil {
		return err
	}
	shutdownTracing, err := observability.SetupTracing(ctx, os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"), env["OPENFGA_ENVIRONMENT_ID"], envDefault("REEFOPS_VERSION", "development"))
	if err != nil {
		return err
	}
	defer func() {
		shutdownContext, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := shutdownTracing(shutdownContext); err != nil {
			slog.Error("shutdown tracing", "error", err)
		}
	}()
	openFGAClient, err := openfga.New(env["OPENFGA_API_URL"], env["OPENFGA_STORE_ID"], env["OPENFGA_AUTHORIZATION_MODEL_ID"], env["OPENFGA_API_TOKEN"], 2*time.Second)
	if err != nil {
		return err
	}
	signer, err := actorcontext.New(env["ACTOR_CONTEXT_PRIVATE_KEY_PKCS8_PEM_B64"], env["ACTOR_CONTEXT_PUBLIC_KEY_PEM_B64"], env["ACTOR_CONTEXT_ACTIVE_KID"], env["ACTOR_CONTEXT_ALGORITHM"], 20*time.Second)
	if err != nil {
		return err
	}
	poolConfig, err := pgxpool.ParseConfig(env["AUDIT_DATABASE_URL"])
	if err != nil {
		return fmt.Errorf("parse audit database URL: %w", err)
	}
	poolConfig.MaxConns = 8
	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return fmt.Errorf("open audit database: %w", err)
	}
	defer pool.Close()
	pingCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if err := pool.Ping(pingCtx); err != nil {
		return fmt.Errorf("connect to audit database: %w", err)
	}
	service, err := authorization.NewService(authorization.Config{EnvironmentID: env["OPENFGA_ENVIRONMENT_ID"], StoreID: env["OPENFGA_STORE_ID"], ModelID: env["OPENFGA_AUTHORIZATION_MODEL_ID"], ModelSHA256: env["OPENFGA_MODEL_SHA256"], ActorContextAudience: envDefault("ACTOR_CONTEXT_AUDIENCE", "reefops-services")}, openFGAClient, auditdb.NewRepository(pool), signer)
	if err != nil {
		return err
	}
	listener, err := net.Listen("tcp", listenAddress)
	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}
	metrics, registry := observability.NewMetrics()
	server := grpc.NewServer(grpc.StatsHandler(otelgrpc.NewServerHandler()))
	authv3.RegisterAuthorizationServer(server, grpcauthz.New(service, metrics))
	healthServer := health.NewServer()
	healthServer.SetServingStatus("envoy.service.auth.v3.Authorization", grpc_health_v1.HealthCheckResponse_SERVING)
	grpc_health_v1.RegisterHealthServer(server, healthServer)
	managementServer := &http.Server{Addr: managementAddress, Handler: observability.Handler(registry, func(readinessContext context.Context) error {
		if err := pool.Ping(readinessContext); err != nil {
			return fmt.Errorf("audit database: %w", err)
		}
		return openFGAClient.Ready(readinessContext)
	}), ReadHeaderTimeout: 2 * time.Second, IdleTimeout: 30 * time.Second}
	errCh := make(chan error, 1)
	go func() { errCh <- server.Serve(listener) }()
	go observability.ServeHTTP(managementServer, errCh)
	slog.Info("authorizer listening", "grpc_address", listenAddress, "management_address", managementAddress)
	select {
	case <-ctx.Done():
		healthServer.SetServingStatus("envoy.service.auth.v3.Authorization", grpc_health_v1.HealthCheckResponse_NOT_SERVING)
		server.GracefulStop()
		shutdownContext, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := managementServer.Shutdown(shutdownContext); err != nil {
			return fmt.Errorf("shutdown management HTTP: %w", err)
		}
		return nil
	case err := <-errCh:
		return fmt.Errorf("serve gRPC: %w", err)
	}
}
func requiredEnvironment(names ...string) (map[string]string, error) {
	values := make(map[string]string, len(names))
	for _, name := range names {
		value := os.Getenv(name)
		if value == "" {
			return nil, fmt.Errorf("%s is required", name)
		}
		values[name] = value
	}
	return values, nil
}
func envDefault(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}
