package service

import (
	"drpshop/internal/apps/admin/biz"
	"drpshop/internal/apps/admin/data"
	"drpshop/internal/apps/admin/middle"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerDictRouter)
}

func registerDictRouter(server *AdminService) {
	dictDataRepo := data.NewDictDataRepo(server.DataData, server.Logger)
	dictDatac := biz.NewDictDataController(dictDataRepo, server.Logger)

	dictTypeRepo := data.NewDictTypeRepo(server.DataData, server.Logger)
	dictTypec := biz.NewDictTypeController(dictTypeRepo, server.Logger)
	router := server.Router.Group("/api/v1").Use(middle.CasbinMiddleware(server.Casbinc))
	{
		//字典数据 TODO ---
		//dictType获取
		router.GET("/dict-data/select", dictDatac.DictDataSelect)

		router.GET("/dict-data/list", dictDatac.DictDataList)
		//dictCode获取
		router.GET("/dict-data/detail", dictDatac.DictDataDetail)
		router.POST("/dict-data/create", dictDatac.DictDataCreate)
		router.PUT("/dict-data/update", dictDatac.DictDataUpdate)
		router.DELETE("/dict-data/delete", dictDatac.DictDataDelete)


		//字典分类数据
		router.GET("/dict-type/select", dictTypec.DictTypeSelect)
		router.GET("/dict-type/detail", dictTypec.DictTypeDetail)
		router.GET("/dict-type/list", dictTypec.DictTypeList)
		router.POST("/dict-type/create", dictTypec.DictTypeCreate)
		router.PUT("/dict-type/update", dictTypec.DictTypeUpdate)
		router.DELETE("/dict-type/delete", dictTypec.DictTypeDelete)

	}

}
