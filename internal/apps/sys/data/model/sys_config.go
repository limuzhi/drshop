package model

type SysConfig struct {
	ConfigId    int64  `gorm:"column:config_id;primaryKey;autoIncrement" json:"configId"`
	ConfigName  string `gorm:"column:config_name;type:varchar(100);comment:'参数名称'" json:"configName"`
	ConfigKey   string `gorm:"column:config_key;type:varchar(100);comment:'参数键名'" json:"configKey"`
	ConfigValue string `gorm:"column:config_value;type:varchar(500);comment:'参数键值'" json:"configValue"`
	ConfigType  int    `gorm:"column:config_type;type:tinyint;comment:'系统内置（1是 0否）'" json:"configType"`
	IsFrontend  int    `gorm:"column:is_frontend;type:tinyint;comment:'是否前台 （1是 0否）'" json:"isFrontend"`
	Remark      string `gorm:"column:remark;type:varchar(500);comment:'备注'" json:"remark"`
	BaseModelTime
	ControlBy
}

func (e *SysConfig) TableName() string {
	return "sys_config"
}
