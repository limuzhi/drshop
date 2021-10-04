package model


//分页请求参数
type SysConfigSearchReq struct {
	ConfigName string
	ConfigKey  string
	ConfigType string
	BeginTime  string
	EndTime    string
	PageNum    int
	PageSize   int
}

type SysDeptSearchParams struct {
	ExcludeId int64
	DeptName  string
	Status    int
}

type SelectDictPageReq struct {
	DictType  string
	DictLabel string
	Status    int
	PageNum   int
	PageSize  int
}

type ListSysDictTypeReq struct {
	DictName string
	DictType string
	Status   int
	PageNum  int
	PageSize int
}

type SysJobSearchReq struct {
	JobName  string
	JobGroup string
	Status   int
	PageNum  int
	PageSize int
}



type SysPostPageReq struct {
	PostName string
	PostCode string
	Status   int
	PageNum  int
	PageSize int
}

// SysUserOnlineSearchReq 列表搜索参数
type SysUserOnlineSearchReq struct {
	Username string
	Ip       string
	PageNum  int
	PageSize int
}

type SysTaskReq struct {
	TaskId   int64
	Name     string
	Tag      string
	Protocol int
	Status   int
	HostId   int64
	PageNum  int
	PageSize int
}

type SysHostReq struct {
	Name     string
	PageNum  int
	PageSize int
}
