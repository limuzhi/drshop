package data

import (
	"context"
	sysv1 "drpshop/api/sys/v1"
	"drpshop/internal/conf"
	"drpshop/pkg/token"
	consul "github.com/go-kratos/consul/registry"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
	consulAPI "github.com/hashicorp/consul/api"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(
	NewData, NewTokenPublic, NewDiscovery, NewRegistrar, NewSysServiceClient,
	NewUserRepo,
)

// Data .
type Data struct {
	conf *conf.Data
	log  *log.Helper
	rdb  *redis.Client
	Jk   *token.Token
	sc   sysv1.SysServiceClient
}

// NewData .
func NewData(conf *conf.Data, logger log.Logger, rdb *redis.Client, sc sysv1.SysServiceClient) (*Data, error) {
	jk, err := NewTokenPublic(conf.JwtCert)
	if err != nil {
		return nil, err
	}
	return &Data{
		log:  log.NewHelper(log.With(logger, "module", "data")),
		rdb:  rdb,
		sc:   sc,
		Jk:   jk,
		conf: conf,
	}, nil
}

func NewTokenPublic(jwtCert []byte) (*token.Token, error) {
	jk, err := token.NewPublic(jwtCert)
	if err != nil {
		return nil, err
	}
	return jk, nil
}

func NewDiscovery(conf *conf.Registry) registry.Discovery {
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

func NewSysServiceClient(r registry.Discovery, tp *tracesdk.TracerProvider) sysv1.SysServiceClient {
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint("discovery:///drshop.sys.service"),
		grpc.WithDiscovery(r),
		grpc.WithMiddleware(
			metadata.Client(),
			tracing.Client(tracing.WithTracerProvider(tp)),
		),
	)
	if err != nil {
		panic(err)
	}
	return sysv1.NewSysServiceClient(conn)
}
