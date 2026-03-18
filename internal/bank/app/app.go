package app

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	slogotel "github.com/remychantenay/slog-otel"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"

	bankapi "github.com/romangurevitch/go-training/internal/bank/api"
	"github.com/romangurevitch/go-training/internal/bank/config"
	"github.com/romangurevitch/go-training/internal/bank/repository/postgres"
	"github.com/romangurevitch/go-training/internal/bank/service"
)

// Run initializes and starts the application.
// It returns an error if any part of the initialization or execution fails.
func Run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg := config.Values

	// 1. Infrastructure: Logging & Tracing
	logger := SetupLogger(cfg.ServiceName)
	slog.SetDefault(logger)

	tp, err := InitTracer(ctx, cfg)
	if err != nil {
		slog.ErrorContext(ctx, "failed to init tracer", slog.Any("error", err))
		return err
	}
	defer ShutdownTracer(tp)

	// 2. Data Access
	db, err := InitDB(ctx, cfg.DatabaseURL)
	if err != nil {
		return err // InitDB already logged it
	}
	defer func() {
		_ = db.Close()
	}()

	// 3. Application Logic & API Wiring
	srv, err := WireServer(db, logger, cfg)
	if err != nil {
		slog.ErrorContext(ctx, "failed to wire server", slog.Any("error", err))
		return err
	}

	// 4. Lifecycle Management
	return Serve(ctx, srv)
}

// SetupLogger configures and returns a structured logger.
func SetupLogger(serviceName string) *slog.Logger {
	jsonHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	otelBridge := otelslog.NewHandler(serviceName)
	otelEnricher := slogotel.OtelHandler{Next: slog.NewMultiHandler(jsonHandler, otelBridge)}
	return slog.New(otelEnricher)
}

// InitTracer initializes the OpenTelemetry tracer provider.
func InitTracer(ctx context.Context, cfg config.Config) (*sdktrace.TracerProvider, error) {
	exporter, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint(cfg.OTelEndpoint),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(cfg.ServiceName),
			semconv.ServiceVersion("0.1.0"),
		),
	)
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)
	return tp, nil
}

// ShutdownTracer handles graceful shutdown of the tracer provider.
func ShutdownTracer(tp *sdktrace.TracerProvider) {
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := tp.Shutdown(shutdownCtx); err != nil {
		slog.Error("tracer shutdown failed", slog.Any("error", err))
	}
}

// InitDB initializes the database connection and verifies health.
func InitDB(ctx context.Context, url string) (*sql.DB, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		slog.ErrorContext(ctx, "could not open db", slog.Any("error", err))
		return nil, err
	}

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		slog.ErrorContext(ctx, "could not ping db", slog.Any("error", err))
		return nil, err
	}
	return db, nil
}

// WireServer assembles the HTTP server with all its dependencies.
func WireServer(db *sql.DB, logger *slog.Logger, cfg config.Config) (*http.Server, error) {
	repo := postgres.New(db)
	svc := service.NewBankService(repo)

	apiCfg := bankapi.Config{
		JWTSecret:   cfg.JWTSecret,
		ServiceName: cfg.ServiceName,
	}

	router := bankapi.NewServer(svc, logger, apiCfg)

	return &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}, nil
}

// Serve starts the server and waits for the context to be canceled for graceful shutdown.
func Serve(ctx context.Context, srv *http.Server) error {
	errChan := make(chan error, 1)
	go func() {
		slog.Info("server starting", slog.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	select {
	case err := <-errChan:
		slog.ErrorContext(ctx, "http server error", slog.Any("error", err))
		return err
	case <-ctx.Done():
		slog.Info("shutting down server")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.ErrorContext(ctx, "server forced to shutdown", slog.Any("error", err))
		return err
	}

	slog.Info("server exited")
	return nil
}
