package model

type SysUserRoles struct {
	UserId int64 `gorm:"column:user_id;" json:"userId"`
	RoleId int64 `gorm:"column:role_id;" json:"role_id"`
}

func (e *SysUserRoles) TableName() string {
	return "sys_user_roles"
}
