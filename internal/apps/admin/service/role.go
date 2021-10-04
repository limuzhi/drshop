package service

import (
	"drpshop/internal/apps/admin/biz"
	"drpshop/internal/apps/admin/data"
	"drpshop/internal/apps/admin/middle"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerRoleRouter)
}

func registerRoleRouter(server *AdminService) {
	roleRepo := data.NewRoleRepo(server.DataData, server.Logger)
	rolec := biz.NewRoleController(roleRepo, server.Logger)
	router := server.Router.Group("/api/v1").Use(middle.CasbinMiddleware(server.Casbinc))
	{
		//角色
		//角色列表
		router.GET("/role/list", rolec.RoleList)
		//创建角色
		router.POST("/role/create", rolec.RoleCreate)
		//更新角色
		router.PUT("/role/update", rolec.RoleUpdate)
		//删除角色
		router.DELETE("/role/batchDelete", rolec.BatchDelete)
		//修改角色状态
		router.PUT("/role/changeStatus", rolec.ChangeStatus)
		//根据角色ID获取角色的权限菜单
		router.GET("/role/menus", rolec.RoleMenusList)
		//更新角色的权限菜单
		router.PUT("/role/menus/update", rolec.RoleMenusUpdate)

		//根据角色ID获取角色的权限接口
		router.GET("/role/apis", rolec.RoleApisList)
		//更新角色的权限接口
		router.PUT("/role/apis/update", rolec.RoleApisUpdate)
	}
}
