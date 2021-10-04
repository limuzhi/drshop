package service

import (
	"drpshop/internal/apps/admin/biz"
	"drpshop/internal/apps/admin/data"
	"drpshop/internal/apps/admin/middle"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerUserRouter)
}

func registerUserRouter(server *AdminService) {
	userRepo := data.NewUserRepo(server.DataData, server.Logger)
	userc := biz.NewUserController(userRepo, server.Logger)

	router := server.Router.Group("/api/v1").Use(middle.CasbinMiddleware(server.Casbinc))
	{
		router.POST("/logout", userc.Logout)
		//用户---
		router.GET("/user/info", userc.GetUserDetail)
		router.GET("/user/list", userc.UserList)
		router.PUT("/user/changeStatus", userc.ChangeStatus)
		router.GET("/user/postroleList", userc.GetPostAndRoleList)
		router.DELETE("/user/batchDelete", userc.BatchDeleteUser)
		router.POST("/user/create", userc.UserCreate)
		router.PUT("/user/update", userc.UpdateUser)
		router.PUT("/user/changePwd", userc.ChangePwd)
		router.GET("/user/profile", userc.UserPofile)
		router.PUT("/user/profileSet", userc.UserProfileSet)
		router.POST("/user/avatar", userc.UploadAvatar)
	}
}
