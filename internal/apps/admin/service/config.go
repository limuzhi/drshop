package service

import (
	"drpshop/internal/apps/admin/biz"
	"drpshop/internal/apps/admin/data"
	"drpshop/internal/apps/admin/middle"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerConfigRouter)
}

func registerConfigRouter(server *AdminService) {
	configRepo := data.NewConfigRepo(server.DataData, server.Logger)
	configc := biz.NewConfigController(configRepo, server.Logger)
	router := server.Router.Group("/api/v1").Use(middle.CasbinMiddleware(server.Casbinc))
	{
		//配置 TODO
		router.GET("/config/key", configc.GetConfigByKey)
		router.GET("/config/list", configc.ConfigList)
		router.GET("/config/detail", configc.ConfigDetail)
		router.POST("/config/create", configc.ConfigCreate)
		router.PUT("/config/update", configc.ConfigUpdate)
		router.DELETE("/config/delete", configc.ConfigDelete)
	}
}
