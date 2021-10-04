package service

import (
	"drpshop/internal/apps/admin/biz"
	"drpshop/internal/apps/admin/data"
	"github.com/gin-gonic/gin"
	"mime"
)

func init() {
	routerNoCheckRole = append(routerNoCheckRole, registerBaseRouter)
}

func registerBaseRouter(server *AdminService) {
	userRepo := data.NewUserRepo(server.DataData, server.Logger)
	userc := biz.NewUserController(userRepo, server.Logger)
	router := server.Router.Group("/api/v1")
	{
		// 登录登出刷新token无需鉴权
		router.POST("/login", userc.Login)
		router.GET("/captcha-img", userc.CaptchaImg)
		router.POST("/captcha", userc.Captcha)
	}
	r := server.Router.Group("")
	staticFileRouter(r)
}

func staticFileRouter(r *gin.RouterGroup) {
	err := mime.AddExtensionType(".js", "application/javascript")
	if err != nil {
		return
	}
	// 静态文件
	r.Static("/static", "./static")
	r.Static("/form-generator", "./static/form-generator")
}
