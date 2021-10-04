package vo

type DictTypeListReq struct {
	DictName string `json:"dictName" form:"dictName"`
	DictType string `json:"dictType" form:"dictType"`
	Status   string `json:"status" form:"status"`
	Pagination
}

type DictTypeDetailReq struct {
	DictId int64 `form:"dictId" json:"dictId"  validate:"required,gte=1"`
}

type DictTypeCreateReq struct {
	DictName string `json:"dictName" form:"dictName" validate:"required,min=1,max=100" comment:"字典名称"`
	DictType string `json:"dictType" form:"dictType" validate:"required,min=1,max=100" comment:"字典类型"`
	Status   string `json:"status" form:"status" validate:"oneof=1 2" comment:"状态（2正常 1停用）"`
	Remark   string `json:"remark" form:"remark" comment:"备注"`
	CreateBy int64  `json:"-" form:"-"`
}

type DictTypeUpdateReq struct {
	DictTypeCreateReq
	DictTypeDetailReq
	UpdateBy int64  `json:"-" form:"-"`
}

type DictTypeDeleteReq struct {
	DictIds []int64 `form:"dictIds" json:"dictIds"  validate:"required"`
}
