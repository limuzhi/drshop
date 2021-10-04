package vo

//用户登录
type LoginReq struct {
	Account  string `form:"account" json:"account" validate:"required"`
	Password string `form:"password" json:"password" validate:"required"`
	Code     string `form:"code" json:"code" validate:"required"`
	UUID     string `form:"uuid" json:"uuid" validate:"required"`
}

//短信验证码
type CaptchaReq struct {
	To        string `form:"to" json:"to" validate:"required"`               //接收对象 email
	Captcha   string `form:"captcha" json:"captcha" validate:"required"`     //本地验证码
	CaptchaId string `form:"captchaId" json:"captchaId" validate:"required"` //本地验证码
}

type UserSearchReq struct {
	Pagination
	DeptId   int    `form:"deptId" json:"deptId"`
	Mobile   string `form:"mobile" json:"mobile"`
	Status   string `form:"status" json:"status"` //用户状态;1:禁用,2:正常,3:未验证
	Username string `form:"username" json:"username"`
}

type UserChangeStatusReq struct {
	UserId int64 `form:"userId" json:"userId" comment:"用户ID" validate:"required" example:""`
	Status int64 `form:"status" json:"status" comment:"状态 2-正常、1-停用" validate:"required,oneof=1 2" example:""`
}

type BatchDeleteUserReq struct {
	UserIds []int64 `form:"userIds" json:"userIds" comment:"用户ID" validate:"required" example:""`
}

// SetUserReq 添加修改用户公用请求字段
type SetUserReq struct {
	DeptId   int64   `form:"deptId" json:"deptId" validate:"required,gte=1"`
	Nickname string  `form:"nickname" json:"nickname" validate:"required"`
	Mobile   string  `form:"mobile" json:"mobile" validate:"required"`
	Email    string  `form:"mobile" json:"email" validate:"email"`
	Sex      string  `form:"sex" json:"sex" validate:""`
	Avatar   string  `form:"avatar" json:"avatar"`
	Remark   string  `form:"remark" json:"remark"`
	Status   string  `form:"status" json:"status" validate:""`
	IsAdmin  int64   `form:"isAdmin" json:"isAdmin"`
	PostIds  []int64 `form:"postIds" json:"postIds"`
	RoleIds  []int64 `form:"roleIds" json:"roleIds"`
}

//  添加用户参数
type CreateUserReq struct {
	SetUserReq
	Username string `form:"username" json:"username" validate:"required"`
	Password string `form:"password" json:"password" validate:"required"`
	CreateBy int64
}

type UpdateUserReq struct {
	SetUserReq
	UpdateBy int64
	UserId   int64 `form:"userId" json:"userId" validate:"required,gte=1"`
}

// 更新密码结构体
type ChangePwdReq struct {
	OldPassword string `json:"oldPassword" form:"oldPassword" validate:"required,min=6,max=50"`
	NewPassword string `json:"newPassword" form:"newPassword" validate:"required,min=6,max=50"`
	UserId      int64  `json:"-" form:"-"'`
}

// UserProfileSetReq 个人信息修改
type UserProfileSetReq struct {
	Nickname string `form:"nickname" json:"nickname" validate:"required"`
	Mobile   string `form:"mobile" json:"mobile" validate:"required"`
	Email    string `form:"mobile" json:"email" validate:"email"`
	Sex      int64  `form:"sex" json:"sex" validate:""`
	Avatar   string `form:"avatar" json:"avatar"`
	UserId   int64  `json:"-" form:"-"'`
}
