package grpc

import (
	"context"
	"fmt"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	grpc_insecure "google.golang.org/grpc/credentials/insecure"
)

// ClientOption is gRPC client option.
type ClientOption func(o *clientOptions)

// WithEndpoint with client endpoint.
func WithEndpoint(endpoint string) ClientOption {
	return func(o *clientOptions) {
		o.endpoint = endpoint
	}
}

// WithUnaryInterceptor returns a DialOption that specifies the interceptor for unary RPCs.
func WithUnaryInterceptor(in ...grpc.UnaryClientInterceptor) ClientOption {
	return func(o *clientOptions) {
		o.ints = in
	}
}

// WithOptions with gRPC options.
func WithOptions(opts ...grpc.DialOption) ClientOption {
	return func(o *clientOptions) {
		o.grpcOpts = opts
	}
}

func WithTracing(opt bool) ClientOption {
	return func(o *clientOptions) {
		o.tracing = opt
	}
}

// clientOptions is gRPC Client
type clientOptions struct {
	endpoint string
	ints     []grpc.UnaryClientInterceptor
	grpcOpts []grpc.DialOption

	tracing bool
}

// Dial returns a GRPC connection.
func Dial(ctx context.Context, opts ...ClientOption) (*grpc.ClientConn, error) {
	return dial(ctx, false, opts...)
}

// DialInsecure returns an insecure GRPC connection.
func DialInsecure(ctx context.Context, opts ...ClientOption) (*grpc.ClientConn, error) {
	return dial(ctx, true, opts...)
}

func dial(ctx context.Context, insecure bool, opts ...ClientOption) (*grpc.ClientConn, error) {
	options := clientOptions{}
	for _, o := range opts {
		o(&options)
	}

	unaryInterceptors := []grpc.UnaryClientInterceptor{}
	streamInterceptors := []grpc.StreamClientInterceptor{}

	if len(options.ints) > 0 {
		unaryInterceptors = append(unaryInterceptors, options.ints...)
	}

	if options.tracing {
		unaryInterceptors = append(unaryInterceptors, otelgrpc.UnaryClientInterceptor())
		streamInterceptors = append(streamInterceptors, otelgrpc.StreamClientInterceptor())
	}

	grpcOpts := []grpc.DialOption{
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": %q}`, roundrobin.Name)),
		grpc.WithChainUnaryInterceptor(unaryInterceptors...),
		grpc.WithChainStreamInterceptor(streamInterceptors...),
	}

	if insecure {
		// grpc.WithInsecure()
		grpcOpts = append(grpcOpts, grpc.WithTransportCredentials(grpc_insecure.NewCredentials()))
	}
	if len(options.grpcOpts) > 0 {
		grpcOpts = append(grpcOpts, options.grpcOpts...)
	}

	return grpc.DialContext(ctx, options.endpoint, grpcOpts...)
}
