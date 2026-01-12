package platform

import (
	"context"
	"errors"
	"log"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/propagation"
	otellog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
)

// setupOTelSDK bootstraps the OpenTelemetry pipeline.
// If it does not return an error, make sure to call shutdown for proper cleanup.
func SetupOTelSDK(ctx context.Context, serviceName, logURI, metricURI, traceURI string) (func(context.Context) error, error) {
	var shutdownFuncs []func(context.Context) error
	var err error

	// shutdown calls cleanup functions registered via shutdownFuncs.
	// The errors from the calls are joined.
	// Each registered cleanup will be invoked once.
	shutdown := func(ctx context.Context) error {
		var err error
		for _, fn := range shutdownFuncs {
			err = errors.Join(err, fn(ctx))
		}
		shutdownFuncs = nil
		return err
	}

	res, _ := resource.Merge(resource.Default(), resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(serviceName),
	))

	// handleErr calls shutdown for cleanup and makes sure that all errors are returned.
	handleErr := func(inErr error) {
		err = errors.Join(inErr, shutdown(ctx))
	}

	// Set up propagator.
	prop := newPropagator()
	otel.SetTextMapPropagator(prop)

	// 1. Tracer Provider
	if traceURI != "" {
		log.Printf("[OTel] Configurando Trace Exporter gRPC no endpoint: %s", traceURI)
		tp, err := newTracerProvider(ctx, traceURI, res)
		if err != nil {
			handleErr(err)
			return shutdown, err
		}
		shutdownFuncs = append(shutdownFuncs, tp.Shutdown)
		otel.SetTracerProvider(tp)
	} else {
		log.Println("[OTel] Trace Exporter desabilitado (URI vazia)")
	}

	// 2. Meter Provider
	if metricURI != "" {
		log.Printf("[OTel] Configurando Metric Exporter gRPC no endpoint: %s", metricURI)
		mp, err := newMeterProvider(ctx, metricURI, res)
		if err != nil {
			handleErr(err)
			return shutdown, err
		}
		shutdownFuncs = append(shutdownFuncs, mp.Shutdown)
		otel.SetMeterProvider(mp)
	} else {
		log.Println("[OTel] Metric Exporter desabilitado (URI vazia)")
	}

	// 3. Logger Provider
	if logURI != "" {
		log.Printf("[OTel] Configurando Log Exporter gRPC no endpoint: %s (Stdout também ativo)", logURI)
	} else {
		log.Println("[OTel] Log Exporter remoto desabilitado. Usando apenas Stdout.")
	}
	loggerProvider, err := newLoggerProvider(ctx, logURI, res)
	if err != nil {
		handleErr(err)
		return shutdown, err
	}

	shutdownFuncs = append(shutdownFuncs, loggerProvider.Shutdown)
	global.SetLoggerProvider(loggerProvider)

	return shutdown, err
}

func newPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}

func newTracerProvider(ctx context.Context, uri string, res *resource.Resource) (*trace.TracerProvider, error) {
	exporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithEndpoint(uri), otlptracegrpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(res),
	), nil
}

func newMeterProvider(ctx context.Context, uri string, res *resource.Resource) (*metric.MeterProvider, error) {
	exporter, err := otlpmetricgrpc.New(ctx, otlpmetricgrpc.WithEndpoint(uri), otlpmetricgrpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(metric.NewPeriodicReader(exporter, metric.WithInterval(3*time.Second))),
	), nil
}

func newLoggerProvider(ctx context.Context, uri string, res *resource.Resource) (*otellog.LoggerProvider, error) {
	var processors []otellog.Processor

	// Processador gRPC (Opcional)
	if uri != "" {
		remoteExporter, err := otlploggrpc.New(ctx, otlploggrpc.WithEndpoint(uri), otlploggrpc.WithInsecure())
		if err != nil {
			return nil, err
		}
		processors = append(processors, otellog.NewBatchProcessor(remoteExporter))
	}

	// Processador Stdout (Sempre Ativo)
	consoleExporter, err := stdoutlog.New(stdoutlog.WithPrettyPrint())
	if err != nil {
		return nil, err
	}
	processors = append(processors, otellog.NewBatchProcessor(consoleExporter))

	// Criamos o provider com todos os processadores coletados
	options := []otellog.LoggerProviderOption{otellog.WithResource(res)}
	for _, p := range processors {
		options = append(options, otellog.WithProcessor(p))
	}

	return otellog.NewLoggerProvider(options...), nil
}
