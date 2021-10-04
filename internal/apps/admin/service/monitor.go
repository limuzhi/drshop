package service

import (
	"drpshop/internal/apps/admin/biz"
	"drpshop/internal/apps/admin/data"
	"drpshop/internal/apps/admin/middle"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerMonitorRouter)
}

func registerMonitorRouter(server *AdminService) {
	monitorRepo := data.NewMonitorRepo(server.DataData, server.Logger)
	monitorc := biz.NewMonitorController(monitorRepo, server.Logger)
	router := server.Router.Group("/api/v1").Use(middle.CasbinMiddleware(server.Casbinc))
	{
		router.GET("/monitor/server", monitorc.ServerInfo)
	}
}
