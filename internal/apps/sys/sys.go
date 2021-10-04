package sys

import (
	v1 "drpshop/api/sys/v1"
	"drpshop/internal/server"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
)

type Sys struct {
	*server.App
}

func NewSysApp(hs *http.Server, gs *grpc.Server, m middleware.Middleware,
	sys v1.SysServiceServer,
) *Sys {
	// grpc
	v1.RegisterSysServiceServer(gs, sys)
	return &Sys{App: server.NewApp(hs, gs)}
}
