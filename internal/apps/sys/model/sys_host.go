package model

type SysHost struct {
	HostId int64  `gorm:"column:host_id;primaryKey;autoIncrement" json:"hostId"`
	Name   string `gorm:"column:name;type:varchar(64);comment:'主机名称'" json:"name"`
	Alias  string `gorm:"column:alias;type:varchar(64);comment:'主机别名'" json:"alias"`
	Port   int    `gorm:"column:port;type:int;default:5921;comment:'主机别名'" json:"port"`
	Remark string `gorm:"column:remark;type:varchar(255);default:'';comment:'备注'" json:"remark"`

	Tasks []*SysTask `gorm:"many2many:sys_task_host;foreignKey:host_id;joinForeignKey:host_id;References:task_id;JoinReferences:task_id;" json:"tasks"`
}

func (e *SysHost) TableName() string {
	return "sys_host"
}
