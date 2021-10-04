package model

type SysApis struct {
	ApiId      int64  `gorm:"column:api_id;primaryKey;autoIncrement" json:"apiId"`
	Handle     string `gorm:"column:handle;type:varchar(128);comment:'handle'" json:"handle"`
	Title      string `gorm:"column:title;type:varchar(128);comment:'标题说明'" json:"title"`
	Permission string `gorm:"column:permission;type:varchar(128);comment:'按钮权限标识'" json:"permission"`
	Path       string `gorm:"column:path;type:varchar(128);comment:'访问路径'" json:"path"`
	Method     string `gorm:"column:method;type:varchar(128);comment:'接口类型-请求方式'" json:"method"`
	Category   string `gorm:"column:category;type:varchar(50);comment:'所属类别'" json:"category"`
	BaseModelTime
	ControlBy
}

func (e *SysApis) TableName() string {
	return "sys_apis"
}
