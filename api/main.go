package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/emersonmatsumoto/clean-go/api/internal/handlers"
	"github.com/emersonmatsumoto/clean-go/api/internal/platform"
	"github.com/emersonmatsumoto/clean-go/orders"
	"github.com/emersonmatsumoto/clean-go/payments"
	"github.com/emersonmatsumoto/clean-go/products"
	"github.com/emersonmatsumoto/clean-go/users"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type Config struct {
	MongoURI    string
	StripeKey   string
	LogURI      string
	MetricURI   string
	TraceURI    string
	ServiceName string
	Port        string
}

func main() {
	cfg := Config{
		MongoURI:    os.Getenv("MONGO_URI"),
		StripeKey:   os.Getenv("STRIPE_KEY"),
		LogURI:      os.Getenv("LOG_URI"),
		MetricURI:   os.Getenv("METRIC_URI"),
		TraceURI:    os.Getenv("TRACE_URI"),
		ServiceName: os.Getenv("SERVICE_NAME"),
		Port:        ":8080",
	}
	if cfg.MongoURI == "" {
		log.Fatal("MONGO_URI não foi configurada nas variáveis de ambiente")
	}

	if cfg.StripeKey == "" {
		log.Fatal("STRIPE_KEY não foi configurada nas variáveis de ambiente")
	}

	if err := run(cfg); err != nil {
		log.Fatalf("Erro fatal na aplicação: %v", err)
	}
}

func run(cfg Config) (err error) {
	// Handle SIGINT (CTRL+C) gracefully.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// Set up OpenTelemetry.
	otelShutdown, err := platform.SetupOTelSDK(ctx, cfg.ServiceName, cfg.LogURI, cfg.MetricURI, cfg.TraceURI)
	if err != nil {
		return err
	}
	// Handle shutdown properly so nothing leaks.
	defer func() {
		shCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err = errors.Join(err, otelShutdown(shCtx))
	}()

	mongoClient, mongoCleanup, err := platform.NewMongoClient(cfg.MongoURI)
	if err != nil {
		return err
	}
	defer mongoCleanup()

	userComp := users.NewComponent(mongoClient)
	productComp := products.NewComponent(mongoClient)
	paymentComp := payments.NewComponent(cfg.StripeKey)
	orderComp := orders.NewComponent(mongoClient, productComp, paymentComp, userComp)
	orderCtrl := handlers.NewOrderController(orderComp)

	mux := http.NewServeMux()
	handlers.RegisterRoutes(mux, orderCtrl)
	handler := otelhttp.NewHandler(mux, "api-server")
	srv := platform.NewServer(cfg.Port, handler)
	srvErr := make(chan error, 1)
	go func() {
		slog.Info("Iniciando servidor", "port", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			srvErr <- err
		}
	}()

	select {
	case err := <-srvErr:
		return fmt.Errorf("erro no servidor: %w", err)
	case <-ctx.Done():
		slog.Info("Sinal recebido, encerrando servidor...")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return srv.Shutdown(shutdownCtx)
}
