package model


type SysCasbinRule struct {
	PType   string `gorm:"column:p_type;type:varchar(100);"`
	VO   string `gorm:"column:v0;type:varchar(100);"`
	V1   string `gorm:"column:v1;type:varchar(100);"`
	V2  string `gorm:"column:v2;type:varchar(100);"`
	V3   string `gorm:"column:v3;type:varchar(100);"`
	V4   string `gorm:"column:v4;type:varchar(100);"`
	V5   string `gorm:"column:v5;type:varchar(100);"`
}

func (e *SysCasbinRule) TableName() string {
	return "sys_casbin_rule"
}