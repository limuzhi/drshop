package model

type SysMenu struct {
	MenuId     int64  `gorm:"column:menu_id;primaryKey;autoIncrement" json:"menuId"`
	Pid        int64  `gorm:"column:pid;type:bigint;default:0;comment:'父ID'" json:"pid"`
	Name       string `gorm:"column:name;type:varchar(128);default:'';comment:'菜单名称(英文名, 可用于国际化)''" json:"name"`
	Title      string `gorm:"column:title;type:varchar(50);default:'';comment:'菜单标题(无法国际化时使用)'" json:"title"`
	Icon       string `gorm:"column:icon;type:varchar(50);default:'';comment:'菜单图标'" json:"icon"`
	Sort       int    `gorm:"column:sort;type:int;comment:'排序 越小越前'" json:"sort"`
	Hidden     int    `gorm:"column:hidden;type:tinyint;default:2;comment:'显示状态 1隐藏 2显示'" json:"hidden"`
	NoCache    int    `gorm:"column:no_cache;type:tinyint;default:2;comment:'菜单是否被 <keep-alive> 缓存(1不缓存，2缓存)'" json:"noCache"`
	Breadcrumb int    `gorm:"column:breadcrumb;type:tinyint;default:1;comment:'面包屑可见性(可见/隐藏, 默认可见)'" json:"breadcrumb"`
	AlwaysShow int    `gorm:"column:always_show;type:tinyint;default:2;comment:'忽略之前定义的规则，一直显示根路由(1忽略，2不忽略)'" json:"alwaysShow"`
	Path       string `gorm:"column:path;type:varchar(255);comment:'路由地址'" json:"path"`
	Redirect   string `gorm:"column:redirect;type:varchar(255);comment:'重定向路径'" json:"redirect"`
	ActiveMenu string `gorm:"column:active_menu;type:varchar(255);comment:'在其它路由时，想在侧边栏高亮的路由'" json:"activeMenu"`
	JumpPath   string `gorm:"column:jump_path;type:varchar(125);comment:'跳转路由'" json:"jumpPath"`
	Component  string `gorm:"column:component;type:varchar(125);comment:'组件路径'" json:"component"`
	IsFrame    int    `gorm:"column:is_frame;type:tinyint;default:0;comment:'是否外链 1是 0否'" json:"isFrame"`
	ModuleType string `gorm:"column:module_type;type:varchar(30);default:'';comment:'所属模块'" json:"moduleType"`
	ModelId    int    `gorm:"column:model_id;type:int;default:0;comment:'模型ID'" json:"modelId"`
	Status     int    `gorm:"column:status;type:tinyint(1);default:2;comment:'状态（2正常 1停用）'" json:"status"`

	BaseModelTime
	ControlBy
	Children []*SysMenu `gorm:"-" json:"children"` // 子菜单集合
	RoleList []*SysRole `gorm:"many2many:sys_role_menus;foreignKey:menu_id;joinForeignKey:menu_id;References:role_id;JoinReferences:role_id;" json:"roleList"`
}

func (e *SysMenu) TableName() string {
	return "sys_menu"
}

type SysMenuSlice []*SysMenu

func (a SysMenuSlice) Len() int { // 重写 Len() 方法
	return len(a)
}
func (a SysMenuSlice) Swap(i, j int) { // 重写 Swap() 方法
	a[i], a[j] = a[j], a[i]
}
func (a SysMenuSlice) Less(i, j int) bool { // 重写 Less() 方法， 从小到大排序
	return a[j].Sort > a[i].Sort
}

//sort.Sort(SysMenuSlice(list))    // 按照 Age 的逆序排序
//    fmt.Println(people)
//
// sort.Sort(sort.Reverse(SysMenuSlice(list)))    // 按照 sort 的排序
//fmt.Println(people)
