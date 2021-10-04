package model

type SysRole struct {
	RoleId    int64  `gorm:"column:role_id;primaryKey;autoIncrement" json:"roleId"`
	Status    int64  `gorm:"column:status;type:tinyint;default:2;comment:'状态（2正常 1停用）'" json:"status"`
	Sort      int64  `gorm:"column:sort;type:int;default:50;comment:'排序'" json:"sort"`
	PID       int64  `gorm:"column:pid;type:int;default:50;comment:'父角色ID'" json:"pid"`
	Name      string `gorm:"column:name;type:varchar(20);not null;comment:'角色名称'" json:"name"`
	RoleKey   string `gorm:"column:role_key;type:varchar(128);not null;comment:'权限字符'" json:"roleKey"`
	Remark    string `gorm:"column:remark;type:varchar(255);not null;comment:'备注'" json:"remark"`
	DataScope int32  `gorm:"column:data_scope;type:varchar(255);not null;comment:'数据范围（1：全部数据权限 2：自定数据权限 3：本部门数据权限 4：本部门及以下数据权限）'" json:"dataScope"`

	Users []*SysUser `gorm:"many2many:sys_user_roles;foreignKey:role_id;joinForeignKey:role_id;References:user_id;JoinReferences:user_id;" json:"users"`
	Depts []*SysDept `gorm:"many2many:sys_role_dept;foreignKey:role_id;joinForeignKey:role_id;References:dept_id;JoinReferences:dept_id;" json:"depts"`
	Menus []*SysMenu `gorm:"many2many:sys_role_menus;foreignKey:role_id;joinForeignKey:role_id;References:menu_id;JoinReferences:menu_id;" json:"menus"`
}

func (e *SysRole) TableName() string {
	return "sys_role"
}
