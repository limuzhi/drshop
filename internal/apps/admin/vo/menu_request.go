package vo

// 创建结构体
type CreateMenuReq struct {
	Name       string `json:"name" form:"name" validate:"required,min=1,max=100"`
	Title      string `json:"title" form:"title" validate:"required,min=1,max=50"`
	Icon       string `json:"icon" form:"icon" validate:"min=0,max=50"`
	Path       string `json:"path" form:"path" validate:"required,min=1,max=100"`
	Redirect   string `json:"redirect" form:"redirect" validate:"min=0,max=250"`
	Component  string `json:"component" form:"component" validate:"required,min=1,max=100"`
	Sort       int    `json:"sort" form:"sort" validate:"gte=1,lte=999"`
	Status     int    `json:"status" form:"status" validate:"oneof=1 2"`
	Hidden     int    `json:"hidden" form:"hidden" validate:"oneof=1 2"`
	NoCache    int    `json:"noCache" form:"noCache" validate:"oneof=1 2"`
	AlwaysShow int    `json:"alwaysShow" form:"alwaysShow" validate:"oneof=1 2"`
	Breadcrumb int    `json:"breadcrumb" form:"breadcrumb" validate:"oneof=1 2"`
	ActiveMenu string `json:"activeMenu" form:"activeMenu" validate:"min=0,max=100"`
	Pid        int64  `json:"pid" form:"pid"`
	UserId     int64  `json:"-" form:"-"`

	//JumpPath   string `json:"jumpPath" form:"jumpPath" validate:""`
	//IsFrame    int    `json:"isFrame" form:"isFrame" validate:""`
	//ModuleType string `json:"moduleType" form:"moduleType" validate:""`
	//ModelId    int    `json:"modelId" form:"modelId" validate:""`
}

// 更新结构体
type UpdateMenuReq struct {
	MenuId     int64  `json:"menuId" form:"menuId" validate:"required,gte=1"`
	Name       string `json:"name" form:"name" validate:"required,min=1,max=100"`
	Title      string `json:"title" form:"title" validate:"required,min=1,max=50"`
	Icon       string `json:"icon" form:"icon" validate:"min=0,max=50"`
	Path       string `json:"path" form:"path" validate:"required,min=1,max=100"`
	Redirect   string `json:"redirect" form:"redirect" validate:"min=0,max=250"`
	Component  string `json:"component" form:"component" validate:"required,min=1,max=100"`
	Sort       int    `json:"sort" form:"sort" validate:"gte=1,lte=999"`
	Status     int    `json:"status" form:"status" validate:"oneof=1 2"`
	Hidden     int    `json:"hidden" form:"hidden" validate:"oneof=1 2"`
	NoCache    int    `json:"noCache" form:"noCache" validate:"oneof=1 2"`
	AlwaysShow int    `json:"alwaysShow" form:"alwaysShow" validate:"oneof=1 2"`
	Breadcrumb int    `json:"breadcrumb" form:"breadcrumb" validate:"oneof=1 2"`
	ActiveMenu string `json:"activeMenu" form:"activeMenu" validate:"min=0,max=100"`
	Pid        int64  `json:"pid" form:"pid"`
	UserId     int64  `json:"-" form:"-"`
	//JumpPath   string `json:"jumpPath" form:"jumpPath" validate:""`
	//IsFrame    int    `json:"isFrame" form:"isFrame" validate:""`
	//ModuleType string `json:"moduleType" form:"moduleType" validate:""`
	//ModelId    int    `json:"modelId" form:"modelId" validate:""`
}

// 删除结构体
type DeleteMenuReq struct {
	MenuIds []int64 `json:"menuIds" form:"menuIds"`
}
