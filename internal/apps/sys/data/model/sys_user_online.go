package model

type SysUserOnline struct {
	ID         int64  `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Username   string `gorm:"column:username;type:varchar(255);comment:'用户名'" json:"username"`
	UserId     int64  `gorm:"column:user_id;type:bigint;comment:'标题说明'" json:"userId"`
	UUID       string `gorm:"column:uuid;type:char(32);comment:'用户标识'" json:"uuid"`
	Token      string `gorm:"column:token;type:varchar(255);comment:'用户token'" json:"token"`
	Ip         string `gorm:"column:ip;type:varchar(120);comment:'登录ip'" json:"ip"`
	Explorer   string `gorm:"column:explorer;type:varchar(30);comment:'浏览器'" json:"explorer"`
	Os         string `gorm:"column:os;type:varchar(30);comment:'操作系统'" json:"os"`
	CreateTime int64  `gorm:"column:create_time;type:int;comment:'操作系统'" json:"createTime"`
}

func (e *SysUserOnline) TableName() string {
	return "sys_user_online"
}
