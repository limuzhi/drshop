package service

import (
	"drpshop/internal/apps/admin/biz"
	"drpshop/internal/apps/admin/data"
	"drpshop/internal/apps/admin/middle"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerOperLogRouter)
}

func registerOperLogRouter(server *AdminService) {
	operLogRepo := data.NewOperLogRepo(server.DataData, server.Logger)
	operLogc := biz.NewOperLogController(operLogRepo, server.Logger)
	//操作日志记录
	for i := 0; i < 3; i++ {
		go operLogRepo.SaveOperlogChannel(middle.OperLogChan)
	}
	router := server.Router.Group("/api/v1").Use(middle.CasbinMiddleware(server.Casbinc))
	{
		//操作日志
		router.GET("/operLog/list", operLogc.OperLogList)
		router.GET("/operLog/detail", operLogc.OperLogDetail)
		router.DELETE("/operLog/batchDelete", operLogc.OperLogDelete)
		router.DELETE("/operLog/clear", operLogc.OperLogClear)
	}
}
