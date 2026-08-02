package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/reefops/reefops/services/authorizer/internal/actorcontext"
	auditdb "github.com/reefops/reefops/services/authorizer/internal/audit/postgresql"
	"github.com/reefops/reefops/services/authorizer/internal/authorization"
	"github.com/reefops/reefops/services/authorizer/internal/openfga"
	"github.com/reefops/reefops/services/authorizer/internal/transport/grpcauthz"
	"google.golang.org/grpc"
)

func main() {
	if err := run(); err != nil {
		slog.Error("authorizer stopped", "error", err)
		os.Exit(1)
	}
}
func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	listenAddress := envDefault("AUTHORIZER_GRPC_LISTEN_ADDRESS", ":9002")
	env, err := requiredEnvironment("OPENFGA_API_URL", "OPENFGA_STORE_ID", "OPENFGA_AUTHORIZATION_MODEL_ID", "OPENFGA_API_TOKEN", "OPENFGA_ENVIRONMENT_ID", "OPENFGA_MODEL_SHA256", "AUDIT_DATABASE_URL", "ACTOR_CONTEXT_PRIVATE_KEY_PKCS8_PEM_B64", "ACTOR_CONTEXT_PUBLIC_KEY_PEM_B64", "ACTOR_CONTEXT_ACTIVE_KID", "ACTOR_CONTEXT_ALGORITHM")
	if err != nil {
		return err
	}
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
	server := grpc.NewServer()
	authv3.RegisterAuthorizationServer(server, grpcauthz.New(service))
	errCh := make(chan error, 1)
	go func() { errCh <- server.Serve(listener) }()
	slog.Info("authorizer listening", "address", listenAddress)
	select {
	case <-ctx.Done():
		server.GracefulStop()
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
