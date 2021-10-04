package model

type SysDept struct {
	DeptId    int64  `gorm:"column:dept_id;primaryKey;autoIncrement" json:"deptId"`
	ParentId  int64  `gorm:"column:parent_id;type:bigint;default:0;comment:'父部门id'" json:"parentId"`
	Ancestors string `gorm:"column:ancestors;type:varchar(255);comment:'祖级列表 多个以逗号分割'" json:"ancestors"`
	DeptName  string `gorm:"column:dept_name;type:varchar(30);comment:'部门名称'" json:"deptName"`
	Sort      int    `gorm:"column:sort;type:int;comment:'显示顺序'" json:"sort"`
	Leader    string `gorm:"column:leader;type:varchar(30);comment:'负责人'" json:"leader"`
	Phone     string `gorm:"column:phone;type:varchar(20);comment:'联系电话'" json:"phone"`
	Email     string `gorm:"column:email;type:varchar(128);comment:'邮箱'" json:"email"`
	Status    int    `gorm:"column:status;type:int;default:2;comment:'部门状态（2正常 1停用）'" json:"status"`
	BaseModelTime
	ControlBy

	RoleList []*SysRole `gorm:"many2many:sys_role_dept;foreignKey:dept_id;joinForeignKey:dept_id;References:role_id;JoinReferences:role_id;" json:"roleList"`

	Children []*SysDept `gorm:"-" json:"children"`
}

func (e *SysDept) TableName() string {
	return "sys_dept"
}

type SysRoleDept struct {
	RoleId int64 `gorm:"column:role_id;" json:"roleId"`
	DeptId int64 `gorm:"column:dept_id;" json:"deptId"`
}

func (e *SysRoleDept) TableName() string {
	return "sys_role_dept"
}
