package vo

type RoleSearchReq struct {
	Pagination
	Name    string `json:"name" form:"name"`       //角色名称
	RoleKey string `json:"roleKey" form:"roleKey"` //权限字符
	Status  string `json:"status" form:"status"`   //状态
}

type RoleAddReq struct {
	Name    string `json:"name" form:"name" validate:"required,min=1,max=30"`
	RoleKey string `json:"roleKey" form:"roleKey" validate:"required,min=1,max=20"`
	Remark  string `json:"remark" form:"remark" validate:"min=0,max=100"`
	Status  string `json:"status" form:"status" validate:"oneof=1 2"`
	Sort    int64  `json:"sort" form:"sort" validate:"gte=1,lte=999"`
}

type RoleUpdateReq struct {
	RoleId  int64  `json:"roleId" form:"roleId" validate:"required,gte=1"`
	Name    string `json:"name" form:"name" validate:"required,min=1,max=30"`
	RoleKey string `json:"roleKey" form:"roleKey" validate:"required,min=1,max=20"`
	Remark  string `json:"remark" form:"remark" validate:"min=0,max=100"`
	Status  string `json:"status" form:"status" validate:"oneof=1 2"`
	Sort    int64  `json:"sort" form:"sort" validate:"gte=1,lte=999"`
}

type RoleBatchDeleteReq struct {
	RoleIds []int64 `json:"roleIds" form:"roleIds" validate:"required"` //角色ID
}

type RoleStatusReq struct {
	RoleId int64 `json:"roleId" form:"roleId" validate:"required"` //角色ID
	Status int64 `json:"status" form:"status" validate:"required"` //状态
}

type RoleMenusApisListReq struct {
	RoleId int64 `json:"roleId" form:"roleId" validate:"required,gte=1"` //角色ID
}

type RoleMenusUpdateReq struct {
	RoleId  int64   `json:"roleId" form:"roleId" validate:"required,gte=1"` //角色ID
	MenuIds []int64 `json:"menuIds" form:"menuIds" validate:""`     //菜单Ids
}

type RoleApisUpdateReq struct {
	RoleId int64   `json:"roleId" form:"roleId" validate:"required,gte=1"` //角色ID
	ApiIds []int64 `json:"apiIds" form:"apiIds" validate:""`       //接口Ids
}
