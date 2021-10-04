package model

type SysLoginLog struct {
	LoginId       int64  `gorm:"column:login_id;primaryKey;autoIncrement" json:"loginId"`
	LoginName     string `gorm:"column:login_name;type:varchar(50);not null;comment:'登录账号'" json:"loginName"`
	LoginUid      int64  `gorm:"column:login_uid;type:bigint;not null;comment:'登录账号ID'" json:"loginUid"`
	Ipaddr        string `gorm:"column:ipaddr;type:varchar(50);not null;comment:'登录IP地址'" json:"ipaddr"`
	LoginLocation string `gorm:"column:login_location;type:varchar(255);not null;comment:'登录地点'" json:"loginLocation"`
	Browser       string `gorm:"column:browser;type:varchar(50);not null;comment:'浏览器类型'" json:"browser"`
	Os            string `gorm:"column:os;type:varchar(50);not null;comment:'操作系统'" json:"os"`
	Msg           string `gorm:"column:msg;type:varchar(255);not null;comment:'提示消息'" json:"msg"`
	Status        int32  `gorm:"column:status;type:tinyint;default:2;comment:'登录状态（2成功 1失败）'" json:"status"`
	LoginTime     int64  `gorm:"column:login_time;type:int;default:0;comment:'登录时间'" json:"loginTime"`
}

func (e *SysLoginLog) TableName() string {
	return "sys_login_log"
}

// SysLoginLogSearchReq 查询列表请求参数
type SysLoginLogSearchReq struct {
	LoginName     string //登陆名
	Status        int32  //登录状态（2成功 1失败）'
	Ipaddr        string //登录地址
	LoginLocation string
	BeginTime     int64
	EndTime       int64
	SortName      string
	SortOrder     string
	PageNum       int
	PageSize      int
}
