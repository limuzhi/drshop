package model

import (
	"crypto/hmac"
	"crypto/md5"
	"drpshop/pkg/util"
	"encoding/hex"
	"fmt"
)

type SysUser struct {
	UserId        int64  `gorm:"column:user_id;primaryKey;autoIncrement" json:"userId"`
	Username      string `gorm:"column:username;type:varchar(60);not null;unique;comment:'用户名'" json:"username"`
	Mobile        string `gorm:"column:mobile;type:varchar(20);not null;comment:'手机号'" json:"mobile"`
	Nickname      string `gorm:"column:nickname;type:varchar(50);not null;comment:'用户昵称'" json:"nickname"`
	Birthday      int64  `gorm:"column:birthday;type:int(11);not null;default:0;comment:'生日'" json:"birthday"`
	DeptId        int64  `gorm:"column:dept_id;type:bigint;not null;default:0;comment:'部门id'" json:"deptId"`
	Password      string `gorm:"column:password;type:varchar(255);not null;default:'';comment:'登录密码'" json:"password"`
	Salt          string `gorm:"column:salt;type:varchar(50);comment:'加密盐'" json:"salt"`
	Status        int64  `gorm:"column:status;type:tinyint;default:2;comment:'用户状态;1:禁用,2:正常,3:未验证'" json:"status"`
	Avatar        string `gorm:"column:avatar;type:varchar(255);not null;default:'';comment:'用户头像'" json:"avatar"`
	Sex           int    `gorm:"column:sex;type:tinyint;not null;default:0;comment:'性别;0:保密,1:男,2:女'" json:"sex"`
	Email         string `gorm:"column:email;type:varchar(100);not null;default:'';comment:'用户登录邮箱'" json:"email"`
	Remark        string `gorm:"column:remark;type:varchar(255);not null;default:'';comment:'备注'" json:"remark"`
	IsAdmin       int32  `gorm:"column:is_admin;type:tinyint;not null;default:1;comment:'是否后台管理员1 是 0否'" json:"isAdmin"`
	Address       string `gorm:"column:address;type:varchar(255);not null;default:'';comment:'联系地址'" json:"address"`
	LastLoginIp   string `gorm:"column:last_login_ip;type:varchar(15);not null;default:'';comment:'最后登录ip'" json:"lastLoginIp"`
	LastLoginTime int64  `gorm:"column:last_login_time;type:int;comment:'最后登录时间'" json:"lastLoginTime"`
	LoginCount    int    `gorm:"column:login_count;type:int;comment:'登录次数'" json:"loginCount"`
	BaseModelTime
	ControlBy

	Roles       []*SysRole `gorm:"many2many:sys_user_roles;foreignKey:user_id;joinForeignKey:user_id;References:role_id;JoinReferences:role_id;" json:"roles"`
	Posts       []*SysPost `gorm:"many2many:sys_user_post;foreignKey:user_id;joinForeignKey:user_id;References:post_id;JoinReferences:post_id;" json:"posts"`
	Permissions []string   `gorm:"-" json:"permissions"`
}

func (e *SysUser) TableName() string {
	return "sys_user"
}


// 检查密码
func (e *SysUser) CheckPassword(password string) bool {
	fmt.Println("=====e.GenPassword(password) == e.Password:",e.GenPassword(password) == e.Password,"====")
	return e.GenPassword(password) == e.Password
}

// 创建密码
func (e *SysUser) SetPassword(password string) {
	e.Salt = util.RandomString(4)
	e.Password = e.GenPassword(password)
}

// 创建密码
func (e *SysUser) GenPassword(password string) string {
	hmacEnt := hmac.New(md5.New, []byte(e.Salt))
	hmacEnt.Write([]byte(password))
	str := hex.EncodeToString(hmacEnt.Sum([]byte("")))
	fmt.Println("=====password:",password,"==pwd:",str,"=====e.password:",e.Password,"===")
	return str
}
