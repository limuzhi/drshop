package model

type SysDictData struct {
	DictCode  int64  `gorm:"column:dict_code;primaryKey;autoIncrement" json:"dictCode"`
	DictSort  int    `gorm:"column:dict_sort;type:int(11);comment:'字典排序'" json:"dict_sort"`
	DictLabel string `gorm:"column:dict_label;type:varchar(128);comment:'字典标签'" json:"dictLabel"`
	DictValue string `gorm:"column:dict_value;type:varchar(128);comment:'字典键值'" json:"dictValue"`
	DictType  string `gorm:"column:dict_type;type:varchar(128);comment:'字典类型'" json:"dictType"`
	CssClass  string `gorm:"column:css_class;type:varchar(128);comment:'样式属性（其他样式扩展）'" json:"cssClass"`
	ListClass string `gorm:"column:list_class;type:varchar(128);comment:'表格回显样式'" json:"listClass"`
	Remark    string `gorm:"column:remark;type:varchar(255);comment:'备注'" json:"remark"`
	IsDefault int    `gorm:"column:is_default;type:tinyint(1);default:2;comment:'是否默认（1是 2否）'" json:"isDefault"`
	Status    int    `gorm:"column:status;type:tinyint(1);default:2;comment:'状态（2正常 1停用）'" json:"status"`
	BaseModelTime
	ControlBy
}

func (e *SysDictData) TableName() string {
	return "sys_dict_data"
}
