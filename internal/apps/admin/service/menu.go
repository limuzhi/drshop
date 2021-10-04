package service

import (
	"drpshop/internal/apps/admin/biz"
	"drpshop/internal/apps/admin/data"
	"drpshop/internal/apps/admin/middle"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerMenuRouter)
}

func registerMenuRouter(server *AdminService) {
	menuRepo := data.NewMenuRepo(server.DataData, server.Logger)
	menuc := biz.NewMenuController(menuRepo, server.Logger)
	router := server.Router.Group("/api/v1").Use(middle.CasbinMiddleware(server.Casbinc))
	{
		//菜单---
		//获取用户的可访问菜单树
		router.GET("/menu/tree", menuc.GetMenuTree)
		router.GET("/menu/list", menuc.GetMenuList)
		router.POST("/menu/create", menuc.MenuCreate)
		router.PUT("/menu/update", menuc.MenuUpdate)
		router.DELETE("/menu/batchDelete", menuc.MenubatchDelete)
		router.GET("/menu/access/list", menuc.GetUserMenuListByUserId)
		router.GET("/menu/access/tree", menuc.GetUserMenuTreeByUserId)
	}
}
