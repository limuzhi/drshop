package biz

import (
	"context"
	v1 "drpshop/api/sys/v1"
	"drpshop/internal/apps/sys/data/model"
	"drpshop/pkg/token"
	"errors"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

type SysUserRepo interface {
	//根据用户账号获取信息
	GetUsersForAccounts(ctx context.Context, account string) (*model.SysUser, error)
	//获取用户
	GetInfoById(ctx context.Context, uid int64) (*model.SysUser, error)
	//创建用户
	CreateUser(ctx context.Context, in *model.SysUser) error
	//更新用户
	UpdateUser(ctx context.Context, in *model.SysUser) error
	//列表用户信息
	ListUser(ctx context.Context, in *v1.UserListReq) ([]*model.SysUser, int64, error)
	//批量删除用户数据
	BatchDeleteUserByIds(ctx context.Context, ids []int64) error
	//更新密码
	ChangePwd(ctx context.Context, uid int64, salt, hashNewPasswd string) error
	//修改用户状态
	UpdateUserStatus(ctx context.Context, uid, status int64) error
	//用户权限接口
	GetPermissions(roleIds []int64) ([]string, error)
	//更新
	LoginUpdateUser(ctx context.Context, userId, loginTime int64, ip string) error
	////废弃掉了
	CreateUserPostRole(ctx context.Context, in *model.SysUser, postIds []int64, roleIds []int64) error
	// 创建jwttoken
	CreateToken(claims *token.UserClaims) (string, error)
	//获取个人信息
	GetUserProfile(ctx context.Context, uid int64) (*model.SysUser, error)
	//更新用个人信息
	UpdateUserProfile(ctx context.Context, uid int64, data map[string]interface{}) error
}

type SysUserUsecase struct {
	repo SysUserRepo
	log  *log.Helper
}

func NewSysUserUsecase(repo SysUserRepo, logger log.Logger) *SysUserUsecase {
	return &SysUserUsecase{repo: repo, log: log.NewHelper(log.With(logger, "module", "sys/biz/sys_user"))}
}

//根据用户账号获取信息
func (uc *SysUserUsecase) GetUsersForAccounts(ctx context.Context, account string) (*model.SysUser, error) {
	return uc.repo.GetUsersForAccounts(ctx, account)
}

//获取用户
func (uc *SysUserUsecase) GetInfoById(ctx context.Context, uid int64) (*model.SysUser, error) {
	return uc.repo.GetInfoById(ctx, uid)
}

//创建用户
func (uc *SysUserUsecase) CreateUser(ctx context.Context, in *model.SysUser) error {
	return uc.repo.CreateUser(ctx, in)
}

//更新用户
func (uc *SysUserUsecase) UpdateUser(ctx context.Context, in *model.SysUser) error {
	return uc.repo.UpdateUser(ctx, in)
}

//列表用户信息
func (uc *SysUserUsecase) ListUser(ctx context.Context, in *v1.UserListReq) ([]*model.SysUser, int64, error) {
	return uc.repo.ListUser(ctx, in)
}

//批量删除用户数据
func (uc *SysUserUsecase) BatchDeleteUserByIds(ctx context.Context, ids []int64) error {
	return uc.repo.BatchDeleteUserByIds(ctx, ids)
}

//更新密码
func (uc *SysUserUsecase) ChangePwd(ctx context.Context, uid int64, salt, hashNewPasswd string) error {
	return uc.repo.ChangePwd(ctx, uid, salt, hashNewPasswd)
}

//修改用户状态
func (uc *SysUserUsecase) UpdateUserStatus(ctx context.Context, uid, status int64) error {
	return uc.repo.UpdateUserStatus(ctx, uid, status)
}

//用户权限接口
func (uc *SysUserUsecase) GetPermissions(roleIds []int64) ([]string, error) {
	return uc.repo.GetPermissions(roleIds)
}

//更新
func (uc *SysUserUsecase) LoginUpdateUser(ctx context.Context, userId, loginTime int64, ip string) error {
	return uc.repo.LoginUpdateUser(ctx, userId, loginTime, ip)
}

////废弃掉了
func (uc *SysUserUsecase) CreateUserPostRole(ctx context.Context, in *model.SysUser, postIds []int64, roleIds []int64) error {
	return uc.repo.CreateUserPostRole(ctx, in, postIds, roleIds)
}

// 创建jwttoken
func (uc *SysUserUsecase) CreateToken(claims *token.UserClaims) (string, error) {
	return uc.repo.CreateToken(claims)
}

func (uc *SysUserUsecase) GetUserProfile(ctx context.Context, uid int64) (*model.SysUser, error) {
	return uc.repo.GetUserProfile(ctx, uid)
}

func (uc *SysUserUsecase) UpdateUserProfile(ctx context.Context, uid int64, data map[string]interface{}) error {
	return uc.repo.UpdateUserProfile(ctx, uid, data)
}

//===============

//用户登录
func (uc *SysUserUsecase) Login(ctx context.Context, account, password string) (outToken string, expireTime int64, err error) {
	userInfo, err := uc.repo.GetUsersForAccounts(ctx, account)
	if err != nil {
		return "", 0, err
	}
	if userInfo.Status != int64(model.Enabled) {
		return "", 0, errors.New("用户被禁用")
	}
	if userInfo.IsAdmin != 1 {
		return "", 0, errors.New("抱歉!您不属于后台管理员")
	}
	roleIds := make([]int64, 0)
	roleKeyArr := make([]string, 0)
	for _, role := range userInfo.Roles {
		if role.Status == 2 {
			roleIds = append(roleIds, role.RoleId)
			roleKeyArr = append(roleKeyArr, role.RoleKey)
		}
	}
	if len(roleIds) == 0 {
		return "", 0, errors.New("用户角色被禁用")
	}
	if userInfo.CheckPassword(password) {
		//生成TOKEN
		expireTime = time.Now().Add(time.Hour * 24 * 3).Unix()
		claims := &token.UserClaims{
			UserID:   userInfo.UserId,
			Username: userInfo.Username,
			RoleKey:  roleKeyArr,
			RoleIds:  roleIds,
		}
		claims.ExpiresAt = expireTime
		outToken, err = uc.repo.CreateToken(claims)
		return outToken, expireTime, err
	}else{
		return "", 0, errors.New("密码错误")
	}
}

func (uc *SysUserUsecase) UserInfo(ctx context.Context, uid int64) (*model.SysUser, error) {
	info, err := uc.repo.GetInfoById(ctx, uid)
	if err != nil {
		return nil, err
	}
	roleIds := make([]int64, 0)
	for _, v := range info.Roles {
		if v.Status != 2 {
			continue
		}
		roleIds = append(roleIds, v.RoleId)
	}
	permissions, _ := uc.repo.GetPermissions(roleIds)
	info.Permissions = permissions
	return info, nil
}

func (uc *SysUserUsecase) LastIdCreateUser(ctx context.Context, in *model.SysUser) (int64, error) {
	if err := uc.repo.CreateUser(ctx, in); err != nil {
		return 0, err
	} else {
		return in.UserId, nil
	}
}
