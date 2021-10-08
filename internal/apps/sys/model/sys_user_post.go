package model

type SysUserPost struct {
	UserId        int64  `gorm:"column:user_id;" json:"userId"`
	PostId        int64  `gorm:"column:post_id;" json:"post_id"`
}

func (e *SysUserPost) TableName() string {
	return "sys_user_post"
}
