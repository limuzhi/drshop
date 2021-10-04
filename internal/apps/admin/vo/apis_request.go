package vo

// 获取接口列表结构体
type ApiListReq struct {
	Method   string `json:"method" form:"method"`
	Path     string `json:"path" form:"path"`
	Category string `json:"category" form:"category"`
	Pagination
}

// 创建接口结构体
type CreateApiReq struct {
	Handle   string `json:"handle" form:"handle" validate:""`
	Permission string `json:"permission" form:"permission" validate:""`
	Method   string `json:"method" form:"method" validate:"required,min=1,max=20"`
	Path     string `json:"path" form:"path" validate:"required,min=1,max=100"`
	Category string `json:"category" form:"category" validate:"required,min=1,max=50"`
	Title    string `json:"title" form:"title" validate:"min=0,max=100"`
	CreateBy int64  `json:"-" form:"-"`
}

// 更新接口结构体
type UpdateApiReq struct {
	ApiId      int64  `json:"apiId" form:"apiId" validate:"required,gte=1"`
	Handle     string `json:"handle" form:"handle" validate:""`
	Permission string `json:"permission" form:"permission" validate:""`
	Method     string `json:"method" form:"method" validate:"min=1,max=20"`
	Path       string `json:"path" form:"path" validate:"min=1,max=100"`
	Category   string `json:"category" form:"category" validate:"min=1,max=50"`
	Title      string `json:"title" form:"title" validate:"min=0,max=100"`
	UpdateBy   int64  `json:"-" form:"-"`
}

// 批量删除接口结构体
type BatchDeleteApiReq struct {
	ApiIds []int64 `json:"apiIds" form:"apiIds"`
}
