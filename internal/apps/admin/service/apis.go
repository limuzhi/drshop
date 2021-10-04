package service

import (
	"drpshop/internal/apps/admin/biz"
	"drpshop/internal/apps/admin/data"
	"drpshop/internal/apps/admin/middle"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerApisRouter)
}

func registerApisRouter(server *AdminService) {
	apisRepo := data.NewApisRepo(server.DataData, server.Logger)
	apisc := biz.NewApisController(apisRepo, server.Logger)
	router := server.Router.Group("/api/v1").Use(middle.CasbinMiddleware(server.Casbinc))
	{
		//接口
		//获取接口树
		router.GET("/api/tree", apisc.GetApiTree)
		//获取接口列表
		router.GET("/api/list", apisc.GetApiList)
		//创建接口
		router.POST("/api/create", apisc.ApiCreate)
		//更新接口
		router.PUT("/api/update", apisc.ApiUpdate)
		//批量删除接口
		router.DELETE("/api/batchDelete", apisc.ApiBatchDelete)
	}
}
