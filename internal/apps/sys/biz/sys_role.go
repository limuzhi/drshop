package biz

import (
	"context"
	v1 "drpshop/api/sys/v1"
	"drpshop/internal/apps/sys/data/model"
	"github.com/go-kratos/kratos/v2/log"
)

type SysRoleRepo interface {
	//获取
	GetInfo(ctx context.Context, roleIds []int64) ([]*model.SysRole, error)
	//创建
	CreateRole(ctx context.Context, in *model.SysRole) error
	//修改
	UpdateRole(ctx context.Context, in *model.SysRole) error
	//修改状态
	ChangeRoleStatus(ctx context.Context, id, status int64) error
	//列表
	ListRole(ctx context.Context, in *v1.RoleListReq) ([]*model.SysRole, int64, error)
	//删除
	BatchDeleteRoleByIds(ctx context.Context, ids []int64) error

	//获取角色对应的菜单
	GetRoleMenuList(ctx context.Context, roleIds []int64) ([]*model.SysRole, error)
	//更新角色的权限菜单
	UpdateRoleMenus(ctx context.Context, in *model.SysRole) error

	// 根据角色关键字获取角色的权限接口--casbinEnforcer
	GetRoleApisByRoleKeyword(roleKey string) ([]*model.SysApis, error)
	//根据角色ID获取角色
	GetRoleApisList(ctx context.Context, roleIds []int64) ([]*model.SysRole, error)
	// 更新角色的权限接口（先全部删除再新增）
	UpdateRoleApis(roleKeyword string, reqRolePolicies [][]string) error
	//所有开启的角色列表
	AllListRole(ctx context.Context) ([]*model.SysRole, error)
	//根据角色ID获取
	ListRoleByIds(ctx context.Context, ids []int64) ([]*model.SysRole, error)
	//获取策略中的所有授权规则
	GetFilteredPolicy(roleKey string) [][]string

	CheckCasbinRole(subs []string, obj string, act string) bool
	//获取按钮权限列表
	GetPermissionsByRoleKey(roleKey string) ([]string, error)
}

type SysRoleUsecase struct {
	repo SysRoleRepo
	log  *log.Helper
}

func NewSysRoleUsecase(repo SysRoleRepo, logger log.Logger) *SysRoleUsecase {
	return &SysRoleUsecase{repo: repo, log: log.NewHelper(log.With(logger, "module", "sys/biz/sys_role"))}
}

//获取
func (uc *SysRoleUsecase) GetInfo(ctx context.Context, roleIds []int64) ([]*model.SysRole, error) {
	return uc.repo.GetInfo(ctx, roleIds)
}

//创建
func (uc *SysRoleUsecase) CreateRole(ctx context.Context, in *model.SysRole) error {
	return uc.repo.CreateRole(ctx, in)
}

//修改
func (uc *SysRoleUsecase) UpdateRole(ctx context.Context, in *model.SysRole) error {
	return uc.repo.UpdateRole(ctx, in)
}

//修改状态
func (uc *SysRoleUsecase) ChangeRoleStatus(ctx context.Context, id, status int64) error {
	return uc.repo.ChangeRoleStatus(ctx, id, status)
}

//列表
func (uc *SysRoleUsecase) ListRole(ctx context.Context, in *v1.RoleListReq) ([]*model.SysRole, int64, error) {
	return uc.repo.ListRole(ctx, in)
}

//删除
func (uc *SysRoleUsecase) BatchDeleteRoleByIds(ctx context.Context, ids []int64) error {
	return uc.repo.BatchDeleteRoleByIds(ctx, ids)
}

//获取角色对应的菜单
func (uc *SysRoleUsecase) GetRoleMenuList(ctx context.Context, roleIds []int64) ([]*model.SysRole, error) {
	return uc.repo.GetRoleMenuList(ctx, roleIds)
}

//更新角色的权限菜单
func (uc *SysRoleUsecase) UpdateRoleMenus(ctx context.Context, in *model.SysRole) error {
	return uc.repo.UpdateRoleMenus(ctx, in)
}

// 根据角色关键字获取角色的权限接口--casbinEnforcer
func (uc *SysRoleUsecase) GetRoleApisByRoleKeyword(roleKey string) ([]*model.SysApis, error) {
	return uc.repo.GetRoleApisByRoleKeyword(roleKey)
}

//根据角色ID获取角色
func (uc *SysRoleUsecase) GetRoleApisList(ctx context.Context, roleIds []int64) ([]*model.SysRole, error) {
	return uc.repo.GetRoleApisList(ctx, roleIds)
}

// 更新角色的权限接口（先全部删除再新增）
func (uc *SysRoleUsecase) UpdateRoleApis(roleKey string, reqRolePolicies [][]string) error {
	return uc.repo.UpdateRoleApis(roleKey, reqRolePolicies)
}

//所有开启的角色列表
func (uc *SysRoleUsecase) AllListRole(ctx context.Context) ([]*model.SysRole, error) {
	return uc.repo.AllListRole(ctx)
}

//根据角色ID获取
func (uc *SysRoleUsecase) ListRoleByIds(ctx context.Context, ids []int64) ([]*model.SysRole, error) {
	return uc.repo.ListRoleByIds(ctx, ids)
}

//获取策略中的所有授权规则
func (uc *SysRoleUsecase) GetFilteredPolicy(roleKey string) [][]string {
	return uc.repo.GetFilteredPolicy(roleKey)
}

func (uc *SysRoleUsecase) CheckCasbinRole(subs []string, obj string, act string) bool {
	return uc.repo.CheckCasbinRole(subs, obj, act)
}

//获取按钮权限列表
func (uc *SysRoleUsecase) GetPermissionsByRoleKey(roleKey string) ([]string, error) {
	return uc.repo.GetPermissionsByRoleKey(roleKey)
}

//=========
