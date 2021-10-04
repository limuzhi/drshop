package service

import (
	"drpshop/internal/apps/admin/biz"
	"drpshop/internal/apps/admin/data"
	"drpshop/internal/apps/admin/middle"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerPostRouter)
}

func registerPostRouter(server *AdminService)  {
	postRepo := data.NewPostRepo(server.DataData, server.Logger)
	postc := biz.NewPostController(postRepo, server.Logger)
	router := server.Router.Group("/api/v1").Use(middle.CasbinMiddleware(server.Casbinc))
	{
		//岗位 TODO
		router.GET("/post/list", postc.PostList)
		router.GET("/post/detail", postc.PostDetail)
		router.POST("/post/create", postc.PostCreate)
		router.PUT("/post/update", postc.PostUpdate)
		router.DELETE("/post/delete", postc.PostDelete)
	}

}
