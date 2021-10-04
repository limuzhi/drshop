package vo

type DictDataListReq struct {
	DictLabel string `json:"dictLabel" form:"dictLabel"`
	DictValue string `json:"dictValue" form:"dictValue"`
	DictType  string `json:"dictType" form:"dictType"`
	Status    string `json:"status" form:"status"`
	Pagination
}

type DictDataSelectReq struct {
	DictType string `json:"dictType" form:"dictType" validate:"required"`
}

type DictDataDetailReq struct {
	DictCode int64 `form:"dictCode" json:"dictCode"  validate:"required"`
}

type DictDataCreateReq struct {
	DictLabel string `json:"dictLabel" form:"dictLabel" validate:"required,min=1,max=128" comment:"字典标签"`
	DictValue string `json:"dictValue" form:"dictValue" validate:"required,min=1,max=128" comment:"字典键值"`
	DictType  string `json:"dictType" form:"dictType" validate:"required,min=1,max=128" comment:"字典类型"`
	Status    string `json:"status" form:"status" validate:"oneof=1 2" comment:"状态（2正常 1停用）"`
	DictSort  int    `json:"dictSort" form:"dictSort" validate:"gte=1,lte=999" comment:"字典排序"`
	IsDefault int    `json:"isDefault" form:"isDefault" validate:"oneof=1 2" comment:"是否默认（1是 2否）"`
	CssClass  string `json:"cssClass" form:"cssClass" comment:"样式属性（其他样式扩展）"`
	Remark    string `json:"remark" form:"remark" comment:"备注"`
	ListClass string `json:"listClass" form:"listClass" comment:"表格回显样式"`
	CreateBy  int64  `json:"-" form:"-"`
}

type DictDataUpdateReq struct {
	DictDataCreateReq
	DictCode int64 `json:"dictCode" form:"dictCode" validate:"required,gte=1" comment:"编号"`
	UpdateBy int64 `json:"-" form:"-"`
}

type DictDataDeleteReq struct {
	DictCodes []int64 `json:"dictCodes" form:"dictCodes" validate:"required"`
}
