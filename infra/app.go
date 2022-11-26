package infra

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/pkg/errors"
	"github.com/segmentio/ksuid"
	"github.com/weiqiangxu/common-config/logger"
	"golang.org/x/sync/errgroup"
)

// AppInfo is application context value.
type AppInfo interface {
	ID() string
	Name() string
	Version() string
}

// App is an application components lifecycle manager
type App struct {
	opts   options
	ctx    context.Context
	cancel func()
}

// New create an application lifecycle manager.
func New(opts ...Option) *App {
	options := options{
		ctx:  context.Background(),
		sigs: []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT},
	}
	options.id = ksuid.New().String()
	for _, o := range opts {
		o(&options)
	}
	ctx, cancel := context.WithCancel(options.ctx)
	return &App{
		ctx:    ctx,
		cancel: cancel,
		opts:   options,
	}
}

// ID returns app instance id.
func (a *App) ID() string { return a.opts.id }

// Name returns service name.
func (a *App) Name() string { return a.opts.name }

// Version returns app version.
func (a *App) Version() string { return a.opts.version }

// Run executes all OnStart hooks registered with the application's Lifecycle.
func (a *App) Run() error {
	ctx := NewContext(a.ctx, a)
	if a.opts.agentAddr != "" {
		tp, err := configAgent(ctx, a.opts.agentAddr, a.opts.name, a.opts.version, a.opts.attributes...)
		if err != nil {
			return err
		}
		defer func() {
			if err := tp.Shutdown(context.Background()); err != nil {
				logger.Errorf("Shutting down tracer provider (%v)", err)
			}
		}()
	}
	eg, ctx := errgroup.WithContext(ctx)
	wg := sync.WaitGroup{}
	for _, srv := range a.opts.servers {
		srv := srv
		eg.Go(func() error {
			<-ctx.Done() // wait for stop signal
			return srv.Stop(ctx)
		})
		wg.Add(1)
		eg.Go(func() error {
			wg.Done()
			return srv.Start(ctx)
		})
	}
	wg.Wait()
	c := make(chan os.Signal, 1)
	signal.Notify(c, a.opts.sigs...)
	eg.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-c:
				return a.Stop()
			}
		}
	})
	if err := eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}
	return nil
}

// Stop gracefully stops the application.
func (a *App) Stop() error {
	if a.cancel != nil {
		a.cancel()
	}
	return nil
}

type appKey struct{}

// NewContext returns a new Context that carries value.
func NewContext(ctx context.Context, s AppInfo) context.Context {
	return context.WithValue(ctx, appKey{}, s)
}

// FromContext returns the Transport value stored in ctx, if any.
func FromContext(ctx context.Context) (s AppInfo, ok bool) {
	s, ok = ctx.Value(appKey{}).(AppInfo)
	return
}
