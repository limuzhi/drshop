package vo

//操作日志搜索
type OperLogSearchReq struct {
	Title        string `json:"title" form:"title"`
	OperName     string `json:"operName" form:"operName"`
	Status       string `json:"status" form:"status"`
	BusinessType string `json:"businessType" form:"businessType"`
	Pagination
}

type OperLogDetailReq struct {
	OperId int64 `json:"operId" form:"operId" validate:"required,gte=1"`
}

type OperLogDeleteReq struct {
	OperIds []int64 `json:"operIds" form:"operIds" validate:"required"`
}
