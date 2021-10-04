package service

import (
	"drpshop/internal/apps/admin/biz"
	"drpshop/internal/apps/admin/data"
	"drpshop/internal/apps/admin/middle"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerLoginLogRouter)
}

func registerLoginLogRouter(server *AdminService) {
	loginLogRepo := data.NewLoginLogRepo(server.DataData, server.Logger)
	loginLogc := biz.NewLoginLogController(loginLogRepo, server.Logger)
	router := server.Router.Group("/api/v1").Use(middle.CasbinMiddleware(server.Casbinc))
	{
		//登录日志
		router.GET("/loginLog/list", loginLogc.LoginList)
		router.DELETE("/loginLog/batchDelete", loginLogc.BatchDeleteLog)
		router.DELETE("/loginLog/clear", loginLogc.Clearlog)
	}
}
