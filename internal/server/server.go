package server

import (
	"context"
	"drpshop/internal/conf"
	"drpshop/pkg/token"
	"fmt"
	consul "github.com/go-kratos/consul/registry"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/middleware/metrics"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	prom "github.com/go-kratos/prometheus/metrics"
	"github.com/google/wire"
	consulAPI "github.com/hashicorp/consul/api"
	"github.com/prometheus/client_golang/prometheus"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(NewHttpMiddleware, NewGrpcMiddleware, NewHTTPServer, NewGRPCServer, NewRegistrar)

var (
	_metricSeconds = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "server",
		Subsystem: "requests",
		Name:      "duration_ms",
		Help:      "server requests duration(ms).",
		Buckets:   []float64{5, 10, 25, 50, 100, 250, 500, 1000},
	}, []string{"kind", "operation"})

	_metricRequests = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "client",
		Subsystem: "requests",
		Name:      "code_total",
		Help:      "The total number of processed requests",
	}, []string{"kind", "operation", "code", "reason"})
)

func getOperation(handler middleware.Handler) middleware.Handler {
	return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
		if tr, ok := transport.FromServerContext(ctx); ok {
			fmt.Println("==tr.Operation====", tr.Operation())
		}
		return handler(ctx, req)
	}
}

// 全局Middleware
func NewHttpMiddleware(c *conf.Data, logger log.Logger, tp *tracesdk.TracerProvider) (middleware.Middleware, error) {
	// 获取公钥实例
	jk, err := token.NewPublic(c.JwtCert)
	if err != nil {
		return nil, err
	}

	return middleware.Chain(
		recovery.Recovery(),
		tracing.Server(tracing.WithTracerProvider(tp)),
		logging.Server(logger),
		metrics.Server(
			metrics.WithSeconds(prom.NewHistogram(_metricSeconds)),
			metrics.WithRequests(prom.NewCounter(_metricRequests)),
		),
		metadata.Client(),
		// jwt tk验证为whole模式下使用，微服务模式下在网关鉴权
		//token.AuthTokenGinMiddleware(jk),
		selector.Server(token.AuthMiddleware(jk)).Prefix("/").Build(),
		validate.Validator(),
		getOperation,
	), nil
}

// 全局Middleware
func NewGrpcMiddleware(c *conf.Data, logger log.Logger, tp *tracesdk.TracerProvider) (middleware.Middleware, error) {
	jk, err := token.NewPublic(c.JwtCert)
	if err != nil {
		return nil, err
	}
	return middleware.Chain(
		recovery.Recovery(),
		tracing.Server(
			tracing.WithTracerProvider(tp)),
		logging.Server(logger),
		metrics.Server(
			metrics.WithSeconds(prom.NewHistogram(_metricSeconds)),
			metrics.WithRequests(prom.NewCounter(_metricRequests)),
		),
		metadata.Server(),
		validate.Validator(),
		getOperation,
		//selector.Server(token.AuthMiddleware(jk)).Build(),
		selector.Server(token.AuthMiddleware(jk)).Prefix("/").Build(),
	), nil
}

type App struct {
	hs *http.Server
	gs *grpc.Server
}

func NewApp(hs *http.Server, gs *grpc.Server) *App {
	return &App{
		hs: hs,
		gs: gs,
	}
}

func (a *App) Server() []transport.Server {
	//data := []transport.Server{a.hs, a.gs}
	data := make([]transport.Server, 0)
	if a.hs != nil {
		data = append(data, a.hs)
	}
	if a.gs != nil {
		data = append(data, a.gs)
	}
	return data
}

func NewRegistrar(conf *conf.Registry) registry.Registrar {
	c := consulAPI.DefaultConfig()
	c.Address = conf.Consul.Address
	c.Scheme = conf.Consul.Scheme
	cli, err := consulAPI.NewClient(c)
	if err != nil {
		panic(err)
	}
	r := consul.New(cli, consul.WithHealthCheck(false))
	return r
}
