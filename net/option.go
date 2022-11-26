package net

import (
	"context"
	"os"

	"github.com/weiqiangxu/user/net/transport"
)

// Option is an application option.
type Option func(o *options)

type TracerFn func(o *options) error

// options is an application options.
type options struct {
	id      string
	name    string
	version string

	ctx  context.Context
	sigs []os.Signal

	servers []transport.Server

	attributes []KeyValue
	agentAddr  string
}

// ID with service id.
func ID(id string) Option {
	return func(o *options) { o.id = id }
}

// Name with service name.
func Name(name string) Option {
	return func(o *options) { o.name = name }
}

// Version with service version.
func Version(version string) Option {
	return func(o *options) { o.version = version }
}

// Context with service context.
func Context(ctx context.Context) Option {
	return func(o *options) { o.ctx = ctx }
}

// Server with transport servers.
func Server(srv ...transport.Server) Option {
	return func(o *options) { o.servers = srv }
}

// Signal with exit signals.
func Signal(sigs ...os.Signal) Option {
	return func(o *options) { o.sigs = sigs }
}

// Tracing with agent address and span attributes
func Tracing(agentAddr string, attributes ...KeyValue) Option {
	return func(o *options) {
		o.agentAddr = agentAddr
		o.attributes = attributes
	}
}
