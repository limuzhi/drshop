package model

type SysDictType struct {
	DictId   int64  `gorm:"column:dict_id;primaryKey;autoIncrement" json:"dictId"`
	DictName string `gorm:"column:dict_name;type:varchar(128);comment:'字典名称'" json:"dictName"`
	DictType string `gorm:"column:dict_type;type:varchar(128);comment:'字典类型'" json:"dictType"`
	Remark    string `gorm:"column:remark;type:varchar(255);comment:'备注'" json:"remark"`
	Status   int    `gorm:"column:status;type:tinyint(1);default:2;comment:'状态（2正常 1停用）'" json:"status"`

	BaseModelTime
	ControlBy
}

func (e *SysDictType) TableName() string {
	return "sys_dict_type"
}
