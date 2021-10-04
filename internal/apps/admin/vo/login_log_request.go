package vo

type LoginLogSearchReq struct {
	LoginName     string `json:"loginName" form:"loginName"`
	Ipaddr        string `json:"ipaddr" form:"ipaddr"`
	LoginLocation string `json:"loginLocation" form:"loginLocation"`
	Status        int32  `json:"status" form:"status"`
	Pagination
}

type LoginLogDeleteReq struct {
	LoginIds []int64 `json:"loginIds" form:"loginIds"`
}
