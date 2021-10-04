package service

import (
	"drpshop/internal/apps/admin/biz"
	"drpshop/internal/apps/admin/data"
	"drpshop/internal/apps/admin/middle"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerDeptRouter)
}

func registerDeptRouter(server *AdminService)  {
	deptRepo := data.NewDeptRepo(server.DataData, server.Logger)
	deptc := biz.NewDeptController(deptRepo, server.Logger)
	router := server.Router.Group("/api/v1").Use(middle.CasbinMiddleware(server.Casbinc))
	{
		//部门 TODO
		router.GET("/dept/list", deptc.DeptList)
		router.GET("/dept/tree", deptc.DeptTree)
		router.GET("/dept/detail", deptc.DeptDetail)
		router.POST("/dept/create", deptc.DeptCreate)
		router.PUT("/dept/update", deptc.DeptUpdate)
		router.DELETE("/dept/delete", deptc.DeptDelete)
	}

}
