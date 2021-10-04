package vo

type DeptListReq struct {
	DeptName string `json:"deptName" form:"deptName"`
	Status   string `json:"status" form:"status"`
}

type DeptDetailReq struct {
	DeptId int64 `json:"deptId" form:"deptId"`
}

type DeptCreateReq struct {
	ParentId  int64  `json:"parentId" form:"parentId" validate:""`
	DeptName  string `json:"deptName" form:"deptName" validate:"required,min=1,max=30`
	Sort      int    `json:"sort" form:"sort" validate:"gte=1,lte=999"`
	Leader    string `json:"leader" form:"leader" validate:"required,min=1,max=30`
	Phone     string `json:"phone" form:"phone" validate:"required"`
	Email     string `json:"email" form:"email" validate:"required,email"`
	Status    string `json:"status" form:"status" validate:"oneof=1 2"`
	Ancestors string `json:"ancestors" form:"ancestors"`
	CreateBy  int64  `json:"-" form:"-"`
}

type DeptUpdateReq struct {
	DeptCreateReq
	DeptId   int64 `json:"deptId" form:"deptId" validate:"required,gte=1"`
	UpdateBy int64 `json:"-" form:"-"`
}

type DeptDeleteReq struct {
	DeptIds []int64 `json:"deptIds" form:"deptIds" validate:"required"`
}
