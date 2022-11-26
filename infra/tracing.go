package infra

import (
	"context"
	"fmt"
	"time"

	"github.com/weiqiangxu/common-config/logger"
	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdkTrace "go.opentelemetry.io/otel/sdk/trace"
	semConv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

var Tracer trace.Tracer

// KeyValue holds a key and value pair.
type KeyValue struct {
	Key   string
	Value string
}

type configAgentErrorHandler struct{}

func (e *configAgentErrorHandler) Handle(err error) {
	logger.Errorf("[tracing] %v", err)
}

func configAgent(ctx context.Context, agentAddr, service, version string, attributes ...KeyValue) (*sdkTrace.TracerProvider, error) {
	otel.SetErrorHandler(&configAgentErrorHandler{})

	expOptions := []otlptracegrpc.Option{
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(agentAddr),
	}

	grpcConnectionTimeout := 3 * time.Second
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, grpcConnectionTimeout)
	defer cancel()

	traceExp, err := otlptracegrpc.New(ctx, expOptions...)
	if err != nil {
		return nil, fmt.Errorf("failed to create the collector trace exporter (%s)", err)
	}

	attrs := []attribute.KeyValue{
		semConv.ServiceNameKey.String(service),
		semConv.ServiceVersionKey.String(version),
	}
	for _, attr := range attributes {
		if attr.Key != "" && attr.Value != "" {
			attrs = append(attrs, attribute.String(attr.Key, attr.Value))
		}
	}
	res, err := resource.New(ctx,
		resource.WithAttributes(attrs...),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource (%s)", err)
	}

	tp := sdkTrace.NewTracerProvider(
		sdkTrace.WithSampler(sdkTrace.ParentBased(sdkTrace.TraceIDRatioBased(1.0))),
		sdkTrace.WithBatcher(traceExp,
			sdkTrace.WithBatchTimeout(5*time.Second),
			sdkTrace.WithMaxExportBatchSize(10)),
		sdkTrace.WithResource(res),
	)
	otel.SetTracerProvider(tp)

	propagator := b3.New(b3.WithInjectEncoding(b3.B3MultipleHeader))
	otel.SetTextMapPropagator(propagator)

	Tracer = tp.Tracer("gyu")

	return tp, nil
}

func TraceID(ctx context.Context) string {
	if span := trace.SpanContextFromContext(ctx); span.HasTraceID() {
		return span.TraceID().String()
	}
	return ""
}

func SpanID(ctx context.Context) string {
	if span := trace.SpanContextFromContext(ctx); span.HasSpanID() {
		return span.SpanID().String()
	}
	return ""
}

func Span(ctx context.Context) trace.Span {
	if span := trace.SpanFromContext(ctx); span.SpanContext().IsValid() {
		return span
	}
	return nil
}
