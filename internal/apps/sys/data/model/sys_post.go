package model

type SysPost struct {
	PostId   int64  `gorm:"column:post_id;primaryKey;autoIncrement" json:"postId"`
	PostCode string `gorm:"column:post_code;type:varchar(64);comment:'岗位编码'" json:"postCode"`
	PostName string `gorm:"column:post_name;type:varchar(50);comment:'岗位名称'" json:"postName"`
	PostSort int64  `gorm:"column:post_sort;type:int;comment:'显示顺序'" json:"sort"`
	Remark   string `gorm:"column:remark;type:varchar(255);comment:'备注'" json:"remark"`
	Status   int64  `gorm:"column:status;type:tinyint(1);default:2;comment:'状态（2正常 1停用）'" json:"status"`

	Users []*SysUser `gorm:"many2many:sys_user_post;foreignKey:post_id;joinForeignKey:post_id;References:user_id;JoinReferences:user_id;" json:"users"`

	BaseModelTime
	ControlBy
}

func (e *SysPost) TableName() string {
	return "sys_post"
}
