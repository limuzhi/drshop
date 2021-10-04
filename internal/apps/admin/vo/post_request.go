package vo

type PostListReq struct {
	PostCode string `json:"postCode" form:"postCode"`
	PostName string `json:"postName" form:"postName"`
	Status   string `json:"status" form:"status"`
	Pagination
}

type PostCreateReq struct {
	PostCode string `json:"postCode" form:"postCode" validate:"required,min=1,max=50`
	PostName string `json:"postName" form:"postName" validate:"required,min=1,max=50`
	PostSort int    `json:"postSort" form:"postSort" validate:"gte=1,lte=999"`
	Status   int    `json:"status" form:"status" validate:"oneof=1 2"`
	Remark   string `json:"remark" form:"remark"`
	CreateBy int64  `json:"-" form:"-"`
}

type PostUpdateReq struct {
	PostCreateReq
	PostId   int64 `json:"postId" form:"postId" validate:"required,gte=1"`
	UpdateBy int64 `json:"-" form:"-"`
}

type PostDeleteReq struct {
	PostIds []int64 `json:"postIds" form:"postIds" validate:"required"`
}

type PostDetail struct {
	PostId   int64 `json:"postId" form:"postId" validate:"required,gte=1"`
}
