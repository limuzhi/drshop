package admin

import (
	"drpshop/internal/apps/admin/service"
	"drpshop/internal/server"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
)

type Admin struct {
	*server.App
}

func NewAdminApp(hs *http.Server, gs *grpc.Server, admin *service.AdminService) *Admin {

	router := admin.RegisterSysServiceServer()
	hs.HandlePrefix("/", router)

	return &Admin{App: server.NewApp(hs, gs)}
}
